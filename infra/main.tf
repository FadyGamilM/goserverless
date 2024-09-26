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


provider "aws" {
  access_key = var.aws_access_key
  secret_key = var.aws_secret_access_key
  region = var.aws_region

  skip_credentials_validation = true
  skip_metadata_api_check = true
  skip_requesting_account_id = true


  endpoints {
    dynamodb = var.localstack_endpoint
  }
}


resource "aws_dynamodb_table" "test_table" {
  name           = var.dynamodb_table_name
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"
    range_key      = "createdAt"  # Using createdAt as the sort key


  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "createdAt"
    type = "S" # aws.String(time.Now().Format(time.RFC3339))
  }

#! Why commented ? becasue in terraform we don't have to define all attributed, we just define the hash key and the sort key only 
  # attribute {
  #   name = "content"
  #   type = "S"
  # }

  # attribute {
  #   name = "password"
  #   type = "S"
  # }

  # attribute {
  #   name = "url"
  #   type = "S"
  # }
}