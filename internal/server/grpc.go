package server

import (
	"context"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	peer2 "google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpc_api "github.com/koneal2013/sensorsphere/api/v1/grpc"
	"github.com/koneal2013/sensorsphere/internal/db"
	"github.com/koneal2013/sensorsphere/internal/models"
)

const (
	objectWildCard = "*"
)

type Authorizer interface {
	Authorize(subject, object, action string) error
}

type GrpcConfig struct {
	Db db.Database
	Authorizer
}

func NewGRPCServer(config *GrpcConfig, opts ...grpc.ServerOption) (*grpc.Server, error) {
	logger := zap.L().Named("grpc_server")
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(
			func(duration time.Duration) zapcore.Field {
				return zap.Int64("grpc.time_ns", duration.Nanoseconds())
			}),
	}
	opts = append(opts, grpc.StreamInterceptor(
		grpcmiddleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger, zapOpts...),
			grpcauth.StreamServerInterceptor(authenticate),
			otelgrpc.StreamServerInterceptor(),
		)), grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
		grpcauth.UnaryServerInterceptor(authenticate),
		otelgrpc.UnaryServerInterceptor(),
	)))
	gsrv := grpc.NewServer(opts...)
	if srv, err := newGrpcServer(config); err != nil {
		return nil, err
	} else {
		grpc_api.RegisterSensorSphereServiceServer(gsrv, srv)
		return gsrv, nil
	}
}

type grpcServer struct {
	grpc_api.UnimplementedSensorSphereServiceServer
	*GrpcConfig
	grpcTracer trace.Tracer
	database   db.Database
}

func newGrpcServer(config *GrpcConfig) (srv *grpcServer, err error) {
	srv = &grpcServer{
		GrpcConfig: config,
		grpcTracer: otel.GetTracerProvider().Tracer("GrpcTracer"),
		database:   config.Db,
	}
	return srv, nil
}

func (s *grpcServer) CreateSensor(ctx context.Context, in *grpc_api.Sensor) (*grpc_api.Sensor, error) {
	ctx, span := s.grpcTracer.Start(ctx, "CreateSensor")
	defer span.End()

	newSensor := apiSensorToModel(in)
	sensor, err := s.database.CreateSensor(ctx, newSensor)
	if err != nil {
		return nil, err
	}

	return modelSensorToAPI(sensor), nil
}

func (s *grpcServer) GetSensor(ctx context.Context, in *grpc_api.GetSensorRequest) (*grpc_api.Sensor, error) {
	ctx, span := s.grpcTracer.Start(ctx, "GetSensor")
	defer span.End()

	sensor, err := s.database.GetSensor(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	return modelSensorToAPI(sensor), nil
}

func (s *grpcServer) UpdateSensor(ctx context.Context, in *grpc_api.Sensor) (*grpc_api.UpdateSensorResponse, error) {
	ctx, span := s.grpcTracer.Start(ctx, "UpdateSensor")
	defer span.End()

	updatedSensor := apiSensorToModel(in)
	rows, err := s.database.UpdateSensor(ctx, updatedSensor)
	if err != nil {
		return nil, err
	}

	return &grpc_api.UpdateSensorResponse{RowsAffected: rows}, nil
}

func (s *grpcServer) GetNearestSensor(ctx context.Context, in *grpc_api.Location) (*grpc_api.Sensor, error) {
	ctx, span := s.grpcTracer.Start(ctx, "GetNearestSensor")
	defer span.End()

	location := apiLocationToModel(in)
	sensor, err := s.database.GetNearestSensor(ctx, location)
	if err != nil {
		return nil, err
	}

	return modelSensorToAPI(sensor), nil
}

func (s *grpcServer) CreateSensorReading(ctx context.Context, in *grpc_api.SensorReading) (*grpc_api.SensorReading, error) {
	ctx, span := s.grpcTracer.Start(ctx, "CreateSensorReading")
	defer span.End()

	reading := apiReadingToModel(in)
	sensorReading, err := s.database.CreateSensorReading(ctx, reading)
	if err != nil {
		return nil, err
	}

	return modelReadingToAPI(sensorReading), nil
}

func (s *grpcServer) GetSensorReadingsForTimeRange(ctx context.Context, in *grpc_api.TimeRangeQuery) (*grpc_api.SensorReadingsResponse, error) {
	ctx, span := s.grpcTracer.Start(ctx, "GetSensorReadingsForTimeRange")
	defer span.End()

	timeRange := apiTimeRangeQueryToModel(in)
	sensorReadings, err := s.database.GetSensorReadingsForTimeRange(ctx, timeRange)
	if err != nil {
		return nil, err
	}

	return &grpc_api.SensorReadingsResponse{SensorReadings: modelReadingsToAPI(sensorReadings)}, nil
}

func apiSensorToModel(in *grpc_api.Sensor) *models.Sensor {
	return &models.Sensor{
		Name: in.Name,
		Location: models.Location{
			Longitude: in.Location.Longitude,
			Latitude:  in.Location.Latitude,
		},
		Tags: in.Tags,
	}
}

func modelSensorToAPI(sensor *models.Sensor) *grpc_api.Sensor {
	return &grpc_api.Sensor{
		Name: sensor.Name,
		Location: &grpc_api.Location{
			Longitude: sensor.Location.Longitude,
			Latitude:  sensor.Location.Latitude,
		},
		Tags: sensor.Tags,
	}
}

func apiLocationToModel(in *grpc_api.Location) *models.Location {
	return &models.Location{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
}

func apiReadingToModel(in *grpc_api.SensorReading) *models.SensorReading {
	return &models.SensorReading{
		SensorName: in.SensorName,
		Value:      in.Value,
		Time:       in.Time.AsTime(),
	}
}

func modelReadingToAPI(reading *models.SensorReading) *grpc_api.SensorReading {
	return &grpc_api.SensorReading{
		SensorName: reading.SensorName,
		Value:      reading.Value,
		Time:       timestamppb.New(reading.Time),
	}
}

func apiTimeRangeQueryToModel(in *grpc_api.TimeRangeQuery) models.TimeRangeQuery {
	return models.TimeRangeQuery{
		SensorName: in.SensorName,
		StartTime:  in.StartTime.AsTime(),
		EndTime:    in.EndTime.AsTime(),
	}
}

func modelReadingsToAPI(readings []*models.SensorReading) []*grpc_api.SensorReading {
	apiReadings := make([]*grpc_api.SensorReading, len(readings))
	for i, reading := range readings {
		apiReadings[i] = modelReadingToAPI(reading)
	}
	return apiReadings
}

func authenticate(ctx context.Context) (context.Context, error) {
	if peer, ok := peer2.FromContext(ctx); !ok {
		return ctx, status.New(codes.Unknown, "couldn't find peer info").Err()
	} else if peer.AuthInfo == nil {
		return context.WithValue(ctx, subjectContextKey{}, ""), nil
	} else {
		tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
		subject := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
		ctx = context.WithValue(ctx, subjectContextKey{}, subject)
		return ctx, nil
	}
}

type subjectContextKey struct{}

func subject(ctx context.Context) string {
	return ctx.Value(subjectContextKey{}).(string)
}
