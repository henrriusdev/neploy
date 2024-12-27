// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "6oMwz@example.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/applications": {
            "get": {
                "description": "List all applications",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Application"
                ],
                "summary": "List all applications",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Application"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Application"
                ],
                "summary": "Create a new application",
                "parameters": [
                    {
                        "description": "Application details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateApplicationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/applications/{id}": {
            "get": {
                "description": "Get an application by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Application"
                ],
                "summary": "Get an application by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Application ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Application"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Application"
                ],
                "summary": "Delete an application",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Application ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/applications/{id}/deploy": {
            "post": {
                "description": "Deploy an application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Application"
                ],
                "summary": "Deploy an application",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Application ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Deployment details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.DeployApplicationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/applications/{id}/start": {
            "post": {
                "description": "Start an application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Application"
                ],
                "summary": "Start an application",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Application ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/applications/{id}/stop": {
            "post": {
                "description": "Stop an application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Application"
                ],
                "summary": "Stop an application",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Application ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/github": {
            "get": {
                "description": "Start GitHub OAuth flow",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Start GitHub OAuth flow",
                "parameters": [
                    {
                        "type": "string",
                        "description": "State parameter",
                        "name": "state",
                        "in": "query"
                    }
                ],
                "responses": {
                    "302": {
                        "description": "Redirects to GitHub OAuth flow",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/github/callback": {
            "get": {
                "description": "Handle GitHub OAuth callback",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Handle GitHub OAuth callback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "State parameter",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Authorization code",
                        "name": "code",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.OAuthResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/gitlab": {
            "get": {
                "description": "Start GitLab OAuth flow",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Start GitLab OAuth flow",
                "parameters": [
                    {
                        "type": "string",
                        "description": "State parameter",
                        "name": "state",
                        "in": "query"
                    }
                ],
                "responses": {
                    "302": {
                        "description": "Redirects to GitLab OAuth flow",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/gitlab/callback": {
            "get": {
                "description": "Handle GitLab OAuth callback",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Handle GitLab OAuth callback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "State parameter",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Authorization code",
                        "name": "code",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.OAuthResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/gateways": {
            "post": {
                "description": "Create a new gateway",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gateway"
                ],
                "summary": "Create a new gateway",
                "parameters": [
                    {
                        "description": "Gateway details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Gateway"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/gateways/app/{appId}": {
            "get": {
                "description": "List gateways by app ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gateway"
                ],
                "summary": "List gateways by app ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "App ID",
                        "name": "appId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Gateway"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/gateways/{id}": {
            "get": {
                "description": "Get a gateway by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gateway"
                ],
                "summary": "Get a gateway by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Gateway ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Gateway"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "put": {
                "description": "Update a gateway",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gateway"
                ],
                "summary": "Update a gateway",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Gateway ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Gateway details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Gateway"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a gateway",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gateway"
                ],
                "summary": "Delete a gateway",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Gateway ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/gateways/{id}/health": {
            "get": {
                "description": "Check the health of a gateway",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gateway"
                ],
                "summary": "Check the health of a gateway",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Gateway ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Login a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth",
                    "User"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "Login Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "description": "Logout a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth",
                    "User"
                ],
                "summary": "Logout a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Application": {
            "type": "object",
            "properties": {
                "appName": {
                    "type": "string"
                },
                "createdAt": {
                    "$ref": "#/definitions/model.Date"
                },
                "deletedAt": {
                    "$ref": "#/definitions/model.Date"
                },
                "deployLocation": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "storageLocation": {
                    "type": "string"
                },
                "techStackId": {
                    "type": "string"
                },
                "updatedAt": {
                    "$ref": "#/definitions/model.Date"
                }
            }
        },
        "model.CreateApplicationRequest": {
            "type": "object",
            "properties": {
                "appName": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "techStack": {
                    "type": "string"
                }
            }
        },
        "model.Date": {
            "type": "object",
            "properties": {
                "time.Time": {
                    "type": "string"
                }
            }
        },
        "model.DeployApplicationRequest": {
            "type": "object",
            "properties": {
                "repoUrl": {
                    "type": "string"
                }
            }
        },
        "model.Gateway": {
            "type": "object",
            "properties": {
                "applicationId": {
                    "type": "string"
                },
                "createdAt": {
                    "$ref": "#/definitions/model.Date"
                },
                "deletedAt": {
                    "$ref": "#/definitions/model.Date"
                },
                "domain": {
                    "type": "string"
                },
                "endpointType": {
                    "description": "\"subdomain\" or \"path\"",
                    "type": "string"
                },
                "endpointUrl": {
                    "type": "string"
                },
                "httpMethod": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "integrationType": {
                    "type": "string"
                },
                "loggingLevel": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "port": {
                    "type": "string"
                },
                "stage": {
                    "type": "string"
                },
                "status": {
                    "description": "\"active\", \"inactive\", \"error\"",
                    "type": "string"
                },
                "subdomain": {
                    "type": "string"
                },
                "updatedAt": {
                    "$ref": "#/definitions/model.Date"
                }
            }
        },
        "model.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 8
                }
            }
        },
        "model.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/model.User"
                }
            }
        },
        "model.OAuthResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "provider": {
                    "$ref": "#/definitions/model.Provider"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.Provider": {
            "type": "string",
            "enum": [
                "github",
                "gitlab"
            ],
            "x-enum-varnames": [
                "Github",
                "Gitlab"
            ]
        },
        "model.User": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "createdAt": {
                    "$ref": "#/definitions/model.Date"
                },
                "deletedAt": {
                    "$ref": "#/definitions/model.Date"
                },
                "dob": {
                    "$ref": "#/definitions/model.Date"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "updatedAt": {
                    "$ref": "#/definitions/model.Date"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Neploy API",
	Description:      "Neploy API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
