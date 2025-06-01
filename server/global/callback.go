package global

// GenServer behavior needs to implement this interface
type GenServer interface {
	Start()
	HandleCall(GenReq) Reply
	HandleCast(GenReq)
	HandleInfo(GenReq)
	Terminate()
}
