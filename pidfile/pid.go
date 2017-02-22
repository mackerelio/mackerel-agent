package pidfile

// ExistsPid checks if pid exists
func ExistsPid(pid int) bool {
	return existsPid(pid)
}
