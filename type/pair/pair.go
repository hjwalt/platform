package pair

func Of[L any, R any](left L, right R) Pair[L, R] {
	return Pair[L, R]{
		Left:  left,
		Right: right,
	}
}

type Pair[L any, R any] struct {
	Left  L
	Right R
}
