package slice

func Map[E any, R any](src []E, fn func(E) R) []R {
	dst := make([]R, 0, len(src))
	for _, v := range src {
		dst = append(dst, fn(v))
	}
	return dst
}
