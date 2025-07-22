package srverr

type serverError struct {
	error   string
	details string

	servError Error
}

func NewServerError(servError Error) ServerError {
	return &serverError{
		servError: servError,
	}
}

func (s *serverError) Error() string {
	return s.error
}

func (s *serverError) SetError(err string) ServerError {
	s.error = err
	return s
}

func (s *serverError) SetDetails(details string) ServerError {
	s.details = details
	return s
}

func (s *serverError) GetServerError() Error {
	return s.servError
}

func (s *serverError) GetDetails() string {
	return s.details
}
