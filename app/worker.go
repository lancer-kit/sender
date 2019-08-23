package app

type Worker interface {
	Init() error
	Run(chan<- error)
}
