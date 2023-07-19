FROM golang:1.20.6-alpine as build
# Set the working directory
WORKDIR /go/src/sensorsphere
# Copy and download dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download
# Copy the source files from the host
COPY . /go/src/sensorsphere

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/sensorsphere ./cmd/sensorsphere

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/sensorsphere /bin/sensorsphere
ENTRYPOINT ["/bin/sensorsphere", "--config-file=/etc/config.json"]
