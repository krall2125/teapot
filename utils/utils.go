package utils

func Filter[E any](s []E, f func(E) bool) []E {
	s2 := make([]E, 0, len(s))
	for _, e := range s {
		if f(e) {
			 s2 = append(s2, e)
		}
	}
	return s2
}
