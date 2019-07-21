package domain

type CommandHandler interface {
	CommandName() string
	Decode(rawCommand *RawCommand) (Command, error)
	Process(cache map[string]interface{}, command Command) interface{}
}
