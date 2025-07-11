package api

type Result[T any] struct {
	data T
	err  error
}

func NewResult[T any](data T, err error) Result[T] {
	return Result[T]{
		data: data,
		err:  err,
	}
}

func (r Result[T]) Sucess() bool {
	return r.err == nil
}

func (r Result[T]) Error() error {
	return r.err
}

func (r Result[T]) Data() T {
	return r.data
}
