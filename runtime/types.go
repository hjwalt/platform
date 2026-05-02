package runtime

type Runtime interface {
	Start() error
	Stop()
}

type Holder interface {
	Add(runtimes ...Runtime)
	Block()
}
