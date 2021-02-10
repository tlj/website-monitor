package notifiers

type Notifier interface {
	Notify(msg string) error
}
