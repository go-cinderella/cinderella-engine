package model

import (
	"github.com/spf13/cast"
	"strings"
)

type MultiInstanceLoopCharacteristics struct {
	IsSequential        bool   `xml:"isSequential,attr"`
	Collection          string `xml:"collection,attr"`
	CompletionCondition string `xml:"completionCondition"`
}

func (receiver MultiInstanceLoopCharacteristics) String() string {
	var sb strings.Builder
	sb.WriteString("MultiInstanceLoopCharacteristics")
	sb.WriteString("{")
	sb.WriteString("IsSequential: ")
	sb.WriteString(cast.ToString(receiver.IsSequential))
	sb.WriteString(", ")
	sb.WriteString("Collection: ")
	sb.WriteString(receiver.Collection)
	sb.WriteString(", ")
	sb.WriteString("CompletionCondition: ")
	sb.WriteString(receiver.CompletionCondition)
	sb.WriteString("}")
	return sb.String()
}
