// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/sensor_readings": {
            "get": {
                "description": "Get sensor readings for a specific time range",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sensor_readings"
                ],
                "summary": "Get sensor readings for a time range",
                "parameters": [
                    {
                        "description": "Time range query",
                        "name": "timeRangeQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TimeRangeQuery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.SensorReading"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new sensor reading with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sensor_readings"
                ],
                "summary": "Create a new sensor reading",
                "parameters": [
                    {
                        "description": "Create sensor reading",
                        "name": "sensorReading",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SensorReading"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SensorReading"
                        }
                    }
                }
            }
        },
        "/sensors": {
            "post": {
                "description": "Create a new sensor with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sensors"
                ],
                "summary": "Create a new sensor",
                "parameters": [
                    {
                        "description": "Create sensor",
                        "name": "sensor",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Sensor"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Sensor"
                        }
                    }
                }
            }
        },
        "/sensors/nearest": {
            "get": {
                "description": "Get the nearest sensor to a specific location",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sensors"
                ],
                "summary": "Get the nearest sensor",
                "parameters": [
                    {
                        "description": "Location",
                        "name": "location",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Location"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Sensor"
                        }
                    }
                }
            }
        },
        "/sensors/{name}": {
            "get": {
                "description": "Get a sensor by its name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sensors"
                ],
                "summary": "Get a sensor",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Sensor name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Sensor"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a sensor with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sensors"
                ],
                "summary": "Update a sensor",
                "parameters": [
                    {
                        "description": "Update sensor",
                        "name": "sensor",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Sensor"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "description": "Returns 200 OK if server is ready to accept requests",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Get server status",
                "responses": {
                    "200": {
                        "description": "Server is running",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Location": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "models.Sensor": {
            "type": "object",
            "properties": {
                "location": {
                    "$ref": "#/definitions/models.Location"
                },
                "name": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.SensorReading": {
            "type": "object",
            "properties": {
                "sensorName": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "models.TimeRangeQuery": {
            "type": "object",
            "properties": {
                "endTime": {
                    "type": "string"
                },
                "sensorName": {
                    "type": "string"
                },
                "startTime": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
