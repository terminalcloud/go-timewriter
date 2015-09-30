// A wrapper over io.Writer to begin each line with a timestamp
package timewriter

import (
	"fmt"
	"io"
	"time"
)

// A Writer that begins each line with a timestamp
type TimeWriter struct {
	w io.Writer
}

func New(w io.Writer) (*TimeWriter, error) {
	// Begin the log with a timestamp
	now := []byte(fmt.Sprintf("%s: ", time.Now().UTC().String()))
	if _, err := w.Write(now); err != nil {
		return nil, err
	}

	return &TimeWriter{w}, nil
}

func (l *TimeWriter) Write(p []byte) (n int, err error) {
	now := []byte(fmt.Sprintf("%s: ", time.Now().UTC().String()))
	i := 0

	// Flush to w on newline and append a timestamp
	for j, b := range p {
		if b == '\r' || b == '\n' {
			if b == '\r' {
				if j+1 < len(p) && p[j+1] == '\n' {
					continue
				}
			}

			w, err := l.w.Write(p[i : j+1])
			n += w
			i = j + 1
			if err != nil {
				return n, err
			}

			_, err = l.w.Write(now)
			if err != nil {
				return n, err
			}
		}
	}

	// Write remainder to w
	if i != len(p) {
		w, err := l.w.Write(p[i:len(p)])
		n += w
		if err != nil {
			return n, err
		}
	}

	return n, nil
}
