package awsrules

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/configload"
	"github.com/hashicorp/terraform/terraform"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/tflint"
)

func Test_AwsRouteNotSpecifiedTarget(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected issue.Issues
	}{
		{
			Name: "route target is not specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
}`,
			Expected: []*issue.Issue{
				{
					Detector: "aws_route_not_specified_target",
					Type:     "ERROR",
					Message:  "The routing target is not specified, each routing must contain either a gateway_id, egress_only_gateway_id a nat_gateway_id, an instance_id or a vpc_peering_connection_id or a network_interface_id.",
					Line:     2,
					File:     "resource.tf",
					Link:     "https://github.com/wata727/tflint/blob/master/docs/aws_route_not_specified_target.md",
				},
			},
		},
		{
			Name: "gateway_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    gateway_id = "igw-1234abcd"
}`,
			Expected: []*issue.Issue{},
		},
		{
			Name: "egress_only_gateway_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    egress_only_gateway_id = "eigw-1234abcd"
}`,
			Expected: []*issue.Issue{},
		},
		{
			Name: "nat_gateway_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    nat_gateway_id = "nat-1234abcd"
}`,
			Expected: []*issue.Issue{},
		},
		{
			Name: "instance_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    instance_id = "i-1234abcd"
}`,
			Expected: []*issue.Issue{},
		},
		{
			Name: "vpc_peering_connection_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    vpc_peering_connection_id = "pcx-1234abcd"
}`,
			Expected: []*issue.Issue{},
		},
		{
			Name: "network_interface_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    network_interface_id = "eni-1234abcd"
}`,
			Expected: []*issue.Issue{},
		},
	}

	dir, err := ioutil.TempDir("", "AwsRouteNotSpecifiedTarget")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, tc := range cases {
		loader, err := configload.NewLoader(&configload.Config{})
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(dir+"/resource.tf", []byte(tc.Content), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		mod, diags := loader.Parser().LoadConfigDir(dir)
		if diags.HasErrors() {
			t.Fatal(diags)
		}
		cfg, tfdiags := configs.BuildConfig(mod, configs.DisabledModuleWalker)
		if tfdiags.HasErrors() {
			t.Fatal(tfdiags)
		}

		runner := tflint.NewRunner(tflint.EmptyConfig(), cfg, map[string]*terraform.InputValue{})
		rule := NewAwsRouteNotSpecifiedTargetRule()

		if err = rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		if !cmp.Equal(tc.Expected, runner.Issues) {
			t.Fatalf("Expected issues are not matched:\n %s\n", cmp.Diff(tc.Expected, runner.Issues))
		}
	}
}
