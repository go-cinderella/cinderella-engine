package model

type FieldExtension struct {
	FieldName   string `xml:"name,attr"`
	StringValue string `xml:"string"`
	Expression  string `xml:"expression"`
}
