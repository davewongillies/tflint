mapping {
  resource {
    type      = "aws_acm_certificate"
    attribute = "certificate_body"
  }
  model {
    path  = "aws-sdk-go/models/apis/acm/2015-12-08/api-2.json"
    shape = "CertificateBody"
  }
}

mapping {
  resource {
    type      = "aws_acm_certificate"
    attribute = "certificate_chain"
  }
  model {
    path  = "aws-sdk-go/models/apis/acm/2015-12-08/api-2.json"
    shape = "CertificateChain"
  }
}

mapping {
  resource {
    type      = "aws_acm_certificate"
    attribute = "private_key"
  }
  model {
    path  = "aws-sdk-go/models/apis/acm/2015-12-08/api-2.json"
    shape = "PrivateKey"
  }
}

mapping {
  resource {
    type      = "aws_launch_template"
    attribute = "name"
  }
  model {
    path  = "aws-sdk-go/models/apis/ec2/2016-11-15/api-2.json"
    shape = "LaunchTemplateName"
  }
}
