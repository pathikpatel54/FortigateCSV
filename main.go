package main

import (
	"bufio"
	"exportcsv/exporters"
	"exportcsv/utilities"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {

	fmt.Print("1. Export Fortigate Configurations\n2. Import Fortigate configuration with Mapped Zones to PaloAlto\nPlease select the operation : ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := scanner.Text()

	switch response {
	case "1":
		inputfile, outputfile, filename, path := utilities.FileReadWrite()

		defer inputfile.Close()

		textscanner := bufio.NewScanner(utilities.ConvertEncoding(inputfile))

		for textscanner.Scan() {
			counter := 1
			if strings.Contains(textscanner.Text(), "config firewall policy6") {
				continue
			} else if strings.Contains(textscanner.Text(), "config firewall policy") {
				exporters.ExportPolicy(textscanner, outputfile, &counter)
			} else if strings.Contains(textscanner.Text(), "config system interface") {
				exporters.ExportInterfaces(textscanner, outputfile, &counter)
			} else if strings.Contains(textscanner.Text(), "config firewall address6") {
				continue
			} else if strings.Contains(textscanner.Text(), "config firewall address") {
				exporters.ExportObjects(textscanner, outputfile, &counter)
			} else if strings.Contains(textscanner.Text(), "config firewall addrgrp") {
				exporters.ExportObjectGroups(textscanner, outputfile, &counter)
			} else if strings.Contains(textscanner.Text(), "config firewall service custom") {
				exporters.ExportServiceObjects(textscanner, outputfile, &counter)
			}
		}

		outputfile.SetSheetRow("Firewall_Interfaces", "K1", &[]interface{}{"Mapped Zone", "Interface Name", "Interface Type", "Parent/Aggregate Group Interface", "New IP Address"})

		dvRange := excelize.NewDataValidation(true)
		dvRange.Sqref = "M2:M51"
		dvRange.SetDropList([]string{"Layer2", "Layer3", "SubInterface", "Aggregate"})
		outputfile.AddDataValidation("Firewall_Interfaces", dvRange)

		if err := textscanner.Err(); err != nil {
			log.Fatal(err)
		}

		if err := outputfile.SaveAs(path + "/configs/" + strings.Split(filename, ".")[0] + ".xlsx"); err != nil {
			log.Println(err)
		}
		fmt.Print("Configuration has been exported successfully to file '" + strings.Split(filename, ".")[0] + ".xlsx'" + " Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	case "2":

		outputfile := utilities.FileOpen()
		Panorama := utilities.PolicyCreation()
		cmd := utilities.PolicyFormation(outputfile, Panorama)
		utilities.WriteConn(Panorama.Hostname, Panorama.Username, Panorama.Password, cmd)
	default:

	}

}
