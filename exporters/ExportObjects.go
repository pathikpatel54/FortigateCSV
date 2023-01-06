package exporters

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ExportObjectGroups(textscanner *bufio.Scanner, outputfile *excelize.File, counter *int) {
	textscanner.Scan()
	if strings.Contains(textscanner.Text(), "edit") {
		stdFields := []string{
			"Name",
			"Members",
			"Description",
		}
		newFields := make([]string, 3)
		copy(newFields, stdFields)
		newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))

		for textscanner.Scan() {

			if strings.Contains(textscanner.Text(), "set") {
				command := strings.TrimSpace(textscanner.Text())

				switch {
				case strings.Contains(command, "set member") && newFields[1] == "Members":
					membersSlc := strings.Split(strings.TrimSpace(strings.ReplaceAll(command, "set member", "")), `" "`)
					newFields[1] = ""
					for i, val := range membersSlc {
						if len(membersSlc) != (i + 1) {
							newFields[1] = newFields[1] + strings.ReplaceAll(val, `"`, "") + ", "
						} else {
							newFields[1] = newFields[1] + strings.ReplaceAll(val, `"`, "")
						}
					}
				case strings.Contains(command, "set comment") && newFields[2] == "Description":
					newFields[2] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set comment", ""))
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
				outputfile.SetSheetRow("Firewall_ObjectGroups", "A"+strconv.Itoa(*counter), &newFields)

				newFields = []string{
					"Name",
					"Members",
					"Description",
				}
			} else if strings.Contains(textscanner.Text(), "edit") {
				newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))
			} else if strings.Contains(textscanner.Text(), "end") && newFields[0] == "Name" {
				break
			}

		}
	}
}

func ExportObjects(textscanner *bufio.Scanner, outputfile *excelize.File, counter *int) {
	textscanner.Scan()
	if strings.Contains(textscanner.Text(), "edit") {
		stdFields := []string{
			"Name",
			"Type",
			"Value",
			"Description",
		}
		newFields := make([]string, 4)
		copy(newFields, stdFields)
		newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))
		setRange := false
		for textscanner.Scan() {

			if strings.Contains(textscanner.Text(), "set") {
				command := strings.TrimSpace(textscanner.Text())

				switch {
				case strings.Contains(command, "set type") && newFields[1] == "Type":
					newFields[1] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set type", ""))
				case strings.Contains(command, "set subnet") && newFields[2] == "Value":
					newFields[2] = strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set subnet", "")), " ")[0] + "/" +
						strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set subnet", "")), " ")[1]
				case strings.Contains(command, "set fqdn") && newFields[2] == "Value":
					newFields[2] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set fqdn", ""))
				case strings.Contains(command, "set start-ip") && newFields[2] == "Value":
					newFields[2] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set start-ip", "")) + " - "
					setRange = true
				case strings.Contains(command, "set end-ip") && setRange:
					newFields[2] = newFields[2] + strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set end-ip", ""))
				case strings.Contains(command, "set comment") && newFields[3] == "Description":
					newFields[3] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(command, "\"", ""), "set comment", ""))
				}
			} else if strings.Contains(textscanner.Text(), "next") {
				*counter++
				for i, val := range newFields {
					if val == stdFields[i] {
						switch val {
						case "Type":
							newFields[i] = "ipmask"
						default:
							newFields[i] = ""
						}
					}
				}
				outputfile.SetSheetRow("Firewall_Objects", "A"+strconv.Itoa(*counter), &newFields)

				newFields = []string{
					"Name",
					"Type",
					"Value",
					"Description",
				}
				setRange = false
			} else if strings.Contains(textscanner.Text(), "edit") {
				newFields[0] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(textscanner.Text(), "\"", ""), "edit", ""))
			} else if strings.Contains(textscanner.Text(), "end") && newFields[0] == "Name" {
				break
			}

		}
	}
}
