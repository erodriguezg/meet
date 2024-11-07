package sliceutils

func Filter[T any](input []T, condition func(T) bool) []T {
	var result []T
	for _, v := range input {
		if condition(v) {
			result = append(result, v)
		}
	}
	return result
}

func Map[T any, U any](input []T, mapper func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = mapper(v)
	}
	return result
}

func MapWithError[T any, U any](input []T, mapper func(T) (U, error)) ([]U, error) {
	result := make([]U, len(input))
	for i, v := range input {
		aux, err := mapper(v)
		if err != nil {
			return nil, err
		}
		result[i] = aux
	}
	return result, nil
}

func Reduce[T any, U any](input []T, initial U, accumulator func(U, T) U) U {
	result := initial
	for _, v := range input {
		result = accumulator(result, v)
	}
	return result
}

func Find[T any](input []T, condition func(T) bool) *T {
	for _, v := range input {
		if condition(v) {
			return &v
		}
	}
	return nil
}

func Any[T any](input []T, condition func(T) bool) bool {
	for _, v := range input {
		if condition(v) {
			return true
		}
	}
	return false
}

func All[T any](input []T, condition func(T) bool) bool {
	for _, v := range input {
		if !condition(v) {
			return false
		}
	}
	return true
}
