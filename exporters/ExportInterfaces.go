package exporters

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ExportInterfaces(textscanner *bufio.Scanner, outputfile *excelize.File, counter *int) {
	textscanner.Scan()
	if strings.Contains(textscanner.Text(), "edit") {
		stdFields := []string{
			"Name",
			"Description",
			"Status",
			"IP Address",
			"Subnet Mask",
			"Allowed Access",
			"Type",
			"Alias",
			"Role",
			"Secondary IP",
		}
		newFields := make([]string, 9)
		copy(newFields, stdFields)

		newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))
		for textscanner.Scan() {
			if strings.Contains(textscanner.Text(), "set") || strings.Contains(textscanner.Text(), "config secondaryip") {
				command := strings.TrimSpace(textscanner.Text())
				switch {
				case strings.Contains(command, "set description") && newFields[1] == "Description":
					newFields[1] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set description", ""))
				case strings.Contains(command, "set status") && newFields[2] == "Status":
					newFields[2] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set status", ""))
				case strings.Contains(command, "set ip") && newFields[3] == "IP Address" && newFields[4] == "Subnet Mask":
					newFields[3] = strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set ip", "")), " ")[0]
					newFields[4] = strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set ip", "")), " ")[1]
				case strings.Contains(command, "set allowaccess") && newFields[5] == "Allowed Access":
					newFields[5] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set allowaccess", ""))
				case strings.Contains(command, "set type") && newFields[6] == "Type":
					newFields[6] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set type", ""))
				case strings.Contains(command, "set alias") && newFields[7] == "Alias":
					newFields[7] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set alias", ""))
				case strings.Contains(command, "set role") && newFields[8] == "Role":
					newFields[8] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set role", ""))
				case strings.Contains(command, "config secondaryip") && newFields[9] == "Secondary IP":
					textscanner.Scan()
					if strings.Contains(textscanner.Text(), "edit") {
						newFields[9] = ""
						for textscanner.Scan() {
							if strings.Contains(textscanner.Text(), "set") {
								if strings.Contains(textscanner.Text(), "set allowaccess") {
									continue
								}
								newFields[9] = newFields[9] + strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "set ip", "")) + ", \n"
							} else if strings.Contains(textscanner.Text(), "next") {
								continue
							} else if strings.Contains(textscanner.Text(), "end") {
								break
							}

						}
					}

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
				outputfile.SetSheetRow("Firewall_Interfaces", "A"+strconv.Itoa(*counter), &newFields)

				newFields = []string{
					"Name",
					"Description",
					"Status",
					"IP Address",
					"Subnet Mask",
					"Allowed Access",
					"Type",
					"Alias",
					"Role",
					"Secondary IP",
				}
			} else if strings.Contains(textscanner.Text(), "edit") {
				newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))
			} else if strings.Contains(textscanner.Text(), "end") && newFields[0] == "Name" {
				break
			}
		}

	}

}
