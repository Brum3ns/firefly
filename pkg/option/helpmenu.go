package option

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type groupField struct {
	name     string
	defValue string
	usage    string
}

func (opt *Options) customUsage() {
	m := make(map[string]string)
	lst_groupOrder := []string{}
	v := reflect.ValueOf(opt).Elem()
	t := v.Type()

	//Extract nested struct from the core struct (Aka: option group):
	for i := 0; i < v.NumField(); i++ {
		groupValue := v.FieldByName(t.Field(i).Name)
		groupType := groupValue.Type()

		lst_groupOrder = append(lst_groupOrder, groupType.Name())

		//Extract the variables (options/flags) inside the struct (options) found in the option struct (option group)
		for i := 0; i < groupValue.NumField(); i++ {
			if tag := string(groupType.Field(i).Tag.Get("flag")); len(tag) > 0 {
				m[tag] = groupType.Name()
			}
		}
	}

	//Add the groups options to the map "groupOption":
	groupOption := map[string][]groupField{}
	flag.VisitAll(func(f *flag.Flag) {
		groupName := m[f.Name]
		groupOption[groupName] = append(groupOption[groupName], groupField{
			name:     f.Name,
			usage:    f.Usage,
			defValue: f.DefValue,
		})
	})

	//Make help menu by adding the "group name", "option", "usage" and default value (if any):
	menu := make(map[string]string)
	space_width := "\t"
	for group_name, strt := range groupOption {
		for _, field := range strt {
			var defaultValue string

			if len(field.defValue) > 0 {

				defaultValue = fmt.Sprintf("(Default: %s%v\033[0m)", colorDefaultValue(field.defValue), field.defValue)
			}
			if len(field.name) <= 4 {
				space_width = "      \t"
			}
			menu[group_name] += fmt.Sprintf("  -%s%s %s\n", field.name, (space_width + field.usage), defaultValue)
		}
	}

	//Print the help menu:
	fmt.Println("Usage: firefly -u 'target.com/query=FUZZ' [OPTION] ...")
	for _, k := range lst_groupOrder {
		fmt.Printf("%s:\n%s\n", strings.ToUpper(k), menu[k])
	}
	exampleUsage()
}

func colorDefaultValue(s string) string {
	if s == "false" || s == "true" {
		return "\033[1;34m"
	} else if _, err := strconv.Atoi(s); err == nil {
		return "\033[1:32m"
	}
	return "\033[33m"
}

func exampleUsage() {
	fmt.Println(`[ Basic Examples ]
  firefly -u 'target.com/?query=FUZZ'
  firefly -u 'target.com/?query=hoodie&sort=DESC' -au replace -H 'Host:localhost'
  firefly -u 'target.com/?query=FUZZ&cachebuster=#RANDOM#' -e url -w wordlist.txt -pr '( ) => (/**/)'
  `)
}
