package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/k1LoW/errors"
)

var errA = errors.New("error a")
var errF = errors.New("error f")
var errI = errors.New("error i")

func a() error { return errors.WithStack(errA) }
func b() error { return errors.WithStack(a()) }
func c() error { return errors.WithStack(b()) }

func d() error { return fmt.Errorf("error d: %w", a()) }
func e() error { return fmt.Errorf("error e: %w", d()) }

func f() error { return errF }
func g() error { return errors.Join(a(), f()) }
func h() error { return errors.Join(f(), a()) }

func i() error { return errors.WithStack(errI) }
func j() error { return errors.Join(a(), i()) }
func k() error { return fmt.Errorf("error k: %w", i()) }
func l() error { return errors.Join(a(), k()) }

func TestWithStack(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		err := errors.WithStack(nil)
		if err != nil {
			t.Errorf("got: %v, want: %v", err, nil)
		}
	})

	t.Run("The first stack is retained.", func(t *testing.T) {
		err := c()
		if !errors.Is(err, errA) {
			t.Error("error not found")
		}
		traces := errors.StackTraces(err)
		if len(traces) != 1 {
			t.Errorf("got: %d, want: %d", len(traces), 1)
		}
		trace := traces[0]
		if !strings.Contains(trace.Frames[0].Name, ".a") {
			t.Errorf("stack trace of a() not found: %v", trace.Frames[0])
		}
		if !strings.Contains(trace.Frames[1].Name, ".b") {
			t.Errorf("stack trace of b() not found: %v", trace.Frames[1])
		}
		if !strings.Contains(trace.Frames[2].Name, ".c") {
			t.Errorf("stack trace of c() not found: %v", trace.Frames[2])
		}
	})

	t.Run("wrapped stack is retained.", func(t *testing.T) {
		// errA -> a:WithStack() -> d:Errorf("%w") -> e:Errorf("%w")
		err := e()
		if want := "error e: error d: error a"; err.Error() != want {
			t.Errorf("got: %s, want: %s", err.Error(), want)
		}
		if !errors.Is(err, errA) {
			t.Error("error not found")
		}
		traces := errors.StackTraces(err)
		if len(traces) != 1 {
			t.Errorf("got: %d, want: %d", len(traces), 1)
		}
		trace := traces[0]
		if want := "error a"; trace.Err.Error() != want {
			t.Errorf("got: %s, want: %s", trace.Err.Error(), want)
		}
		if !strings.Contains(trace.Frames[0].Name, ".a") {
			t.Errorf("stack trace of a() not found: %v", trace.Frames[0])
		}
		if !strings.Contains(trace.Frames[1].Name, ".d") {
			t.Errorf("stack trace of d() not found: %v", trace.Frames[1])
		}
		if !strings.Contains(trace.Frames[2].Name, ".e") {
			t.Errorf("stack trace of e() found: %v", trace.Frames[2])
		}
	})

	t.Run("joined stack is retained.", func(t *testing.T) {
		{
			err := g()
			if want := "error a\nerror f"; err.Error() != want {
				t.Errorf("got: %s, want: %s", err.Error(), want)
			}
			if !errors.Is(err, errA) {
				t.Error("error not found")
			}
			traces := errors.StackTraces(err)
			if len(traces) != 1 {
				t.Errorf("got: %d, want: %d", len(traces), 1)
			}
			trace := traces[0]
			if want := "error a"; trace.Err.Error() != want {
				t.Errorf("got: %s, want: %s", trace.Err.Error(), want)
			}
			if !strings.Contains(trace.Frames[0].Name, ".a") {
				t.Errorf("stack trace of a() not found: %v", trace.Frames[0])
			}
			if !strings.Contains(trace.Frames[1].Name, ".g") {
				t.Errorf("stack trace of g() not found: %v", trace.Frames[1])
			}
		}
		{
			err := h()
			if want := "error f\nerror a"; err.Error() != want {
				t.Errorf("got: %s, want: %s", err.Error(), want)
			}
			if !errors.Is(err, errA) {
				t.Error("error not found")
			}
			traces := errors.StackTraces(err)
			if len(traces) != 1 {
				t.Errorf("got: %d, want: %d", len(traces), 1)
			}
			trace := traces[0]
			if want := "error a"; trace.Err.Error() != want {
				t.Errorf("got: %s, want: %s", trace.Err.Error(), want)
			}
			if !strings.Contains(trace.Frames[0].Name, ".a") {
				t.Errorf("stack trace of a() not found: %v", trace.Frames[0])
			}
			if !strings.Contains(trace.Frames[1].Name, ".h") {
				t.Errorf("stack trace of g() not found: %v", trace.Frames[1])
			}
		}
		{
			err := j()
			if want := "error a\nerror i"; err.Error() != want {
				t.Errorf("got: %s, want: %s", err.Error(), want)
			}
			if !errors.Is(err, errA) {
				t.Error("error not found")
			}
			if !errors.Is(err, errI) {
				t.Error("error not found")
			}
			traces := errors.StackTraces(err)
			if len(traces) != 2 {
				t.Errorf("got: %d, want: %d", len(traces), 1)
			}
			traceA := traces[0]
			traceI := traces[1]

			if want := "error a"; traceA.Err.Error() != want {
				t.Errorf("got: %s, want: %s", traceA.Err.Error(), want)
			}
			if !strings.Contains(traceA.Frames[0].Name, ".a") {
				t.Errorf("stack trace of a() not found: %v", traceA.Frames[0])
			}
			if !strings.Contains(traceA.Frames[1].Name, ".j") {
				t.Errorf("stack trace of g() not found: %v", traceA.Frames[1])
			}

			if want := "error i"; traceI.Err.Error() != want {
				t.Errorf("got: %s, want: %s", traceI.Err.Error(), want)
			}
			if !strings.Contains(traceI.Frames[0].Name, ".i") {
				t.Errorf("stack trace of i() not found: %v", traceI.Frames[0])
			}
			if !strings.Contains(traceI.Frames[1].Name, ".j") {
				t.Errorf("stack trace of j() not found: %v", traceI.Frames[1])
			}
		}

		{
			err := l()
			if want := "error a\nerror k: error i"; err.Error() != want {
				t.Errorf("got: %s, want: %s", err.Error(), want)
			}
			if !errors.Is(err, errA) {
				t.Error("error not found")
			}
			if !errors.Is(err, errI) {
				t.Error("error not found")
			}
			traces := errors.StackTraces(err)
			if len(traces) != 2 {
				t.Errorf("got: %d, want: %d", len(traces), 1)
			}
			traceA := traces[0]
			traceI := traces[1]

			if want := "error a"; traceA.Err.Error() != want {
				t.Errorf("got: %s, want: %s", traceA.Err.Error(), want)
			}
			if !strings.Contains(traceA.Frames[0].Name, ".a") {
				t.Errorf("stack trace of a() not found: %v", traceA.Frames[0])
			}
			if !strings.Contains(traceA.Frames[1].Name, ".l") {
				t.Errorf("stack trace of l() not found: %v", traceA.Frames[1])
			}

			if want := "error i"; traceI.Err.Error() != want {
				t.Errorf("got: %s, want: %s", traceI.Err.Error(), want)
			}
			if !strings.Contains(traceI.Frames[0].Name, ".i") {
				t.Errorf("stack trace of i() not found: %v", traceI.Frames[0])
			}
			if !strings.Contains(traceI.Frames[1].Name, ".k") {
				t.Errorf("stack trace of k() not found: %v", traceI.Frames[1])
			}
			if !strings.Contains(traceI.Frames[2].Name, ".l") {
				t.Errorf("stack trace of l() not found: %v", traceI.Frames[2])
			}
		}
	})
}

