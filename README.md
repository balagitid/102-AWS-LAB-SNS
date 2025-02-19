# 102-AWS-LAB-SNS

## Steps to follow: (manual)

### Install dependencies:

```
vi main.go
go mod init sns-demo
go get github.com/aws/aws-sdk-go
go get github.com/joho/godotenv
```

###  Create a .env file in your project root with:
```
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1
```

### Run the application:

go run main.go

## Steps to follow: Docker build and Run

```
docker build -t sns-demo-app .
docker run -p 8080:8080  sns-demo-app
```
