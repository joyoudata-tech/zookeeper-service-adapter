package adapter

import (
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type SchemaGenerator struct {}

func (s SchemaGenerator) GeneratePlanSchema(plan serviceadapter.Plan) (serviceadapter.PlanSchema, error) {
	bindSchema := serviceadapter.JSONSchemas{
		Parameters: map[string]interface{}{
			"$schema":              "http://json-schema.org/draft-04/schema#",
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]interface{}{
				"director": map[string]interface{}{
					"description": "zookeeper create director.",
					"type":        "string",
				},
			},
		},
	}
	return serviceadapter.PlanSchema{
		ServiceBinding: serviceadapter.ServiceBindingSchema{
			Create: bindSchema,
		},
	},nil
}