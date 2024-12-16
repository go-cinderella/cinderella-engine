package support

import (
	"github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

type BooleanTypedValue struct {
	*spel.TypedValue
}

func (b *BooleanTypedValue) ForValue(bool2 bool) spel.TypedValue {
	boo := spel.TypedValue{}
	if bool2 {
		boo.Value = true
		return boo
	}

	boo.Value = false
	return boo
}
