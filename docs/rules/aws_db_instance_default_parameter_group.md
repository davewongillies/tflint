# aws_db_instance_default_parameter_group

Disallow using default DB parameter group.

## Example

```hcl
resource "aws_db_instance" "mysql" {
  identifier             = "app"
  allocated_storage      = 50
  storage_type           = "gp2"
  engine                 = "mysql"
  engine_version         = "5.7.11"
  instance_class         = "db.m4.large"
  name                   = "app_db"
  port                   = 3306
  publicly_accessible    = false
  vpc_security_group_ids = ["${aws_security_group.mysql.id}"]
  db_subnet_group_name   = "app-subnet-group"
  parameter_group_name   = "default.mysql5.7" // default DB parameter group!
  multi_az               = true
}
```

```
$ tflint
template.tf
        NOTICE:12 "default.mysql5.7" is default parameter group. You cannot edit it. (aws_db_instance_default_parameter_group)

Result: 1 issues  (0 errors , 0 warnings , 1 notices)
```

## Why

You can modify parameter values in a custom DB parameter group, but you can't change the parameter values in a default DB parameter group.

## How To Fix

Create a new parameter group, and change the `parameter_group_name` to that.
