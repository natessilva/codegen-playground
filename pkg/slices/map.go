package slices

func Map[T, V any](ts []T, fn func(t T) V) []V {
	vs := make([]V, 0, len(ts))
	for _, t := range ts {
		vs = append(vs, fn(t))
	}
	return vs
}
