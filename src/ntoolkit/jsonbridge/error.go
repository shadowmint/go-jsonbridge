package jsonbridge

// ErrRead indicates reading from the bridge stream failed.
type ErrRead struct{}

// ErrWrite indicates writing to the bridge stream failed.
type ErrWrite struct{}

// ErrMarshal indicates conversion to json failed.
type ErrMarshal struct{}

// ErrUnmarshal indicates conversion from json failed.
type ErrUnmarshal struct{}

// ErrNoData is raised when a request is made and no pending messages exist.
type ErrNoData struct{}
