package detector

import (
	"fmt"

	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/schema"
)

type AwsInstancePreviousTypeDetector struct {
	*Detector
	previousInstanceTypes map[string]bool
}

func (d *Detector) CreateAwsInstancePreviousTypeDetector() *AwsInstancePreviousTypeDetector {
	nd := &AwsInstancePreviousTypeDetector{
		Detector:              d,
		previousInstanceTypes: map[string]bool{},
	}
	nd.Name = "aws_instance_previous_type"
	nd.IssueType = issue.WARNING
	nd.TargetType = "resource"
	nd.Target = "aws_instance"
	nd.DeepCheck = false
	nd.Link = "https://github.com/wata727/tflint/blob/master/docs/aws_instance_previous_type.md"
	nd.Enabled = true
	return nd
}

func (d *AwsInstancePreviousTypeDetector) PreProcess() {
	d.previousInstanceTypes = map[string]bool{
		"t1.micro":    true,
		"m1.small":    true,
		"m1.medium":   true,
		"m1.large":    true,
		"m1.xlarge":   true,
		"c1.medium":   true,
		"c1.xlarge":   true,
		"c3.large":    true,
		"c3.xlarge":   true,
		"c3.2xlarge":  true,
		"c3.4xlarge":  true,
		"c3.8xlarge":  true,
		"cc2.8xlarge": true,
		"cg1.4xlarge": true,
		"m2.xlarge":   true,
		"m2.2xlarge":  true,
		"m2.4xlarge":  true,
		"m3.medium":   true,
		"m3.large":    true,
		"m3.xlarge":   true,
		"m3.2xlarge":  true,
		"cr1.8xlarge": true,
		"r3.large":    true,
		"r3.xlarge":   true,
		"r3.2xlarge":  true,
		"r3.4xlarge":  true,
		"r3.8xlarge":  true,
		"hi1.4xlarge": true,
		"hs1.8xlarge": true,
		"i2.xlarge":   true,
		"i2.2xlarge":  true,
		"i2.4xlarge":  true,
		"i2.8xlarge":  true,
		"g2.2xlarge":  true,
		"g2.8xlarge":  true,
	}
}

func (d *AwsInstancePreviousTypeDetector) Detect(resource *schema.Resource, issues *[]*issue.Issue) {
	instanceTypeToken, ok := resource.GetToken("instance_type")
	if !ok {
		return
	}
	instanceType, err := d.evalToString(instanceTypeToken.Text)
	if err != nil {
		d.Logger.Error(err)
		return
	}

	if d.previousInstanceTypes[instanceType] {
		issue := &issue.Issue{
			Detector: d.Name,
			Type:     d.IssueType,
			Message:  fmt.Sprintf("\"%s\" is previous generation instance type.", instanceType),
			Line:     instanceTypeToken.Pos.Line,
			File:     instanceTypeToken.Pos.Filename,
			Link:     d.Link,
		}
		*issues = append(*issues, issue)
	}
}
