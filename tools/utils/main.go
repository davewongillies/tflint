package utils

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// ToCamel converts a string to CamelCase
func ToCamel(str string) string {
	exceptions := map[string]string{
		"ami":         "AMI",
		"db":          "DB",
		"alb":         "ALB",
		"elb":         "ELB",
		"vpc":         "VPC",
		"elasticache": "ElastiCache",
		"iam":         "IAM",
	}
	for pattern, conv := range exceptions {
		str = strings.Replace(str, "_"+pattern+"_", "_"+conv+"_", -1)
		str = strings.Replace(str, pattern+"_", conv+"_", -1)
		str = strings.Replace(str, "_"+pattern, "_"+conv, -1)
	}
	return strcase.ToCamel(str)
}
