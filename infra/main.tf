# locals {
#    lambda_build_path = "${path.module}/../../src/lambda/build/safenotes.zip"
# }

locals {
   lambda_build_path = "/../src/build/safenotes.zip"
}

output "lambda_zip_path" {
  value = "${path.module}/${local.lambda_build_path}"
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
    lambda   = var.localstack_endpoint
    iam      = var.localstack_endpoint
  }
}


resource "aws_dynamodb_table" "safenotes" {
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


# IAM role for Lambda
resource "aws_iam_role" "lambda_role_to_access_dynamodb" {
  name = "lambda_role_to_access_dynamodb"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# IAM policy for Lambda to access DynamoDB
resource "aws_iam_role_policy" "lambda_policy" {
  name = "lambda_dynamodb_policy"
  role = aws_iam_role.lambda_role_to_access_dynamodb.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Scan",
          "dynamodb:Query"
        ]
        Resource = aws_dynamodb_table.safenotes.arn
      }
    ]
  })
}

# Lambda function
resource "aws_lambda_function" "crud_lambda" {
  filename      = "${path.module}${local.lambda_build_path}" # TODO : if the .zip file is outside the current directory we need to add this as prefix ${path.module}
  function_name = "safenotes"
  role          = aws_iam_role.lambda_role_to_access_dynamodb.arn
  handler       = "safenotes" // the name of the build binary 
  runtime       = "go1.x"

  source_code_hash = filebase64sha256("${path.module}${local.lambda_build_path}")

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = aws_dynamodb_table.safenotes.name
    }
  }
}
