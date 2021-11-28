package routing

type ResponseStatus struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

func ErrorRenderer(err error, statusText string, statusCode int) *ResponseStatus {
	return &ResponseStatus{
		Err:        err,
		StatusCode: statusCode,
		StatusText: statusText,
		Message:    err.Error(),
	}
}

func ErrorRendererDefault(err error) *ResponseStatus {
	return &ResponseStatus{
		Err:        err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}