package jobs

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func Read10000Items(conn *dynamodb.DynamoDB, tableName string) {
	fmt.Println("-----------sequentially reading from DDB table :"+tableName+ " -------")

	scanInputData := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	scanOutputData, err := conn.Scan(scanInputData)

	if err!=nil{
		fmt.Println("ERROR Reading data from DDB:", err)
	}

	for i, item := range scanOutputData.Items {
		if i==2{
			break
		}

		var itemUUID string = *item["uuid"].S;
		var itemDate string = *item["dateRep"].S;
		fmt.Println(itemDate)

		updateTableItem(conn, tableName, itemUUID, itemDate)
	}

	fmt.Println("-----------COMPLETED : reading from DDB table :"+tableName+ " -------")
}

func updateTableItem(conn *dynamodb.DynamoDB,tableName string,uuid string, oldDate string){
	fmt.Println("----------updating Item-----------")
	newRandomDate := generateNewDateRep();

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {
				S: aws.String(uuid),
			},
		},

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newDateRep": {
				S: aws.String(newRandomDate),
			},
		},

		UpdateExpression: aws.String("set dateRep = :newDateRep"),
	}

	// _ =input
	_, err := conn.UpdateItem(input)

	if err != nil {
		fmt.Println("Error Updating Item In DDB:",err)
	}

}


func generateNewDateRep() string {
		rand.Seed(time.Now().UTC().UnixNano())

		min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		delta := max - min

		sec := rand.Int63n(delta) + min
		randomDate := time.Unix(sec, 0)

		randomDateString := randomDate.Format("02/01/2006")

		// fmt.Println(randDate)
		fmt.Println(randomDateString)

		return randomDateString
	}


