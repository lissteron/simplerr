package simplerr

type ErrCodeInterface interface {
	HTTP() int
	GRPC() int
	Int() int
}

type ErrCode int

func (e ErrCode) HTTP() int {
	return int(e)
}

func (e ErrCode) GRPC() int {
	return int(e)
}

func (e ErrCode) Int() int {
	return int(e)
}
