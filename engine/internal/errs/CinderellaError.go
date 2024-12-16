package errs

import "github.com/unionj-cloud/toolkit/stringutils"

type CinderellaError struct {
	Code string
	Msg  string
}

func (error CinderellaError) Error() string {
	if stringutils.IsNotEmpty(error.Code) {
		return error.Code + "---" + error.Msg
	}
	return error.Msg
}
