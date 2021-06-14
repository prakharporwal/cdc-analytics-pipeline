package main

import (
	"bytes"
	"context"
	"time"

	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var bucketName string = "covidfiledumpbucket"

var prefix string = "input"

var MERGED_FOLDER_NAME string = "merged/"

var MERGE_JSON_SEPARATOR string = ","

var FILE_EXTENSION string = ".json"

func getListOfFiles(connS3 *s3.S3) (objectList *s3.ListObjectsOutput) {

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	resp, err := connS3.ListObjects(params)
	if err != nil {
		log.Fatal("Unable to list items in bucket : ", bucketName, " error : ", err)
	}

	return resp
}

func mergeFiles(connS3 *s3.S3, objectList *s3.ListObjectsOutput) string {

	var mergedJSONstring string

	for _, key := range objectList.Contents {
		fmt.Println(*key.Key)

		numBytes, err := connS3.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(*key.Key),
			})

		if err != nil {
			log.Fatalf("Unable to download item %q, %v", *key.Key, err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(numBytes.Body)
		myFileContentAsString := buf.String()

		mergedJSONstring += myFileContentAsString + MERGE_JSON_SEPARATOR

		// fmt.Println("Downloaded :", myFileContentAsString, "\n\n\n")

	}

	return mergedJSONstring

}

func uploadtoS3(sess *session.Session, mergedJSONstring string) {

	currTime := time.Now().Format("2006-01-02-15-04")

	mergedObjectName := MERGED_FOLDER_NAME + currTime + FILE_EXTENSION

	fmt.Println(mergedObjectName)

	uploader := s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(mergedObjectName),
		Body:   bytes.NewReader([]byte(mergedJSONstring)),
	})

	if err != nil {
		fmt.Println("Error Uploading Merged Object :", err)
	}
}

func LambdaHandler(ctx context.Context, event events.S3Event) {

	sess, _ := session.NewSession()

	connS3 := s3.New(sess)

	for _, record := range event.Records {

		s3 := record.S3

		objectKey := s3.Object.Key
		bucketName := s3.Bucket.Name

		fmt.Println("Object Added:", objectKey, " in ", bucketName)
	}

	objectList := getListOfFiles(connS3)
	mergedJSONstring := mergeFiles(connS3, objectList)

	uploadtoS3(sess, mergedJSONstring)

}

func main() {
	lambda.Start(LambdaHandler)
	// currTime := time.Now().Format("2006-01-02-15-04")
	// print(currTime)

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	// connS3 := s3.New(sess)

	// objectList := getListOfFiles(connS3)

	// mergedJSONstring := mergeFiles(connS3, objectList)

	// uploadtoS3(sess, mergedJSONstring)
}
