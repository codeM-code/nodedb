package nodedb

// uses go text/template for primite string replacement

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

//SQLString provide changed text replacements to form a valid sql string
type SQLString string

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err) // TODO better error handling
	}
	return tmplBytes.String()
}

func processString(str string, vars interface{}) string {
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		panic(err) // TODO better error handling
	}
	return process(tmpl, vars)
}

func intSliceToString(ints []int64) string {
	s := make([]string, len(ints), len(ints))
	for i, n := range ints {
		s[i] = fmt.Sprint(n)
	}
	return strings.Join(s, ",")
}

// helper: replaces collection name template
func (sql SQLString) Collection(collectionname string) SQLString {

	values := baseMap()
	values["collection"] = collectionname
	return processSQLString(sql, values)
}

// helper: replaces commaContent, handles with or without content variations
func (sql SQLString) WithContent(withContent bool) SQLString {

	values := baseMap()
	values["commaContent"] = ""
	values["commaQuestionMark"] = ""
	if withContent {
		values["commaContent"] = ",content"
		values["commaQuestionMark"] = ",?"
	}
	return processSQLString(sql, values)
}

// helper: replaces nodeids, handles 'id in (...)' clauses
func (sql SQLString) NodeIDs(nodeIDs []int64) SQLString {

	values := baseMap()
	values["nodeids"] = intSliceToString(nodeIDs)
	return processSQLString(sql, values)
}

func (sql SQLString) String() string {
	return string(sql)
}

func baseMap() map[string]string { // replacement invariants

	result := make(map[string]string)
	result["collection"] = "{{.collection}}"
	result["commaContent"] = "{{.commaContent}}"
	result["commaQuestionMark"] = "{{.commaQuestionMark}}"
	result["nodeids"] = "{{.nodeids}}"

	return result
}

func processSQLString(sql SQLString, vars interface{}) SQLString {

	str := string(sql)
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		panic(err) // TODO better error handling
	}
	return SQLString(process(tmpl, vars))
}
