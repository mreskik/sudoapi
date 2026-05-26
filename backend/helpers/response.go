package helpers

type Response struct {
	Code    *int   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewResponse() *Response {
	return &Response{}
}
func (res *Response) SetCode(Code int) *Response {
	res.Code = &Code
	return res
}

func (res *Response) SetMessage(Message string) *Response {
	res.Message = Message
	return res
}
func (res *Response) SetData(Data any) *Response {
	res.Data = Data
	return res
}

func (res *Response) Success() *Response {
	code := 0
	res.Code = &code
	res.Message = "success"
	return res
}
func (res *Response) GeneralError() *Response {
	code := 100
	res.Code = &code
	res.Message = "error"
	return res
}
