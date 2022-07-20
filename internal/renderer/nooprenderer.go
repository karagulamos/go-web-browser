package renderer

type NoopRenderer struct {
}

func NewNoopRenderer() Renderer {
	return &NoopRenderer{}
}

func (r *NoopRenderer) Invoke() *Result {
	return &Result{Error: ErrNotImplemented}
}
