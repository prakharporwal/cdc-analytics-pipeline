package services

import (
	"fmt"
	// "log"
	"os"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

import (
	"bfassignment/model"
)

var ddbInstance *dynamodb.DynamoDB = nil

func dbisEmpty() bool {
	return ddbInstance == nil
}

func initDB() {

}

func GetDBInstance() *dynamodb.DynamoDB {
	if ddbInstance == nil {
		// Initialize a session that the SDK will use to load
		// credentials from the shared credentials file ~/.aws/credentials
		// and region from the shared configuration file ~/.aws/config.
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		// Create DynamoDB client connection
		svc := dynamodb.New(sess)

		ddbInstance = svc
	}
	return ddbInstance
}

func AddItemsinDDB(svc *dynamodb.DynamoDB,tableName string, data []model.DataModel) {

	for i,item := range data{
		
		if i==5{
			break;
		}

		av, err := dynamodbattribute.MarshalMap(item)

		// log.Printf("ROW: %s", av)


		// if i== 1000{
			fmt.Println("UUID :",item.UUID)
		// }

		if err != nil {
			fmt.Println("Got error marshalling new cupcake entry:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		inputItem := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		_=inputItem
		// _, err = svc.PutItem(inputItem)

		if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

	}
	fmt.Println("Successfully added Data to " + tableName)
}
