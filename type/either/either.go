package either

import "github.com/hjwalt/platform/type/optional"

func Left[L any, R any](left L) Either[L, R] {
	return Either[L, R]{
		left:  optional.Of(left),
		right: optional.Empty[R](),
	}
}

func Right[L any, R any](right R) Either[L, R] {
	return Either[L, R]{
		left:  optional.Empty[L](),
		right: optional.Of(right),
	}
}

type Either[L any, R any] struct {
	left  optional.Optional[L]
	right optional.Optional[R]
}

func (e *Either[L, R]) IsLeft() bool {
	return e.left.IsPresent()
}

func (e *Either[L, R]) Left() L {
	return e.left.Get()
}

func (e *Either[L, R]) IsRight() bool {
	return e.right.IsPresent()
}

func (e *Either[L, R]) Right() R {
	return e.right.Get()
}
