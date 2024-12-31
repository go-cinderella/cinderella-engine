package marshal

type DefaultBaseElement struct {
	Id   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,attr,omitempty"`
}
