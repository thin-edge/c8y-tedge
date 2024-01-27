package runner

// type Runner struct {
// }

type Runner interface {
	Run() error
}

type Responder interface {
	ExitCode() int
	Stdout() []byte
}
