package exporters

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ExportPolicy(textscanner *bufio.Scanner, outputfile *excelize.File, counter *int) {
	textscanner.Scan()
	if strings.Contains(textscanner.Text(), "edit") {
		stdFields := []string{
			"Name",
			"Status",
			"Schedule",
			"Source Interface",
			"Source IP Address",
			"Source NAT",
			"Destination Interface",
			"Destination IP Address",
			"Action",
			"Service",
			"Applications",
		}
		newFields := []string{
			"Name",
			"Status",
			"Schedule",
			"Source Interface",
			"Source IP Address",
			"Source NAT",
			"Destination Interface",
			"Destination IP Address",
			"Action",
			"Service",
			"Applications",
		}
		for textscanner.Scan() {

			if strings.Contains(textscanner.Text(), "set") {
				command := strings.TrimSpace(textscanner.Text())

				switch {
				case strings.Contains(command, "set name") && newFields[0] == "Name":
					newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set name", ""))
				case strings.Contains(command, "set status") && newFields[1] == "Status":
					newFields[1] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set status", ""))
				case strings.Contains(command, "set schedule") && newFields[2] == "Schedule":
					newFields[2] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set schedule", ""))
				case strings.Contains(command, "set srcintf") && newFields[3] == "Source Interface":
					newFields[3] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set srcintf", ""))
				case strings.Contains(command, "set srcaddr") && newFields[4] == "Source IP Address":
					membersSlc := strings.Split(strings.TrimSpace(strings.ReplaceAll(command, "set srcaddr", "")), `" "`)
					newFields[4] = ""
					for i, val := range membersSlc {
						if len(membersSlc) != (i + 1) {
							newFields[4] = newFields[4] + strings.ReplaceAll(val, `"`, "") + ", \n"
						} else {
							newFields[4] = newFields[4] + strings.ReplaceAll(val, `"`, "")
						}
					}
				case strings.Contains(command, "set poolname") && newFields[5] == "Source NAT":
					newFields[5] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set poolname", ""))
				case strings.Contains(command, "set dstintf") && newFields[6] == "Destination Interface":
					newFields[6] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set dstintf", ""))
				case strings.Contains(command, "set dstaddr") && newFields[7] == "Destination IP Address":
					membersSlc := strings.Split(strings.TrimSpace(strings.ReplaceAll(command, "set dstaddr", "")), `" "`)
					newFields[7] = ""
					for i, val := range membersSlc {
						if len(membersSlc) != (i + 1) {
							newFields[7] = newFields[7] + strings.ReplaceAll(val, `"`, "") + ", \n"
						} else {
							newFields[7] = newFields[7] + strings.ReplaceAll(val, `"`, "")
						}
					}
				case strings.Contains(command, "set action") && newFields[8] == "Action":
					newFields[8] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set action", ""))
				case strings.Contains(command, "set service") && newFields[9] == "Service":
					membersSlc := strings.Split(strings.TrimSpace(strings.ReplaceAll(command, "set service", "")), `" "`)
					newFields[9] = ""
					for i, val := range membersSlc {
						if len(membersSlc) != (i + 1) {
							newFields[9] = newFields[9] + strings.ReplaceAll(val, `"`, "") + ", \n"
						} else {
							newFields[9] = newFields[9] + strings.ReplaceAll(val, `"`, "")
						}
					}
				case strings.Contains(command, "set application-list") && newFields[10] == "Applications":
					newFields[10] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set application-list", ""))
				}
			} else if strings.Contains(textscanner.Text(), "next") {
				*counter++
				for i, val := range newFields {
					if val == stdFields[i] {
						switch val {
						case "Action":
							newFields[i] = "deny"
						default:
							newFields[i] = ""
						}
					}
				}
				outputfile.SetSheetRow("Firewall_Access_Rules", "A"+strconv.Itoa(*counter), &newFields)

				newFields = []string{
					"Name",
					"Status",
					"Schedule",
					"Source Interface",
					"Source IP Address",
					"Source NAT",
					"Destination Interface",
					"Destination IP Address",
					"Action",
					"Service",
					"Applications",
				}

			} else if strings.Contains(textscanner.Text(), "edit") {
				continue
			} else {
				break
			}

		}
	}
}
