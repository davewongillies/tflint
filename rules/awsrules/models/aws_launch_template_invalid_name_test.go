package models

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

func Test_AwsLaunchTemplateInvalidNameRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected issue.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_launch_template" "foo" {
	name = "foo[bar]"
}`,
			Expected: []*issue.Issue{
				{
					Detector: "aws_launch_template_invalid_name",
					Type:     "ERROR",
					Message:  `"foo[bar]" does not match valid pattern ^[a-zA-Z0-9\(\)\.\-/_]+$`,
					Line:     3,
					File:     "resource.tf",
				},
			},
		},
		{
			Name: "It is too short",
			Content: `
resource "aws_launch_template" "foo" {
	name = "go"
}`,
			Expected: []*issue.Issue{
				{
					Detector: "aws_launch_template_invalid_name",
					Type:     "ERROR",
					Message:  `"go" must be 3 characters or higher`,
					Line:     3,
					File:     "resource.tf",
				},
			},
		},
		{
			Name: "It is too long",
			Content: `
resource "aws_launch_template" "foo" {
	name = "Lorem_ipsum_dolor_sit_amet_consectetur_adipisicing_elit_sed_do_eiusmod_tempor_incididunt_ut_labore_et_dolore_magna_aliqua.Ut_enim_ad_minim_veniam_quis_nostrud_exercitation_ullamco_laboris_nisi_ut_aliquip_ex_ea_commodo_consequat.Duis_aute_irure_dolor_in_reprehenderit_in_voluptate_velit_esse_cillum_dolore_eu_fugiat_nulla_pariatur.Excepteur_sint_occaecat_cupidatat_non_proident_sunt_in_culpa_qui_officia_deserunt_mollit_anim_id_est_laborum."
}`,
			Expected: []*issue.Issue{
				{
					Detector: "aws_launch_template_invalid_name",
					Type:     "ERROR",
					Message:  `"Lorem_ipsum_dolor_sit_amet_consectetur_adipisicing_elit_sed_do_eiusmod_tempor_incididunt_ut_labore_et_dolore_magna_aliqua.Ut_enim_ad_minim_veniam_quis_nostrud_exercitation_ullamco_laboris_nisi_ut_aliquip_ex_ea_commodo_consequat.Duis_aute_irure_dolor_in_reprehenderit_in_voluptate_velit_esse_cillum_dolore_eu_fugiat_nulla_pariatur.Excepteur_sint_occaecat_cupidatat_non_proident_sunt_in_culpa_qui_officia_deserunt_mollit_anim_id_est_laborum." must be 128 characters or less`,
					Line:     3,
					File:     "resource.tf",
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_launch_template" "foo" {
	name = "foo"
}`,
			Expected: []*issue.Issue{},
		},
	}

	dir, err := ioutil.TempDir("", "AwsLaunchTemplateInvalidNameRule")
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
		rule := NewAwsLaunchTemplateInvalidNameRule()

		if err = rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		if !cmp.Equal(tc.Expected, runner.Issues) {
			t.Fatalf("Expected issues are not matched:\n %s\n", cmp.Diff(tc.Expected, runner.Issues))
		}
	}
}
