package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
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
	if tr.sr == nil || tr.sr.Len() == 0 {
		if len(tr.token) == 0 {
			return len(b), nil
		}
		tr.sr = strings.NewReader(tr.token[0])
		tr.token = tr.token[1:]
	}
	return tr.sr.Read(b)
}

func TestAggregate(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		info  []item
		warn  []item
		err   []item
	}{
		{
			name: "standard",
			input: []string{
				"2017/01/02 03:04:05 foo.go:1: INFO foo",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo"}},
		},
		{
			name: "over 4096",
			input: []string{
				strings.Repeat("=", 4097) + "\n2017/01/02 03:04:05 foo",
				".go:1: INFO foo",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo"}},
			err:  []item{{1, strings.Repeat("=", 4097)}},
		},
		{
			name: "concated log",
			input: []string{
				"2017/01/02 03:04:05 foo.go:1: INFO foo",
				"2017/01/02 03:04:05 foo.go:1: WARNING foo",
				"2017/01/02 03:04:05 foo.go:1: ERROR foo",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo2017/01/02 03:04:05 foo.go:1: WARNING foo2017/01/02 03:04:05 foo.go:1: ERROR foo"}},
		},
		{
			name: "separated log",
			input: []string{
				"2017/01/02 03:04:05 foo.go:1: INFO foo\n",
				"2017/01/02 03:04:05 foo.go:1: WARNING foo\n",
				"2017/01/02 03:04:05 foo.go:1: ERROR foo\n",
				"2017/01/02 03:04:05 foo.go:1: CRITICAL foo\n",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo"}},
			warn: []item{
				{1, "2017/01/02 03:04:05 foo.go:1: WARNING foo"},
				{1, "2017/01/02 03:04:05 foo.go:1: ERROR foo"},
			},
			err: []item{{1, "2017/01/02 03:04:05 foo.go:1: CRITICAL foo"}},
		},
		{
			name: "separated log, and sleep",
			input: []string{
				"2017/01/02 03:04:05 foo.go:1: INFO foo\n\n2017/01/02 03:04:05 foo.go:1: WARNING foo\n",
				"2017/01/02 03:04:05 foo.go:1: CRITICAL foo\n",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo"}},
			warn: []item{{1, "2017/01/02 03:04:05 foo.go:1: WARNING foo"}},
			err:  []item{{1, "2017/01/02 03:04:05 foo.go:1: CRITICAL foo"}},
		},
		{
			name: "separated log, and sleep",
			input: []string{
				strings.Repeat("=", 4097) + "\n2017/01/02 03:04:05 foo.go:1: INFO foo\n2017/01/02 03:04:05 foo.go:1: ",
				"WARNING foo\n2017/01/02 03:04:05 foo.go:1: CRITICAL foo\n",
			},
			info: []item{{1, "2017/01/02 03:04:05 foo.go:1: INFO foo"}},
			warn: []item{{1, "2017/01/02 03:04:05 foo.go:1: WARNING foo"}},
			err: []item{
				{1, strings.Repeat("=", 4097)},
				{1, "2017/01/02 03:04:05 foo.go:1: CRITICAL foo"},
			},
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
			t.Fatalf("%s: info log: want: %v, got: %v", test.name, test.info, tl.info)
		}
		if !reflect.DeepEqual(tl.warn, test.warn) {
			t.Fatalf("%s: warn log: want: %v, got: %v", test.name, test.warn, tl.warn)
		}
		if !reflect.DeepEqual(tl.err, test.err) {
			t.Fatalf("%s: err log: want: %v, got: %v", test.name, test.err, tl.err)
		}
	}
}
