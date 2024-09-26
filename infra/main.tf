variable "access_key" {
    type = string 
}

variable "secret_access_key" {
    type = string
}

variable "region" {
    type = string
}

variable "localstack_endpoint" {
    type = string
}

variable "dynamodb_table_name" {
  type = string
}


provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_access_key
  region = var.region

  skip_credentials_validation = true
  skip_metadata_api_check = true
  skip_requesting_account_id = true


  endpoints {
     sqs = var.localstack_endpoint
  }
}


resource "aws_dynamodb_table" "test_table" {
  name           = var.dynamodb_table_name
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "createdAt"
    type = "S" # aws.String(time.Now().Format(time.RFC3339))
  }

    attribute {
      name = "content"
      type = "S"
    }

    attribute {
      name = "password"
      type = "S"
    }

    attribute {
      name = "url"
      type = "S"
    }
}