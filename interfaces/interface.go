package interfaces

type IdName interface {
	GetId() int64
	GetName() string
}

type ShouldUpdaterIdName interface {
	ShouldUpdater
	IdName
}

type ShouldUpdater interface {
	ShouldUpdate() bool
}
