package server

type HelloService struct {
}

func (h *HelloService) Hello(request string, hello *string) error {
	*hello = "hello" + request
	return nil
}
