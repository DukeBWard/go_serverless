# go_serverless
API, DynamoDB, Lambda

Note: proper build command for lambda deployment is given in AWS documentation (Building for amd64).
``GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go``
