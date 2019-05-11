# aws_instance_default_standard_volume

Disallow using default volume type.

## Example

```hcl
resource "aws_instance" "web" {
  ami                  = "ami-b73b63a0"
  instance_type        = "m4.2xlarge"
  iam_instance_profile = "app-service"

  root_block_device = {
    volume_size = "16"
  }
}
```

```
$ tflint
template.tf
        WARNING:6 "volume_type" is not specified. Default standard volume type is not recommended. You can use "gp2", "io1", etc instead. (aws_instance_default_standard_volume)

Result: 1 issues  (0 errors , 1 warnings , 0 notices)
```

## Why

If you use EBS as instance volume, you can specify the volume type. If not specified, the "default" volume type will be used. This is an officially deprecated volume type, and it is generally recommended to use "gp2".

## How To Fix

Check the [EBS volume types](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html) and specify volume type. If you want to use the "default", if you explicitly specify "default", TFLint will not report this issue.
