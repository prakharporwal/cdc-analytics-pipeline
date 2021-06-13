package services

import (
	"bfassignment/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "reflect"
	"strings"

	"github.com/google/uuid"
)

func GetDataFromDisk() []model.DataModel {

	data, err := ioutil.ReadFile("./data/source.json")
	if err != nil {
		fmt.Print(err)
	}

	dataModelList, err := convertJsonToModel(data)

	if err != nil {
		fmt.Println("ERROR CONVERTING JSON TO MODEL", err)
	}

	dataModelList = assignUUID(dataModelList)

	return dataModelList
}

func convertJsonToModel(data []byte) ([]model.DataModel, error) {

	var result map[string][]interface{}

	err := json.Unmarshal([]byte(data), &result)

	if err != nil {
		fmt.Println("error:", err)
	}

	objList := result["records"]

	var dataModelList []model.DataModel

	finalList, err := json.Marshal(objList)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = json.Unmarshal(finalList, &dataModelList)
	if err != nil {
		fmt.Println("error:", err)
	}

	return dataModelList, err
}

func assignUUID(dataModelList []model.DataModel)([]model.DataModel) {
	fmt.Println("----Assigning uuid to each object----")

	for i := 0; i < len(dataModelList); i++ {
		uuidWithHyphen := uuid.New()
		uuidPlain := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
		dataModelList[i].UUID = uuidPlain
	}

	fmt.Println("------COMPLETED------")

	return dataModelList
}
