package exception

type DataNotFoundError struct {
	Message string
}

func (dataNotFoundError DataNotFoundError) Error() string {
	return dataNotFoundError.Message
}
