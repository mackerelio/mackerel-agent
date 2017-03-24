package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

type item struct {
	eid uint32
	msg string
}

type testLogger struct {
	info []item
	warn []item
	err  []item
}

func (l *testLogger) Info(eid uint32, msg string) error {
	l.info = append(l.info, item{eid, msg})
	return nil
}

func (l *testLogger) Warning(eid uint32, msg string) error {
	l.warn = append(l.warn, item{eid, msg})
	return nil
}

func (l *testLogger) Error(eid uint32, msg string) error {
	l.err = append(l.err, item{eid, msg})
	return nil
}

type testWriteCloser struct {
}

func (wc *testWriteCloser) Write(b []byte) (int, error) {
	return len(b), nil
}

func (wc *testWriteCloser) Close() error {
	return nil
}

type testReader struct {
	token []string
	sr    *strings.Reader
}

func (tr *testReader) Read(b []byte) (int, error) {
	if len(tr.token) == 0 {
		return 0, io.EOF
	}
	total := len(b)
	rest := total
	for {
		if tr.sr == nil || tr.sr.Len() == 0 {
			if len(tr.token) == 0 {
				return total - rest, nil
			}
			if len(tr.token[0]) == 0 {
				time.Sleep(20 * time.Millisecond)
				tr.token = tr.token[1:]
			}
			tr.sr = strings.NewReader(tr.token[0])
			tr.token = tr.token[1:]
		}
		n, err := tr.sr.Read(b)
		rest -= n
		if err == nil && rest == 0 {
			return total, err
		}
		b = b[n:]
	}
}

func TestAggregate(t *testing.T) {
	tests := []struct {
		input []string
		info  []item
		warn  []item
		err   []item
	}{
		{
			input: []string{
				"2017/01/02 03:04:05 foo.go:1: INFO foo",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo"}},
		},
		{
			input: []string{
				strings.Repeat("=", 4097) + "\n2017/01/02 03:04:05 foo",
				".go:1: INFO foo",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo"}},
		},
	}

	for _, test := range tests {
		tl := &testLogger{}

		h := &handler{
			elog: tl,
			w:    &testWriteCloser{},
			r:    &testReader{test.input, nil},
		}
		h.aggregate()
		h.wg.Wait()

		if !reflect.DeepEqual(tl.info, test.info) {
			t.Fatalf("info log: want: %v, got: %v", test.info, tl.info)
		}
		if !reflect.DeepEqual(tl.warn, test.warn) {
			t.Fatalf("warn log: want: %v, got: %v", test.warn, tl.warn)
		}
		if !reflect.DeepEqual(tl.err, test.err) {
			t.Fatalf("err log: want: %v, got: %v", test.err, tl.err)
		}
	}
}
