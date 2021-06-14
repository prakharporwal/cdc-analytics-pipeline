package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"bfassignment/model"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func LambdaHandler(ctx context.Context, event *events.DynamoDBEvent) (string, error) {

	bucketName := "covidfiledumpbucket"

	for _, record := range event.Records {
		fmt.Println("size of records:: ", len(event.Records))
		sess := session.Must(session.NewSession())
		uploader := s3manager.NewUploader(sess)

		if record.EventName == "INSERT" || record.EventName == "MODIFY" {
			fmt.Println("EventName : ",record.EventName)
			u := record.Change.NewImage

			var item model.DataModel

			item.UUID = u["uuid"].String()
			item.Month = u["month"].String()
			item.Year = u["year"].String()
			item.ContinentExp = u["continentExp"].String()
			item.CountriesAndTerritories = u["countriesAndTerritories"].String()
			item.CountryterritoryCode = u["countryTerritoryCode"].String()
			item.GeoId = u["geoId"].String()
			item.CumulativeNumberFor14DaysofCOVID19CasesPer100000 = u["cumulative_number_for_14_days_of_COVID-19_cases_per_100000"].String()
			item.DateRep = u["dateRep"].String()

			item.Deaths, _ = strconv.Atoi(u["deaths"].Number())
			item.Cases, _ = strconv.Atoi(u["cases"].Number())
			item.PopData2019, _ = strconv.Atoi(u["popData2019"].Number())

			uuid := item.UUID
			dateRep := item.DateRep

			fmt.Println("Item : ",item)


			jsonBod, err := json.Marshal(item)

			if err != nil {
				log.Println(err.Error())
				return "can't create json response", errors.New("error: can't create json response")
			}

			key := fmt.Sprintf("input/%s/%s.json", dateRep, uuid)

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   bytes.NewReader(jsonBod),
			})

			if err != nil {
				log.Printf("There was an issue uploading to s3: %s", err.Error())
				return "Unable to save response", errors.New("error: cant save response")
			}
		}
	}

	return "Updated and uploaded to S3 : " + event.Records[0].EventName, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
