package login

type Request struct {
	Email string

	Password string

	// Optional device information

	DeviceID string

	DeviceName string

	Platform string

	IPAddress string

	UserAgent string
}
