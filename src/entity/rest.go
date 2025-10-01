package entity

type HTTPResp struct {
	Message    HTTPMessage `json:"message"`
	Meta       Meta        `json:"metadata"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type HTTPMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Meta struct {
	Path       string     `json:"path"`
	StatusCode int        `json:"statusCode"`
	Status     string     `json:"status"`
	Message    string     `json:"message"`
	Error      *MetaError `json:"error,omitempty"`
}

type MetaError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	CurrentPage     int64    `json:"currentPage"`
	CurrentElements int64    `json:"currentElements"`
	TotalPages      int64    `json:"totalPages"`
	TotalElements   int64    `json:"totalElements"`
	SortBy          []string `json:"sortBy"`
	CursorStart     *string  `json:"cursorStart,omitempty"`
	CursorEnd       *string  `json:"cursorEnd,omitempty"`
}
