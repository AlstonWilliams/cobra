package callerd

// Maybe we will have a caller which is implemented by GRPC, so this caller interface is here.
type Caller interface {
	Call(method, url string, params map[string]string) (bool, string)
}

func NewCaller() Caller{
	return Caller_http{}
}