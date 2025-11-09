package proto

// Request represents a generic transfer request (placeholder).
type Request struct {
	Path     string
	Upload   bool
	Offset   int64
	Metadata map[string]string
}

// Response represents a generic response (placeholder).
type Response struct {
	OK      bool
	Message string
}
