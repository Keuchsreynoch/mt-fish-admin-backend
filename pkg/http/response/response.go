package utls

import "reflect"

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

func NewResponse(message string, statusCode int, data interface{}) Response {
	return Response{
		Success:    true,
		Message:    message,
		StatusCode: statusCode,
		Data:       normalizeResponseData(data),
	}
}

type ResponseWithPaging struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int         `json:"total"`
}

func NewResponseWithPaging(message string, statusCode int, data interface{}, page int, perpage int, total int) ResponseWithPaging {
	return ResponseWithPaging{
		Success:    true,
		Message:    message,
		StatusCode: statusCode,
		Data:       normalizeResponseData(data),
		Page:       page,
		PerPage:    perpage,
		Total:      total,
	}
}

func normalizeResponseData(data interface{}) interface{} {
	if data == nil {
		return []interface{}{}
	}

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
		if v.IsNil() {
			return []interface{}{}
		}
	}

	return data
}
