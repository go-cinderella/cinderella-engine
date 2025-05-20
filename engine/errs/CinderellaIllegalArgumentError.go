package errs

type CinderellaIllegalArgumentError struct {
	CinderellaError
}

func NewCinderellaIllegalArgumentError(msg string) CinderellaIllegalArgumentError {
	return CinderellaIllegalArgumentError{CinderellaError: CinderellaError{Msg: msg}}
}
