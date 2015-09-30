package timewriter

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

type MockWriter struct {
	buf []byte
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	m.buf = append(m.buf, p...)
	return len(p), nil
}

func (m *MockWriter) Check(expect []byte) error {
	r := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d+)? [+-]{1}\d{4} UTC: `)
	s := r.Split(string(m.buf), -1)

	if len(s) != 3 {
		return errors.New("Could not find write")
	}

	if c := s[1]; c != string(expect) {
		return fmt.Errorf("Found wrong write: %+q", c)
	}

	return nil
}

func testBytes(t *testing.T, buf []byte) {
	m := MockWriter{make([]byte, 0)}

	l, err := New(&m)
	if err != nil {
		t.Fatalf("%s", err)
	}

	n, err := l.Write([]byte(buf))
	if n != len(buf) {
		t.Fatalf("Write did not return the length of the passed buffer: %d/%d", n, len(buf))
	}
	if err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	if err = m.Check(buf); err != nil {
		t.Fatalf("Check failed: %s", err)
	}
}

func TestR(t *testing.T) {
	testBytes(t, []byte("\r"))
}

func TestN(t *testing.T) {
	testBytes(t, []byte("\n"))
}

func TestRN(t *testing.T) {
	testBytes(t, []byte("\r\n"))
}
