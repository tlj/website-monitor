package monitors

type Monitor interface {
	Check(check Check) (bool, error)
}
