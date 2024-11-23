package utils

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelData struct {
	Sno              string
	Name             string
	Email            string
	Company          string
	ApplyingPosition string
	AdditionalInfo   string
	ReasonForContact string
}

func ReadExcelData(filePath *string) ([]ExcelData, error) {
	ff, err := excelize.OpenFile(*filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file at 1: %v", err)
	}
	defer ff.Close()

	rows, err := ff.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("error getting rows at 2: %v", err)
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("excel file is empty")
	}

	var data []ExcelData

	for i, row := range rows {
		if i == 0 {
			continue
		}
		r := ExcelData{
			Sno:              row[0],
			Name:             row[1],
			Email:            row[2],
			Company:          row[3],
			ApplyingPosition: row[4],
			AdditionalInfo:   row[5],
			ReasonForContact: row[6],
		}
		data = append(data, r)
	}

	return data, nil
}

func RemoveFirstAndLastLine(input string) string {
	input = strings.ReplaceAll(input, "```", "")
	lines := strings.Split(input, "\n")
	if len(lines) <= 2 {
		return ""
	}
	return strings.Join(lines[1:len(lines)-1], "\n")
}
