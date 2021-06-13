package jobs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func Read10000Items(conn *dynamodb.DynamoDB, tableName string) {
	fmt.Println("-----------sequentially reading from DDB table :"+tableName+ " -------")
	Last
	scanInputData := &dynamodb.ScanInput{
		TableName: &tableName,
		FilterExpression: ,
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

func updateTableItem(conn *dynamodb.DynamoDB,tableName string,uuid string, oldDate string) error {
	fmt.Println("----------updating Item-----------")
	newRandomDate := generateNewDateRep(oldDate);

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

	_, err := conn.UpdateItem(input)

	if err != nil {
		fmt.Println("Error Updating Item In DDB:",err)
	}

	return err
}


func generateNewDateRep(oldDate string) string {
	return "namaste"+oldDate;
}