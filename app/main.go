package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ehardi19/go-fuzzy-logic"
)

func main() {
	csvFile, _ := os.Open("DataTugas3.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	var data []fuzzy.Number
	reader.Read()
	for {
		var fn fuzzy.Number
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fn.Interview.ID = row[0]
		fn.Interview.Competence, err = strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		fn.Interview.Personality, err = strconv.ParseFloat(row[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, fn)
	}

	file, _ := os.Create("TebakanTugas3.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Fuzzy
	blt := fuzzy.BLT{}
	for i := range data {
		blt.Fuzzification(&data[i])
		blt.Inference(&data[i])
		blt.Defuzzification(&data[i])
	}

	// Inseting Data
	head := []string{
		"ID",
		"Competence",
		"Personality",
		"Diterima",
	}
	if err := writer.Write(head); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	for i := range data {
		csvData := []string{
			fmt.Sprintf("%s", data[i].Interview.ID),
			fmt.Sprintf("%.0f", data[i].Interview.Competence),
			fmt.Sprintf("%.0f", data[i].Interview.Personality),
			fmt.Sprintf("%f", data[i].CrispValue),
		}
		if err := writer.Write(csvData); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
}
