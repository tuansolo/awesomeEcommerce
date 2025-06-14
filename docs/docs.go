// Package docs provides Swagger documentation for the API
package docs

import "github.com/swaggo/swag"

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Awesome E-commerce API",
	Description:      "This is the API documentation for the Awesome E-commerce application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

// This is a placeholder for the generated Swagger documentation
const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "This is the API documentation for the Awesome E-commerce application",
        "title": "Awesome E-commerce API",
        "contact": {},
        "version": "1.0"
    },
    "host": "",
    "basePath": "/api/v1",
    "paths": {}
}`
