package common

type CustomError struct {
	Src string
	Grp string
	Msg string
}

func (err CustomError) Error() string {
	return err.Msg
}

func (err CustomError) Unwrap() error {
	return err
}
