package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	_ error          = (*errorWithStack)(nil)
	_ json.Marshaler = (*errorWithStack)(nil)
	_ fmt.Stringer   = (stackTraces)(nil)
)

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

	var errws *errorWithStack
	if errors.As(err, &errws) {
		return err
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])

	return &errorWithStack{
		Err:   err,
		stack: stack[:length],
	}
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
		return stackTraces{}
	}
	errws.genFrames()
	return stackTraces{errws}
}

// Errors returns all joined errors in the given error.
func Errors(err error) []error {
	je, ok := err.(joinError)
	if !ok {
		return []error{err}
	}
	errs := je.Unwrap()
	var splitted []error
	for _, e := range errs {
		errrs := Errors(e)
		if len(errrs) > 1 {
			splitted = append(splitted, Errors(e)...)
			continue
		}
		splitted = append(splitted, e)
	}
	return splitted
}

type stackTraces []*errorWithStack

type errorWithStack struct {
	Err    error
	Frames []Frame
	stack  []uintptr
}

type Frame struct {
	Name string `json:"name"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func (traces stackTraces) String() string {
	var sb strings.Builder
	for i, errws := range traces {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(errws.Error())
		for _, frame := range errws.Frames {
			sb.WriteString(fmt.Sprintf("\n%s\n\t%s:%d", frame.Name, frame.File, frame.Line))
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
		Frames []Frame `json:"frames"`
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
	if errws.Frames != nil {
		return
	}
	errws.Frames = make([]Frame, len(errws.stack))

	for i, pc := range errws.stack {
		// ref: https://github.com/go-errors/errors/blob/83795c27c02f5cdeaf9a5c3c3fd2709376f20b79/Frame.go#L36-L37
		fn := runtime.FuncForPC(pc - 1)
		name := fn.Name()
		file, line := fn.FileLine(pc - 1)
		errws.Frames[i] = Frame{
			Name: name,
			File: file,
			Line: line,
		}
	}
}
