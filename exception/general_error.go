package exception

type GeneralError struct {
	Message string
}

func (generalError GeneralError) Error() string {
	return generalError.Message
}
