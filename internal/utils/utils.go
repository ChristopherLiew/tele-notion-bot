package main

// Last retrives the last element from an array like data structure (i.e. array, slice).
func Last[E any](s []E) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[len(s)-1], true
}

// Write to toml file
