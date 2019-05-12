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

func Test_AwsAcmCertificateInvalidCertificateBodyRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected issue.Issues
	}{
		{
			Name: "CSR key",
			Content: `
resource "aws_acm_certificate" "cert" {
	private_key      = "${tls_private_key.example.private_key_pem}"
	certificate_body = <<CERTIFICATE
-----BEGIN CERTIFICATE REQUEST-----
MIICjjCCAXYCAQAwSTELMAkGA1UEBhMCSkExCjAIBgNVBAgMAWExCjAIBgNVBAcM
AWIxCjAIBgNVBAoMAWExCjAIBgNVBAsMAWExCjAIBgNVBAMMAWEwggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQC5jqT6lQEWbyNeHMEeLJxRN7wPmfJXfWSJ
L42vYPM8Ny8NxUtUE8eYavsJQSp2PcwBxVYatgVPNeL3saPGxfGpkIuaPgQO1n7a
+O+3Vb/GANe3g9RX/3p280DHdm+pppRp68yhivMrtznjXJCebOix+KdgHBMLUqx2
7FWmqchJ+vQ74wmuDR1Y5fh69NDn79kB8ZJiUpZWQ0CPrgoi8KxWU8FT0JnQxBvE
6CJ81P2/1LtJ//ngasyux37j3R64Q4ZLgTX9VtAX+Bnoy9Nh9wTb6zbjh6xzO1bN
4IwaFvc+y2F5gZvct6m53p/DfskII+WuH6Gc6nHrN/g/V49Vo+3ZAgMBAAGgADAN
BgkqhkiG9w0BAQsFAAOCAQEAA4W/lkp3oTmjIoyhZxUMv7b1zcRU/s9juzvYdfMB
nkty65GIKc8VgRSdgdXHg9LyAmG2fw/Ek7fHzMb10a6AR6nNn8dDmDSJgP/Li/qH
65ufOAZFwaQESmaOKuixXzpOl55k4iJCgWng1ejxZ1CSQczWdchLgW6af+ykUgLK
i2H5CazWnCBtBRonsDKFE6TYH0NEqdFE/kAyWtKiMOXAV8Jyr2p8K5hMG/8Cusux
Oe04sLexs2p1Og6LKAv9aWk0wYKB15Zjgx1EqKGJOwHJ5pOVXyGiQAnkqGaC0Q4N
EUNkhA1s4v7yBuNuulIfhcbyOeLwnzElTz5RrV/1hgMWMg==
-----END CERTIFICATE REQUEST-----
CERTIFICATE
}`,
			Expected: []*issue.Issue{
				{
					Detector: "aws_acm_certificate_invalid_certificate_body",
					Type:     "ERROR",
					Message: `"-----BEGIN CERTIFICATE REQUEST-----
MIICjjCCAXYCAQAwSTELMAkGA1UEBhMCSkExCjAIBgNVBAgMAWExCjAIBgNVBAcM
AWIxCjAIBgNVBAoMAWExCjAIBgNVBAsMAWExCjAIBgNVBAMMAWEwggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQC5jqT6lQEWbyNeHMEeLJxRN7wPmfJXfWSJ
L42vYPM8Ny8NxUtUE8eYavsJQSp2PcwBxVYatgVPNeL3saPGxfGpkIuaPgQO1n7a
+O+3Vb/GANe3g9RX/3p280DHdm+pppRp68yhivMrtznjXJCebOix+KdgHBMLUqx2
7FWmqchJ+vQ74wmuDR1Y5fh69NDn79kB8ZJiUpZWQ0CPrgoi8KxWU8FT0JnQxBvE
6CJ81P2/1LtJ//ngasyux37j3R64Q4ZLgTX9VtAX+Bnoy9Nh9wTb6zbjh6xzO1bN
4IwaFvc+y2F5gZvct6m53p/DfskII+WuH6Gc6nHrN/g/V49Vo+3ZAgMBAAGgADAN
BgkqhkiG9w0BAQsFAAOCAQEAA4W/lkp3oTmjIoyhZxUMv7b1zcRU/s9juzvYdfMB
nkty65GIKc8VgRSdgdXHg9LyAmG2fw/Ek7fHzMb10a6AR6nNn8dDmDSJgP/Li/qH
65ufOAZFwaQESmaOKuixXzpOl55k4iJCgWng1ejxZ1CSQczWdchLgW6af+ykUgLK
i2H5CazWnCBtBRonsDKFE6TYH0NEqdFE/kAyWtKiMOXAV8Jyr2p8K5hMG/8Cusux
Oe04sLexs2p1Og6LKAv9aWk0wYKB15Zjgx1EqKGJOwHJ5pOVXyGiQAnkqGaC0Q4N
EUNkhA1s4v7yBuNuulIfhcbyOeLwnzElTz5RrV/1hgMWMg==
-----END CERTIFICATE REQUEST-----
" does not match valid pattern ^-{5}BEGIN CERTIFICATE-{5}\x{000D}?\x{000A}([A-Za-z0-9/+]{64}\x{000D}?\x{000A})*[A-Za-z0-9/+]{1,64}={0,2}\x{000D}?\x{000A}-{5}END CERTIFICATE-{5}(\x{000D}?\x{000A})?$`,
					Line: 4,
					File: "resource.tf",
				},
			},
		},
		{
			Name: "CRT key",
			Content: `
resource "aws_acm_certificate" "cert" {
	private_key      = "${tls_private_key.example.private_key_pem}"
	certificate_body = <<CERTIFICATE
-----BEGIN CERTIFICATE-----
MIIDDjCCAfYCCQCMlVDEcxV0gDANBgkqhkiG9w0BAQUFADBJMQswCQYDVQQGEwJK
QTEKMAgGA1UECAwBYTEKMAgGA1UEBwwBYjEKMAgGA1UECgwBYTEKMAgGA1UECwwB
YTEKMAgGA1UEAwwBYTAeFw0xOTA1MTQxNTUxMjhaFw0yOTA1MTExNTUxMjhaMEkx
CzAJBgNVBAYTAkpBMQowCAYDVQQIDAFhMQowCAYDVQQHDAFiMQowCAYDVQQKDAFh
MQowCAYDVQQLDAFhMQowCAYDVQQDDAFhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAuY6k+pUBFm8jXhzBHiycUTe8D5nyV31kiS+Nr2DzPDcvDcVLVBPH
mGr7CUEqdj3MAcVWGrYFTzXi97GjxsXxqZCLmj4EDtZ+2vjvt1W/xgDXt4PUV/96
dvNAx3ZvqaaUaevMoYrzK7c541yQnmzosfinYBwTC1KsduxVpqnISfr0O+MJrg0d
WOX4evTQ5+/ZAfGSYlKWVkNAj64KIvCsVlPBU9CZ0MQbxOgifNT9v9S7Sf/54GrM
rsd+490euEOGS4E1/VbQF/gZ6MvTYfcE2+s244escztWzeCMGhb3PstheYGb3Lep
ud6fw37JCCPlrh+hnOpx6zf4P1ePVaPt2QIDAQABMA0GCSqGSIb3DQEBBQUAA4IB
AQCoj/sZfrypif6AoLkqg2WimmK2KvWNf4srEVgI8BBIpnQpmvYdMKm4IBta8eWO
E9Sdh2u8dnTpn9TEwK/hJpisRZey7H4pPXde86QHmJF1YjF+gdwgpsayIHsfCYJ9
LJxew68jxO9YANwHy6RlS3c+hcNIWfSMOoct/P6vVkcMKOgA/hiMfHELlMzBK68U
r+Ae7wRjNF4Whbxc6bdTOLocmhOjy6IvPC8x6K5RdOYaxVpRNgUz6WgQUI1gZ3hu
YjSaGdqonttvSXhhSnoQEAHIpxvHq/PjOc5qEnzOt9nlYp3Ohr6WQAZfF3iwdd3l
Q2V76qkXAhIjADC7VpZKJiij
-----END CERTIFICATE-----
CERTIFICATE
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
		rule := NewAwsAcmCertificateInvalidCertificateBodyRule()

		if err = rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		if !cmp.Equal(tc.Expected, runner.Issues) {
			t.Fatalf("Expected issues are not matched:\n %s\n", cmp.Diff(tc.Expected, runner.Issues))
		}
	}
}
