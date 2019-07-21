package domain

type RawCommand struct {
	Name string
	Args []string
}

type Command interface {
	Name() string
}
