package gobioweb


type AppError struct {
	Error   error
	Message string
	Code    int
	Path string
}
