package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"strings"
)

var (
	command      = kingpin.Arg("command", "Command").Required().String()
	verbose      = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	site         = kingpin.Flag("site", "Site").String()
	tenant       = kingpin.Flag("tenant", "Tenant").String()
	role         = kingpin.Flag("Role", "role").String()
	status       = kingpin.Flag("status", "Status").String()
	manufacturer = kingpin.Flag("manufacturer", "Vendor").String()
	customfield  = CustomFields(kingpin.Flag("customfield", "Custom Field definition as key-value pair IE: core=something"))
)

type customField struct {
	Key   string
	Value string
}
type customFields []customField

func (f *customFields) Set(value string) error {
	parsedField := strings.Split(value, "=")
	if len(parsedField) != 2 {
		return fmt.Errorf("Invalid custom field: %s", value)
	}
	*f = append(*f, customField{parsedField[0], parsedField[1]})
	return nil
}

func (f *customFields) String() string {
	return ""
}

func (f *customFields) IsCumulative() bool {
	return true
}

func CustomFields(s kingpin.Settings) (target *[]customField) {
	target = new([]customField)
	s.SetValue((*customFields)(target))
	return
}

func main() {
	kingpin.Parse()
	fmt.Printf("%v, %s\n", *verbose, *command)
	fmt.Println(customfield)
}
