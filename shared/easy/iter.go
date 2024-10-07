// This is a companion to [[slices]] [[maps]] and [[cmp]] and is built on [[iter]]

package easy

import (
	"iter"
	"slices"
)

func Map[T, U any](it iter.Seq[T], fn func(T) U) iter.Seq[U] {
	next, stop := iter.Pull(it)
	return iter.Seq[U](func(yield func(U) bool) {
		n, ok := next()
		if ok {
			yield(fn(n))
		} else {
			stop()
		}
	})
}

func Map2[K, V, U any](it iter.Seq2[K, V], fn func(K, V) U) iter.Seq2[K, U] {
	next, stop := iter.Pull2(it)
	return iter.Seq2[K, U](func(yield func(K, U) bool) {
		k, v, ok := next()
		if ok {
			yield(k, fn(k, v))
		} else {
			stop()
		}
	})
}

func Filter[T comparable](it iter.Seq[T], fn func(T) bool) iter.Seq[T] {
	next, stop := iter.Pull(it)
	return iter.Seq[T](func(yield func(T) bool) {
		n, ok := next()
		if ok && fn(n) {
			yield(n)
		} else {
			stop()
		}
	})
}

func Filter2[K, V comparable](it iter.Seq2[K, V], fn func(K, V) bool) iter.Seq2[K, V] {
	next, stop := iter.Pull2(it)
	return iter.Seq2[K, V](func(yield func(K, V) bool) {
		k, v, ok := next()
		if ok && fn(k, v) {
			yield(k, v)
		} else {
			stop()
		}
	})
}

func Reduce[T, A any](it iter.Seq[T], a A, fn func(A, T) A) A {
	for v := range it {
		a = fn(a, v)
	}
	return a
}

func Reduce2[K , V, A any](it iter.Seq2[K, V], a A, fn func(A, K, V) A) A {
	for k, v := range it {
		a = fn(a, k, v)
	}
	return a
}

func Flatten[T any](ar [][]T) []T {
	a := make([]T, 0, len(ar)*2)
	for _, v := range ar {
		a = append(a, v...)
	}
	return slices.Clip(a)
}

func KeyOr[K comparable, V any](m map[K]V, key K, or V) V {
	for k, v := range m {
		if k == key {
			return v
		}
	}
	return or
}

func FindOr[T comparable](ar []T, needle T, or T) T {
	i := slices.Index(ar, needle)
	return Ternary(i >= 0, ar[i], or)
}

func PopOr[T any](ar []T, or T) (T, []T) {
	if len(ar) > 0 {
		return ar[0], ar[1:]
	}
	return or, ar
}

func FirstOr[T any](ar []T, or T) T {
	return Ternary(len(ar) > 0, ar[0], or)
}
