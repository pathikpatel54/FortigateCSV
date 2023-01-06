package exporters

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ExportServiceObjects(textscanner *bufio.Scanner, outputfile *excelize.File, counter *int) {
	textscanner.Scan()
	if strings.Contains(textscanner.Text(), "edit") {
		stdFields := []string{
			"Name",
			"Category",
			"Protocol",
			"Protocol Number",
		}
		newFields := make([]string, 4)
		copy(newFields, stdFields)
		newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))

		for textscanner.Scan() {

			if strings.Contains(textscanner.Text(), "set") {
				command := strings.TrimSpace(textscanner.Text())

				switch {
				case strings.Contains(command, "set category") && newFields[1] == "Category":
					newFields[1] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set category", ""))
				case strings.Contains(command, "set protocol") && newFields[2] == "Protocol":
					newFields[2] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set protocol", ""))
				case strings.Contains(command, "set protocol-number") && newFields[3] == "Protocol Number":
					newFields[3] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set protocol-number", ""))
				case strings.Contains(command, "set tcp-portrange") && newFields[3] == "Protocol Number":
					membersSlc := strings.Split(strings.TrimSpace(strings.ReplaceAll(command, "set tcp-portrange", "")), ` `)
					newFields[3] = ""
					for i, val := range membersSlc {
						if len(membersSlc) != (i + 1) {
							newFields[3] = newFields[3] + strings.ReplaceAll(val, `"`, "") + ", \n"
						} else {
							newFields[3] = newFields[3] + strings.ReplaceAll(val, `"`, "")
						}
					}
					newFields[2] = "TCP"
				case strings.Contains(command, "set udp-portrange") && newFields[3] == "Protocol Number":
					membersSlc := strings.Split(strings.TrimSpace(strings.ReplaceAll(command, "set udp-portrange", "")), ` `)
					newFields[3] = ""
					for i, val := range membersSlc {
						if len(membersSlc) != (i + 1) {
							newFields[3] = newFields[3] + strings.ReplaceAll(val, `"`, "") + ", \n"
						} else {
							newFields[3] = newFields[3] + strings.ReplaceAll(val, `"`, "")
						}
					}
					newFields[2] = "UDP"
				}
			} else if strings.Contains(textscanner.Text(), "next") {
				*counter++
				for i, val := range newFields {
					if val == stdFields[i] {
						switch val {
						default:
							newFields[i] = ""
						}
					}
				}
				outputfile.SetSheetRow("Firewall_ServiceObjects", "A"+strconv.Itoa(*counter), &newFields)

				newFields = []string{
					"Name",
					"Category",
					"Protocol",
					"Protocol Number",
				}
			} else if strings.Contains(textscanner.Text(), "edit") {
				newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))
			} else if strings.Contains(textscanner.Text(), "end") && newFields[0] == "Name" {
				break
			}

		}
	}
}
