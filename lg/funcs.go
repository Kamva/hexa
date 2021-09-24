package lg

import (
	"fmt"
	"strings"
	"text/template"
)

var joinResults = func(results []MethodResult) string {
	joined := make([]string, len(results))
	for i, r := range results {
		joined[i] = r.joinNameAndType()
	}

	return strings.Join(joined, ",")
}

var joinResultsInSameNameFormat = func(results []MethodResult) string {
	joined := make([]string, len(results))
	for i, r := range results {
		joined[i] = fmt.Sprintf("%s %s", ResultVar(i), r.Type)
	}

	return strings.Join(joined, ",")
}

func ResultVar(index int) string {
	return fmt.Sprintf("r%d", index+1)
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
		"joinParams": func(params []MethodParam) string { // e.g., // a,b,c
			var joined []string
			for _, p := range params {
				joined = append(joined, fmt.Sprintf("%s", p.Name))
			}

			return strings.Join(joined, ",")
		},
		// join params with unpack operator.
		"joinParamsWithUnpack": func(params []MethodParam) string { // e.g, a,b,c...
			var joined []string
			for _, p := range params {
				mp := fmt.Sprintf("%s", p.Name)
				if strings.Index(p.Type, "...") == 0 {
					mp += "..."
				}
				joined = append(joined, mp)
			}

			return strings.Join(joined, ",")
		},
		// Example for original name is: (*dto.User, error) or (u *dto.User,e error)
		"joinResultsForSignature": func(results []MethodResult) string {
			if len(results) == 0 || (len(results) == 1 && results[0].Name == "") {
				return joinResults(results)
			}
			return fmt.Sprintf("(%s)", joinResults(results))
		},
		// Example for formatted name is : (r1 *dto.User,r2 err error)
		"joinResultsForSignatureInSameNameFormat": func(results []MethodResult) string {
			if len(results) == 0 {
				return ""
			}

			return fmt.Sprintf("(%s)", joinResultsInSameNameFormat(results))
		},
		"genResultsVars": func(results []MethodResult) string {
			genList := make([]string, len(results))
			for i, _ := range results {
				genList[i] = ResultVar(i)
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
		"hasAnnotation": func(annotations Annotations, name string) bool {
			return annotations.Lookup(name) != nil
		},
	}
}
