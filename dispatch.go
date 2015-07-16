package main

import "strings"

func splitSub(argv []string) (string, []string) {
	if len(argv) == 0 || strings.HasPrefix(argv[0], "-") {
		return "", argv
	}
	return argv[0], argv[1:]
}

func dispatch(argv []string) int {
	subCmd, argv := splitSub(argv)
	fn, ok := routes[subCmd]
	if !ok {
		logger.Errorf("subcommand: %s not found", subCmd)
		return exitStatusError
	}
	return fn(argv)
}
