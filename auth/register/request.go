package register

type Request struct {
	Username string
	Email    string
	Password string

	Data any
}
