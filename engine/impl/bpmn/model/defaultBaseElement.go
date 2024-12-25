package model

type DefaultBaseElement struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

func (d DefaultBaseElement) GetId() string {
	return d.Id
}

func (d DefaultBaseElement) GetName() string {
	return d.Name
}
