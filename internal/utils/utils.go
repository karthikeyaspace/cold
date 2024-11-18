package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelData struct {
	Name             string
	Email            string
	Company          string
	Position         string
	CompanyType      string
	ApplyingPosition string
	AdditionalInfo   string
	MutualInterests  string
	ReasonForContact string
	Industry         string
}

func ReadExcelData(filePath string) ([]ExcelData, error) {
	ff, err := excelize.OpenFile(filePath)
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
		if len(row) < 7 {
			return nil, fmt.Errorf("excel file is invalid")
		}
		r := ExcelData{
			Name:             row[0],
			Email:            row[1],
			Company:          row[2],
			Position:         row[3],
			CompanyType:      row[4],
			ApplyingPosition: row[5],
			AdditionalInfo:   row[6],
			MutualInterests:  row[7],
			ReasonForContact: row[8],
			Industry:         row[9],
		}
		data = append(data, r)
	}

	return data, nil

}
