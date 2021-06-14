package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	// "github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

var clusterIdentifier string = "redshift-cluster-1"

func queryRedShift(connRS *redshiftdataapiservice.RedshiftDataAPIService) {
	input := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(clusterIdentifier),
		DbUser:            aws.String("awsuser"),
		Database:          aws.String("dev"),
		Sql:               aws.String("select * from s3data_schema.covidfiles;"),
		// SecretArn:         aws.String("arn:aws:secretsmanager:us-east-2:737216973625:secret:redshift/DataAPI-secret-FZnjh6"),
	}
	out,_ := json.Marshal(input)
	fmt.Println(string(out))

	exec, err := connRS.ExecuteStatement(input)

	if err != nil {
		fmt.Println("Error Executing sql on Redshift :", err)
	}



	for true {
		res, _ := connRS.DescribeStatement(&redshiftdataapiservice.DescribeStatementInput{
			Id: aws.String(*exec.Id),
		})

		if *res.Status =="FINISHED" {
			fmt.Println("Query Status : ", *res.Status)
			break;
		}else{
			fmt.Println("Query Status - ", *res.Status)
			time.Sleep(2000)
		}
	}


	fmt.Println(exec, *exec.Id)

	results, err := connRS.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
		Id: aws.String(*exec.Id),
	})

	if err != nil {
		fmt.Println("Error getting query result on Redshift :", err)
	}

	fmt.Print(results.Records)

	for _, record := range results.Records{
		fmt.Println("Record :",record);
	}

	
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	connRS := redshiftdataapiservice.New(sess)

	queryRedShift(connRS)
}

// aws redshift-data execute-statement
//     --region us-east-2
//     --secret arn:aws:secretsmanager:us-east-2:737216973625:secret:redshift/DataAPI-secret-FZnjh6
//     --cluster-identifier redshift-cluster-1
//     --sql "select * from s3data_schema.covidfiles"
//     --database dev
