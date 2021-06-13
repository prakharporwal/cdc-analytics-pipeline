package main

import (
	"bfassignment/services"
	"bfassignment/jobs"
	"fmt"
)

func main() {
	fmt.Println("..............Starting your Program....................")

	tableName := "covid"

	conn := services.GetDBInstance()

	// fmt.Println(conn);
	// fmt.Println(services.GetDBInstance());

	// data :=  services.GetDataFromDisk()
	// _ = services.GetDataFromDisk()

	// services.AddItemsinDDB(conn, tableName, data)

	jobs.Read10000Items(conn,tableName)
}
