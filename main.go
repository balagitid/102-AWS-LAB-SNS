package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/joho/godotenv"
)

type PublishRequest struct {
	Message string `json:"message"`
}

var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>AWS SNS Demo</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
    <style>
        body {
            background-color: #f8f9fa;
        }
        .container {
            margin-top: 50px;
        }
        .top-header {
            background-color: white;
            padding: 10px;
            text-align: center;
        }
        .top-header img {
            max-height: 80px;
        }
        .running-text {
            background-color: #0b5ed7;
            color: white;
            padding: 10px;
            overflow: hidden;
            white-space: nowrap;
            box-sizing: border-box;
        }
        .running-text p {
            display: inline-block;
            padding-left: 100%;
            animation: marquee 15s linear infinite;
        }
        @keyframes marquee {
            0%   { transform: translate(0, 0); }
            100% { transform: translate(-100%, 0); }
        }
        .form-section {
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            color: gray;
        }
    </style>
</head>
<body>
    <div class="top-header">
        <img src="https://kubelancerlogopublic.s3.us-east-1.amazonaws.com/Kubelancer+Logo-3.png" alt="Kubelancer Logo">
    </div>
    <div class="running-text">
        <p>AWS CLOUD DEMO LABS: SNS - Learn, Experiment, and Innovate with AWS</p>
    </div>

    <div class="container">
        <div class="form-section">
            <h4>How this AWS SNS Demo Lab Works</h4>
            <ol>
                <li>This lab demonstrates how to send messages using Amazon Simple Notification Service (SNS).</li>
                <li>Test event-driven messaging using the form below.</li>
                <li>Send a message to the SNS topic and see the response.</li>
            </ol>

            <form action="/publish" method="POST">
                <div class="form-group">
                    <label for="message">Enter Message:</label>
                    <input type="text" class="form-control" name="message" id="message" required>
                </div>
                <button type="submit" class="btn btn-primary">
                    ✉️ Send to SNS
                </button>
            </form>

            <div class="mt-3 alert alert-info" role="alert">
                {{.}}
            </div>
        </div>
    </div>

    <div class="footer">
        <p>© 2025 Kubelancer Private Limited | <a href="https://kubelancer.com">Visit Website</a></p>
        <p>📩 Contact us at <a href="mailto:connect@kubelancer.com">connect@kubelancer.com</a></p>
    </div>
</body>
</html>
`))

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/publish", publishHandler)
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
		return
	}

	topicArn := os.Getenv("SNS_TOPIC_ARN")
	if topicArn == "" {
		http.Error(w, "SNS_TOPIC_ARN not set in .env file", http.StatusInternalServerError)
		return
	}

	svc := sns.New(sess)
	result, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to publish message: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, fmt.Sprintf("✅ Message sent successfully! Message ID: %s", *result.MessageId))
}

// Instructions:
// 1. Create a `.env` file in the project root with the following content:
//    AWS_ACCESS_KEY_ID=your-access-key
//    AWS_SECRET_ACCESS_KEY=your-secret-key
//    AWS_REGION=us-east-1
//    SNS_TOPIC_ARN=arn:aws:sns:us-east-1:123456789012:YourTopicName
//
// 2. Initialize the Go module:
//    go mod init sns-demo
//    go get github.com/aws/aws-sdk-go
//    go get github.com/joho/godotenv
//
// 3. Run the application:
//    go run main.go
//
// 4. Open your browser and navigate to http://localhost:8080
//    - The logo now appears at the top header with a white background.
//    - The running text banner displays the message: "AWS CLOUD DEMO LABS: SNS - Learn, Experiment, and Innovate with AWS" continuously.
//    - Enter a message in the form and click "Send to SNS" to publish it to the SNS topic.
//    - A confirmation message will appear upon successful publishing.
