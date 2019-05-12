// This file generated by `tools/model-rule-gen/main.go`. DO NOT EDIT

package models

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/tflint"
)

// AwsLaunchTemplateInvalidNameRule checks the pattern is valid
type AwsLaunchTemplateInvalidNameRule struct {
	resourceType  string
	attributeName string
	max           int
	min           int
	pattern       *regexp.Regexp
}

// NewAwsLaunchTemplateInvalidNameRule returns new rule with default attributes
func NewAwsLaunchTemplateInvalidNameRule() *AwsLaunchTemplateInvalidNameRule {
	return &AwsLaunchTemplateInvalidNameRule{
		resourceType:  "aws_launch_template",
		attributeName: "name",
		max:           128,
		min:           3,
		pattern:       regexp.MustCompile(`^[a-zA-Z0-9\(\)\.\-/_]+$`),
	}
}

// Name returns the rule name
func (r *AwsLaunchTemplateInvalidNameRule) Name() string {
	return "aws_launch_template_invalid_name"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLaunchTemplateInvalidNameRule) Enabled() bool {
	return true
}

// Type returns the rule severity
func (r *AwsLaunchTemplateInvalidNameRule) Type() string {
	return issue.ERROR
}

// Link returns the rule reference link
func (r *AwsLaunchTemplateInvalidNameRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsLaunchTemplateInvalidNameRule) Check(runner *tflint.Runner) error {
	log.Printf("[INFO] Check `%s` rule for `%s` runner", r.Name(), runner.TFConfigPath())

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val)

		return runner.EnsureNoError(err, func() error {
			if len(val) > r.max {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`"%s" must be %d characters or less`, val, r.max),
					attribute.Expr.Range(),
				)
			}

			if len(val) < r.min {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`"%s" must be %d characters or higher`, val, r.min),
					attribute.Expr.Range(),
				)
			}

			if !r.pattern.MatchString(val) {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`"%s" does not match valid pattern ^[a-zA-Z0-9\(\)\.\-/_]+$`, val),
					attribute.Expr.Range(),
				)
			}
			return nil
		})
	})
}
