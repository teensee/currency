package dto

type Error struct {
	Message string `json:"errorMessage"`
}

func NewError(err error) Error {
	return Error{Message: err.Error()}
}
