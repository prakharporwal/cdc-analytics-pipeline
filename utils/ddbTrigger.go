package main

import (
	"fmt"
        "context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
        // "github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// func LambdaHandler(ctx context.Context, event events.DynamoDBEvent) (string, error) {
func LambdaHandler(ctx context.Context, event *events.DynamoDBEvent) (string, error) {

	// svc := s3.New(session.New())

	for _, record := range event.Records {
                fmt.Println(record);
		// jsonBod, err := json.Marshal(record)
		// if err != nil {
		// 	log.Println(err.Error())
		// 	return "can't create json response", errors.New("Error: can't create json response")
		// }

		// sess := session.Must(session.NewSession())
		// uploader := s3manager.NewUploader(sess)
		// u := uuid.New()
		// key := fmt.Sprintf("responses/%s", u.String())
		// _, ierr := uploader.Upload(&s3manager.UploadInput{
		// 	Bucket: aws.String("hexaco"),
		// 	Key:    aws.String(key),
		// 	Body:   bytes.NewReader(jsonBod),
		// })

		// if ierr != nil {
		// 	log.Printf("There was an issue uploading to s3: %s", ierr.Error())
		// 	return "Unable to save response", errors.New("Error: cant save response")
		// }

		// return u.String(), nil

	}

	return "Updated " + event.Records[0].EventName , nil;
}

func main() {
	lambda.Start(LambdaHandler)
}
