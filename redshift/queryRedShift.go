// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"

// 	// "github.com/aws/aws-sdk-go/service/redshift"
// 	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
// )

// var clusterIdentifier string = "redshift-cluster-1"

// func queryRedShift(connRS *redshiftdataapiservice.RedshiftDataAPIService) {
// 	input := &redshiftdataapiservice.ExecuteStatementInput{
// 		ClusterIdentifier: aws.String(clusterIdentifier),
// 		DbUser:            aws.String("awsuser"),
// 		Database:          aws.String("dev"),
// 		Sql:               aws.String("select * from s3data_schema.covidfile;"),
// 		// SecretArn:         aws.String("arn:aws:secretsmanager:us-east-2:737216973625:secret:redshift/DataAPI-secret-FZnjh6"),
// 	}
// 	out,_ := json.Marshal(input)
// 	fmt.Println(string(out))

// 	exec, err := connRS.ExecuteStatement(input)

// 	if err != nil {
// 		fmt.Println("Error Executing sql on Redshift :", err)
// 	}

// 	for true {
// 		res, _ := connRS.DescribeStatement(&redshiftdataapiservice.DescribeStatementInput{
// 			Id: aws.String(*exec.Id),
// 		})

// 		if *res.Status =="FINISHED" {
// 			fmt.Println("Query Status : ", *res.Status)
// 			break;
// 		}else{
// 			fmt.Println("Query Status - ", *res.Status)
// 			time.Sleep(2000)
// 		}
// 	}

// 	fmt.Println(exec, *exec.Id)

// 	results, err := connRS.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
// 		Id: aws.String(*exec.Id),
// 	})

// 	if err != nil {
// 		fmt.Println("Error getting query result on Redshift :", err)
// 	}

// 	// fmt.Print(results.Records)

// 	for _, record := range results.Records{
// 		for _,item := range record {
// 			fmt.Print (*item.StringValue, " ");
// 		}
// 		fmt.Println();
// 	}

// }

// func main() {
// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	connRS := redshiftdataapiservice.New(sess)

// 	queryRedShift(connRS)
// }

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

var clusterIdentifier string = "redshift-cluster-1"

func runQuery(connRS *redshiftdataapiservice.RedshiftDataAPIService, query string) (string, error) {

	input := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(clusterIdentifier),
		DbUser:            aws.String("awsuser"),
		Database:          aws.String("dev"),
		Sql:               aws.String(query),
	}

	out, _ := json.Marshal(input)
	fmt.Println(string(out))

	exec, err := connRS.ExecuteStatement(input)

	if err != nil {
		// fmt.Println()
		return "Error Executing sql on Redshift : " + err.Error(), err
	}

	var res *redshiftdataapiservice.DescribeStatementOutput

	for i := 0; i < 50; i++ {
		res, _ = connRS.DescribeStatement(&redshiftdataapiservice.DescribeStatementInput{
			Id: aws.String(*exec.Id),
		})

		if *res.Status == "FINISHED" {
			fmt.Println("Query Status : ", *res.Status)
			break
		} else {
			fmt.Println("Query Status - ", *res.Status)
			time.Sleep(2000)
		}
	}

	if *res.Status != "FINISHED" {
		return "I am Tired ! Check for query results in some time future, ID: " + *exec.Id, nil
	}

	results := getResults(connRS, *exec.Id)

	return results, nil
}

func getResults(connRS *redshiftdataapiservice.RedshiftDataAPIService, queryId string) string {
	results, _ := connRS.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
		Id: aws.String(queryId),
	})

	var output string

	for _, record := range results.Records {
		for _, item := range record {
			output += *item.StringValue + ` `
		}
		output += `\n`
	}

	return output
}

func queryRedShift(query string) (string, error) {
	// func queryRedShift() string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
	 		}))

	connRS := redshiftdataapiservice.New(sess)

	// if query == "get-result" {
	// 	return getResults(connRS, event.RawPath), nil
	// }

	results, err := runQuery(connRS, query)
	if err != nil {
		return results, err
	}
	return results, nil
}

func main() {
	query := "select * from s3data_schema.covidfile;"

	fmt.Println(queryRedShift(query))

}
