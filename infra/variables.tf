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

variable "lambda_build_path" {
    type = string
    default = "${path.module}/../../src/lambda/build/safenotes.zip"
}
