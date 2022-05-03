package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ElectionData struct {
	Provience string `json:"Provience"`
	District  string `json:"District"`
	LocalBody string `json:"Local Body"`
	WardNo    string `json:"Ward No"`
	Post      string `json:"Post"`
	Candidate string `json:"Candidate"`
	Party     string `json:"Party"`
}

func ReadAndParseData() ([]ElectionData, error) {
	jsonFile, err := os.Open("./data/candidate_list.json")

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var electionData []ElectionData

	err = json.Unmarshal(byteValue, &electionData)

	if err != nil {
		return nil, err
	}

	return electionData, nil
}

func convertJSONToCSV(electionData []ElectionData, destination string) error {
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"Provience", "District", "Local Body", "Ward No", "Post", "Candidate", "Party"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, e := range electionData {
		var csvRow []string
		csvRow = append(csvRow, e.Provience, e.District, e.LocalBody, e.WardNo, e.Post, e.Candidate, e.Party)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	electionData, err := ReadAndParseData()

	if err != nil {
		fmt.Println(err)
	}

	var mainMap = make(map[string][]ElectionData)

	for _, data := range electionData {
		mainMap[fmt.Sprintf("%s__%s__%s__%s", data.Provience, data.District, data.LocalBody, data.Post)] = append(mainMap[fmt.Sprintf("%s__%s__%s__%s", data.Provience, data.District, data.LocalBody, data.Post)], data)
	}

	for key, value := range mainMap {
		all := strings.Split(key, "__")
		fileName := fmt.Sprintf("list/%s/%s/%s/", all[0], all[1], all[2])
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			err := os.MkdirAll(fileName, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
		csvFileName := fmt.Sprintf("%s%s.csv", fileName, all[3])
		os.Create(csvFileName)
		if err := convertJSONToCSV(value, csvFileName); err != nil {
			log.Fatal(err)
		}
		fmt.Println("fileName", fileName)
	}
}
