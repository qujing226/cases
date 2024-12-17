package slice

import (
	"errors"
	"golang.org/x/exp/constraints"
)

type Integer interface {
	constraints.Integer | constraints.Float
}

func Sum[T Integer](res ...T) (t T) {
	for _, re := range res {
		t += re
	}
	return
}

func Max[T Integer](slice []T) (t T, err error) {
	if len(slice) == 0 {
		err = errors.New("slice is empty")
		return
	}
	t = slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] > t {
			t = slice[i]
		}
	}
	return
}

func Min[T Integer](slice []T) (t T, err error) {
	if len(slice) == 0 {
		err = errors.New("slice is empty")
		return
	}
	t = slice[0]
	for i := 1; i < len(slice); i++ {
		if t > slice[i] {
			t = slice[i]
		}
	}
	return
}

func Find[T comparable](slice []T, obj T) (b bool, err error) {
	if len(slice) == 0 {
		err = errors.New("slice is empty")
		return
	}
	for i := 0; i < len(slice); i++ {
		if slice[i] == obj {
			b = true
			return
		}
	}
	return
}

func Insert[T any](slice []T, obj T, index int) ([]T, error) {
	if len(slice) < 0 || index > len(slice) {
		return nil, errors.New("index out of range")
	}
	slice = append(slice[:index], append([]T{obj}, slice[index:]...)...)
	return slice, nil
}

func Delete[T comparable](slice []T, obj T) ([]T, error) {
	if len(slice) == 0 {
		return nil, errors.New("slice is empty")
	}
	for i := 0; i < len(slice); i++ {
		if slice[i] == obj {
			copy(slice[i:], slice[i+1:])
			slice = slice[:len(slice)-1]
			if cap(slice) > len(slice)*2 {
				newSlice := make([]T, len(slice))
				copy(newSlice, slice)
				slice = newSlice
			}
			return slice, nil
		}
	}
	return nil, errors.New("not found")
}
