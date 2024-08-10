package errors

import (
	"errors"
	"fmt"
	"runtime"
)

var _ error = (*errorWithStack)(nil)

var MaxStackDepth = 50

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func New(text string) error {
	return errors.New(text)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])

	var errws *errorWithStack
	if errors.As(err, &errws) {
		return errws
	} else {
		if _, ok := err.(error); !ok {
			err = fmt.Errorf("%v", err)
		}
		errws = &errorWithStack{
			Err:   err,
			stack: stack[:length],
		}
	}
	return errws
}

func StackTraces(err error) []*errorWithStack {
	je, ok := err.(joinError)
	if ok {
		// joined error
		var errwss []*errorWithStack
		errs := je.Unwrap()
		for _, e := range errs {
			errwss = append(errwss, StackTraces(e)...)
		}
		return errwss
	}
	var errws *errorWithStack
	if !errors.As(err, &errws) {
		return nil
	}
	errws.genFrames()
	return []*errorWithStack{errws}
}

type errorWithStack struct {
	Err    error   `json:"error"`
	Frames []frame `json:"frames"`
	stack  []uintptr
}

type frame struct {
	Name string `json:"name"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func (errws *errorWithStack) Error() string {
	msg := errws.Err.Error()
	return msg
}

func (errws *errorWithStack) Unwrap() error {
	return errws.Err
}

type joinError interface {
	Unwrap() []error
}

func (errws *errorWithStack) genFrames() {
	if errws.Frames == nil {
		errws.Frames = make([]frame, len(errws.stack))

		for i, pc := range errws.stack {
			// ref: https://github.com/go-errors/errors/blob/83795c27c02f5cdeaf9a5c3c3fd2709376f20b79/Frame.go#L36-L37
			fn := runtime.FuncForPC(pc - 1)
			name := fn.Name()
			file, line := fn.FileLine(pc - 1)
			errws.Frames[i] = frame{
				Name: name,
				File: file,
				Line: line,
			}
		}
	}
}