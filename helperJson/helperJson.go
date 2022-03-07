package helperJson

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"wb_l0/database"
)

func ReadJsonFile(nameFile string) *database.Order {
	var m *database.Order
	jsonFile, err := os.Open(nameFile)
	defer jsonFile.Close()
	if err != nil {
		log.Println(err)
	}
	byteArray, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(byteArray, &m)

	if err != nil {
		log.Println(err)
	}

	//res, err := PrettyStruct(m)
	//
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Println(res)

	return m
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
