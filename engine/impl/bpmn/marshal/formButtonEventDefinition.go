package marshal

type FormButtonEventDefinition struct {
	CandidateGroups string `xml:"candidateGroups,attr,omitempty"`
	OpenForm        bool   `xml:"openForm,attr"`
	OpenConfirm     bool   `xml:"openConfirm,attr"`
	FormKey         string `xml:"formKey,attr,omitempty"`
}
