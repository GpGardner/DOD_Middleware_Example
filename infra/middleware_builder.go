package infra

type MiddlewareBuilder[In any, Out any] struct {
	middlewares []Middleware[In, Out]
}

// Chain composes middlewares around a base operation.
func Chain[In any, Out any](base RepoOp[In, Out], mws ...Middleware[In, Out]) RepoOp[In, Out] {
	for i := len(mws) - 1; i >= 0; i-- {
		base = mws[i](base)
	}
	return base
}

func (b *MiddlewareBuilder[In, Out]) Add(mw Middleware[In, Out]) {
	b.middlewares = append(b.middlewares, mw)
}

func (b *MiddlewareBuilder[In, Out]) Build(base RepoOp[In, Out]) RepoOp[In, Out] {
	return Chain(base, b.middlewares...)
}
