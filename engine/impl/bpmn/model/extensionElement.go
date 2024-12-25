package model

type ExtensionElement struct {
	TaskListener    []ActivitiListener `xml:"taskListener"`
	FieldExtensions []FieldExtension   `xml:"field"`
}

func (receiver ExtensionElement) GetFieldByName(fieldName string) FieldExtension {
	for _, element := range receiver.FieldExtensions {
		if element.FieldName == fieldName {
			return element
		}
	}
	return FieldExtension{}
}
