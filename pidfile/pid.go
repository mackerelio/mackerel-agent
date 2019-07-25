package pidfile

// ExistsPid checks if pid exists
func ExistsPid(pid int) bool {
	return existsPid(pid)
}

// GetCmdName gets the command name of pid
func GetCmdName(pid int) string {
	return getCmdName(pid)
}
