package utilities

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func FileReadWrite() (*os.File, *excelize.File, string, string) {
	fmt.Print("Enter name of Fortigate Config File: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filename := scanner.Text()

	path, err := os.Getwd()

	if err != nil {
		log.Panicln(err)
	}

	inputfile, err := os.Open(path + "/configs/" + filename)

	if err != nil {
		log.Panicln(err)
	}

	outputfile := excelize.NewFile()
	outputfile.SetSheetName("Sheet1", "Firewall_Access_Rules")
	outputfile.NewSheet("Firewall_Interfaces")
	outputfile.NewSheet("Firewall_Objects")
	outputfile.NewSheet("Firewall_ObjectGroups")
	outputfile.NewSheet("Firewall_ServiceObjects")
	outputfile.SetActiveSheet(outputfile.GetSheetIndex("Firewall_Access_Rules"))
	outputfile.SetSheetRow("Firewall_Access_Rules", "A1", &[]interface{}{"Name", "Status", "Schedule", "Source Interface", "Source IP Address", "Source NAT", "Destination Interface", "Destination IP Address", "Action", "Service", "Applications"})
	outputfile.SetSheetRow("Firewall_Interfaces", "A1", &[]interface{}{"Name",
		"Description",
		"Status",
		"IP Address",
		"Subnet Mask",
		"Allowed Access",
		"Type",
		"Alias",
		"Role",
		"Secondary IP"})

	outputfile.SetSheetRow("Firewall_Objects", "A1", &[]interface{}{"Name",
		"Type",
		"Value",
		"Description",
	})

	outputfile.SetSheetRow("Firewall_ObjectGroups", "A1", &[]interface{}{"Name",
		"Members",
		"Description",
	})

	outputfile.SetSheetRow("Firewall_ServiceObjects", "A1", &[]interface{}{
		"Name",
		"Category",
		"Protocol",
		"Protocol Number",
	})

	return inputfile, outputfile, filename, path
}

func FileOpen() *excelize.File {
	fmt.Print("Enter name of Exported config File (Make sure to have Columns added for Interface mapping): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filename := scanner.Text()

	path, err := os.Getwd()

	if err != nil {
		log.Panicln(err)
	}

	inputfile, err := excelize.OpenFile(path + "/configs/" + filename)

	if err != nil {
		log.Panicln(err)
	}

	return inputfile
}
