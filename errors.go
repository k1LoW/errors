package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var _ error = (*errorWithStack)(nil)

// MaxStackDepth is the maximum depth of the stack trace.
var MaxStackDepth = 50

// As is a wrapper for [errors.As].
//
// [errors.As]: https://pkg.go.dev/errors#As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is is a wrapper for [errors.Is].
//
// [errors.Is]: https://pkg.go.dev/errors#Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Join is a wrapper for [errors.Join].
//
// [errors.Join]: https://pkg.go.dev/errors#Join
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// New is a wrapper for [errors.New].
//
// [errors.New]: https://pkg.go.dev/errors#New
func New(text string) error {
	return errors.New(text)
}

// Unwrap is a wrapper for [errors.Unwrap].
//
// [errors.Unwrap]: https://pkg.go.dev/errors#Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// WithStack sets the stack trace for the given error.
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
		errws = &errorWithStack{
			Err:   err,
			stack: stack[:length],
		}
	}
	return errws
}

// StackTraces returns the stack traces of the given error(s).
func StackTraces(err error) stackTraces {
	je, ok := err.(joinError)
	if ok {
		// joined error
		var traces stackTraces
		errs := je.Unwrap()
		for _, e := range errs {
			traces = append(traces, StackTraces(e)...)
		}
		return traces
	}
	var errws *errorWithStack
	if !errors.As(err, &errws) {
		return nil
	}
	errws.genFrames()
	return stackTraces{errws}
}

type stackTraces []*errorWithStack

type errorWithStack struct {
	Err    error
	Frames []frame
	stack  []uintptr
}

type frame struct {
	Name string `json:"name"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func (traces stackTraces) String() string {
	var sb strings.Builder
	for _, errws := range traces {
		sb.WriteString(fmt.Sprintf("%s\n", errws.Error()))
		for _, frame := range errws.Frames {
			sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Name, frame.File, frame.Line))
		}
	}
	return sb.String()
}

func (errws *errorWithStack) Error() string {
	msg := errws.Err.Error()
	return msg
}

func (errws *errorWithStack) Unwrap() error {
	return errws.Err
}

func (errws *errorWithStack) MarshalJSON() ([]byte, error) {
	s := struct {
		Error  string  `json:"error"`
		Frames []frame `json:"frames"`
	}{
		Error:  errws.Error(),
		Frames: errws.Frames,
	}
	return json.Marshal(s)
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
