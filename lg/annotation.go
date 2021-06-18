package lg

import (
	"go/ast"
	"reflect"
	"regexp"
)

var annotationRegex = regexp.MustCompile(`@(\w+)(?:(?:\s+)\x60(.*)\x60)?`) // \x60 is backtick.

// Annotation is a type of comment on any goalng ndoe (struct, method, field,...) that has following format:
// @annotationName	`tagField:"tag val" anotherField:"another val"`
// e.g.,
// @tx	`retryCount:"4"`
type Annotation struct {
	Name string
	Tag  reflect.StructTag
}

type Annotations []Annotation

func (a Annotations) Lookup(name string) *Annotation {
	for _, v := range a {
		if v.Name == name {
			return &v
		}
	}

	return nil
}

func annotationsFromCommentGroup(l []*ast.Comment) Annotations {
	annotations := make([]Annotation, 0)
	for _, c := range l {
		if annotationRegex.Match([]byte(c.Text)) {
			res := annotationRegex.FindStringSubmatch(c.Text)
			if res[2] == "" {
				annotations = append(annotations, Annotation{Name: res[1]})
			} else {
				annotations = append(annotations, Annotation{Name: res[1], Tag: reflect.StructTag(res[2])})
			}
		}
	}
	return annotations
}
