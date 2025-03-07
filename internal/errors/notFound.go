package errors

type NotFoundError struct {
	typ string
}

func CreateNotFoundError(typ string) NotFoundError {
	return NotFoundError{typ: typ}
}

func (nf NotFoundError) Error() string {
	return nf.typ + " not found."
}
