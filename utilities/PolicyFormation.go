package utilities

import (
	"exportcsv/models"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func PolicyFormation(outputfile *excelize.File, Panorama models.Device) string {

	rows, err := outputfile.Rows("Firewall_Access_Rules")

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
		}
		for i, colCell := range row {
			if i == 3 {

			}
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
	return ""
}
