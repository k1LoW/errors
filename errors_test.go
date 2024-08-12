package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
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

func m() (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()

	return errors.New("error m")
}

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
		assertFrames(t, trace.Frames, "errors_test.a", "errors_test.b", "errors_test.c")
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
		assertFrames(t, trace.Frames, "errors_test.a", "errors_test.d", "errors_test.e")
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
			assertFrames(t, trace.Frames, "errors_test.a", "errors_test.g")
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
			assertFrames(t, trace.Frames, "errors_test.a", "errors_test.h")
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
			assertFrames(t, traceA.Frames, "errors_test.a", "errors_test.j")

			if want := "error i"; traceI.Err.Error() != want {
				t.Errorf("got: %s, want: %s", traceI.Err.Error(), want)
			}
			assertFrames(t, traceI.Frames, "errors_test.i", "errors_test.j")
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
			assertFrames(t, traceA.Frames, "errors_test.a", "errors_test.l")

			if want := "error i"; traceI.Err.Error() != want {
				t.Errorf("got: %s, want: %s", traceI.Err.Error(), want)
			}
			assertFrames(t, traceI.Frames, "errors_test.i", "errors_test.k", "errors_test.l")
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

	t.Run("zero frames", func(t *testing.T) {
		err := errors.New("error new")
		b, err := json.Marshal(errors.StackTraces(err))
		if err != nil {
			t.Fatal(err)
		}
		if want := []byte(`[]`); !bytes.Equal(b, want) {
			t.Errorf("got: %s, want: %s", b, want)
		}
	})
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

func TestWithDefer(t *testing.T) {
	err := m()
	traces := errors.StackTraces(err)
	if len(traces) != 1 {
		t.Errorf("got: %d, want: %d", len(traces), 1)
	}
	trace := traces[0]
	if want := "error m"; trace.Err.Error() != want {
		t.Errorf("got: %s, want: %s", trace.Err.Error(), want)
	}
	assertFrames(t, trace.Frames, "errors_test.m", "errors_test.m")
}

func TestWithPallarel(t *testing.T) {
	var err error
	mu := sync.Mutex{}
	sg := sync.WaitGroup{}
	sg.Add(2)
	go func() {
		mu.Lock()
		defer mu.Unlock()
		err = errors.Join(err, c())
		sg.Done()
	}()
	go func() {
		mu.Lock()
		defer mu.Unlock()
		err = errors.Join(err, l())
		sg.Done()
	}()
	sg.Wait()
	traces := errors.StackTraces(err)
	if len(traces) != 3 {
		t.Errorf("got: %d, want: %d", len(traces), 2)
	}
	for _, trace := range traces {
		switch trace.Err.Error() {
		case "error a":
			if strings.Contains(trace.Frames[1].Name, "errors_test.b") {
				assertFrames(t, trace.Frames, "errors_test.a", "errors_test.b", "errors_test.c")
			} else {
				assertFrames(t, trace.Frames, "errors_test.a", "errors_test.l")
			}
		case "error i":
			assertFrames(t, trace.Frames, "errors_test.i", "errors_test.k", "errors_test.l")
		default:
			t.Errorf("unknown error: %v", trace.Err)
		}
	}
}

func assertFrames(t *testing.T, frames []errors.Frame, names ...string) {
	t.Helper()
	for i, name := range names {
		if !strings.Contains(frames[i].Name, name) {
			t.Errorf("stack trace of %s (%d) not found: %v", name, i, frames)
		}
	}
}
