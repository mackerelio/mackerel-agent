package main

import "testing"

func TestSplitSub(t *testing.T) {
	func() {
		cmd, argv := splitSub([]string{})
		if cmd != "" {
			t.Errorf("cmd should be empty")
		}
		if len(argv) != 0 {
			t.Errorf("argv should be empty")
		}
	}()

	func() {
		cmd, argv := splitSub([]string{"-help"})
		if cmd != "" {
			t.Errorf("cmd should be empty")
		}
		if len(argv) != 1 || argv[0] != "-help" {
			t.Errorf("argv[0] should be '-help'")
		}
	}()

	func() {
		cmd, argv := splitSub([]string{"version", "-help"})
		if cmd != "version" {
			t.Errorf("cmd should be version")
		}
		if len(argv) != 1 || argv[0] != "-help" {
			t.Errorf("argv[0] should be '-help'")
		}
	}()
}

func TestDispatch(t *testing.T) {
	func() {
		code := dispatch([]string{"dummmmmmmy"})
		if code != exitStatusError {
			t.Errorf("exit code should be error")
		}
	}()

	func() {
		code := dispatch([]string{"version"})
		if code != exitStatusOK {
			t.Errorf("exit code should be ok")
		}
	}()
}
