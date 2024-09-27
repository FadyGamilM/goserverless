variable "aws_access_key" {
    type = string 
}

variable "aws_secret_access_key" {
    type = string
}

variable "aws_region" {
    type = string
}

variable "localstack_endpoint" {
    type = string
}

variable "dynamodb_table_name" {
  type = string
}

// if we need to add defualt_value and this value contains the $path.module attribute we cannot use it here in variables and access the variable using var.var_name, because this is dynamically and terraform doesn't support this in the default values, instead create a local variable
variable "lambda_build_path" {
    type = string
}