func TestJSON(t *testing.T) {
	err := l()
	b, err := json.Marshal(errors.StackTraces(err))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(b, []byte(`"error":"error a"`)) {
		t.Error(`"error":"error a" not found`)
	}
	if !bytes.Contains(b, []byte(`"frames":[{`)) {
		t.Error(`"frames":[{ not found`)
	}
}

func TestString(t *testing.T) {
	err := l()
	s := errors.StackTraces(err).String()
	t.Log(s)
	if !strings.Contains(s, "error a\n") {
		t.Error(`"error a\n\t" not found`)
	}
	if !strings.Contains(s, ".a\n\t") {
		t.Error(`".a\n\t" not found`)
	}
}

func TestSlogJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewJSONHandler(buf, nil))
	err := l()
	logger.Info("test", slog.Any("stacktracs", errors.StackTraces(err)))
	if !strings.Contains(buf.String(), `"stacktracs"`) {
		t.Error("stacktracs not found")
	}
	if !strings.Contains(buf.String(), `"error":"error a"`) {
		t.Error(`"error":"error a" not found`)
	}
	if !strings.Contains(buf.String(), `"frames":[{`) {
		t.Error(`"frames":[{ not found`)
	}
}

func TestSlogText(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(buf, nil))
	err := l()
	logger.Info("test", slog.Any("stacktracs", errors.StackTraces(err)))
	t.Log(buf.String())
	if !strings.Contains(buf.String(), `stacktracs=`) {
		t.Error("stacktracs= not found")
	}
	if !strings.Contains(buf.String(), "error a\\n") {
		t.Error("error a\\n not found")
	}
}
