package lifecycle

type Initializable interface {
	Initialize() error
	Order() int
}

type Cleanable interface {
	Cleanup()
	Order() int
}
