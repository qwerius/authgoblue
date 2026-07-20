package middleware

type Options struct {

	// Jika true, middleware hanya melakukan parsing token.
	// Validasi access token dilewati.
	// Berguna untuk kebutuhan khusus.
	SkipValidation bool

	// Jika true, middleware tidak melakukan pengecekan session.
	//
	// Cocok untuk mode JWT stateless:
	// - token valid
	// - signature valid
	// - expiration valid
	//
	// Tidak ada Redis/SessionStore lookup.
	SkipSessionCheck bool
}
