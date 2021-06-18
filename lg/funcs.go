package lg

import (
	"fmt"
	"strings"
	"text/template"
)

var joinResults = func(results []MethodResult, formattedName bool) string {
	joined := make([]string, len(results))
	for i, r := range results {
		if formattedName {
			joined[i] = fmt.Sprintf("%s %s", ResultVar(i+1), r.Type) // e.g., r2 *dto.User
		} else {
			joined[i] = r.joinNameAndType() // e.g, *dto.User
		}
	}

	return strings.Join(joined, ",")
}

func ResultVar(index int) string {
	return fmt.Sprintf("r%d", index)
}

func Funcs() template.FuncMap {
	return template.FuncMap{
		"joinParamsWithType": func(params []MethodParam) string {
			var joined []string
			for _, p := range params {
				joined = append(joined, fmt.Sprintf("%s %s", p.Name, p.Type))
			}

			return strings.Join(joined, ",")
		},
		"joinParams": func(params []MethodParam) string {
			var joined []string
			for _, p := range params {
				joined = append(joined, fmt.Sprintf("%s", p.Name))
			}

			return strings.Join(joined, ",")
		},
		// Example for formatted name is : (r1 *dto.User,r2 err error)
		// Example for original name is: (*dto.User, error) or (u *dto.User,e error)
		"joinResultsForSignature": func(results []MethodResult, formattedName bool) string {
			if len(results) == 0 || (len(results) == 1 && results[0].Name == "") {
				return joinResults(results, formattedName)
			}
			return fmt.Sprintf("(%s)", joinResults(results, formattedName))
		},
		"genResultsVars": func(results []MethodResult) string {
			genList := make([]string, len(results))
			for i, _ := range results {
				genList[i] = ResultVar(i+1)
			}

			return strings.Join(genList, ",")
		},
		"hasErrInResults": func(results []MethodResult) bool {
			return len(results) != 0 && IsError(results[len(results)-1].Type)
		},
		"errResultVar": func(results []MethodResult) string {
			for i, r := range results {
				if IsError(r.Type) {
					return ResultVar(i)
				}
			}
			return ""
		},
		"title": func(val string) string {
			return strings.Title(val)
		},
	}
}
