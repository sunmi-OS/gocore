// Portions of this file are derived from go project
// Copyright 2021 The Go Authors.
// License: BSD-style
// Source: https://github.com/golang/go/tree/master/src/slices

// Package slices defines various functions useful with slices of any type.
package slice

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[S ~[]E, E any](s S, f func(E) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[S ~[]E, E comparable](s S, v E) bool {
	return Index(s, v) >= 0
}

// ContainsFunc reports whether at least one element e of s satisfies f(e).
func ContainsFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	return IndexFunc(s, f) >= 0
}

// Diff return the difference set of s1 and s2, that is, the set of all elements that belong to s1 but not s2.
func Diff[S ~[]E, E comparable](s1, s2 S) S {
	set := make(map[E]struct{}) // use an empty struct as the value because it takes up no space

	var diff S

	// add the element from s2 to the map
	for _, item := range s2 {
		set[item] = struct{}{}
	}

	// traverse s1, adding the difference set if the element is not in the map
	for _, item := range s1 {
		if _, found := set[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	return append(S([]E{}), s...)
}

// Intersect return the intersection set of s1 and s2, that is, the set of all elements that belong to both s1 and s2.
func Intersect[S ~[]E, E comparable](s1, s2 S) S {
	set := make(map[E]struct{}) // use an empty struct as the value because it takes up no space
	var intersection S

	// Create a set from the elements of the first slice.
	for _, item := range s1 {
		set[item] = struct{}{}
	}

	// Check if elements of the second slice are in the set.
	for _, item := range s2 {
		if _, exists := set[item]; exists {
			intersection = append(intersection, item)
			// Remove the item from the set to prevent duplicates in the intersection.
			delete(set, item)
		}
	}

	return intersection
}
