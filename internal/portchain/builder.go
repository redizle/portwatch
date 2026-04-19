package portchain

// Builder constructs a Chain using a fluent API.
type Builder struct {
	chain *Chain
	err   error
}

// NewBuilder returns a new Builder.
func NewBuilder() *Builder {
	return &Builder{chain: New()}
}

// Add appends a handler; the first error is retained.
func (b *Builder) Add(h Handler) *Builder {
	if b.err != nil {
		return b
	}
	b.err = b.chain.Use(h)
	return b
}

// Build returns the constructed Chain and any accumulated error.
func (b *Builder) Build() (*Chain, error) {
	return b.chain, b.err
}
