package domain

type Engine interface {
	Register(processor CommandHandler)
	Execute(rawCommand *RawCommand) interface{}
}

type engine struct {
	processors map[string]CommandHandler
	cache      map[string]interface{}
}

func NewEngine() *engine {
	return &engine{}
}

func (e *engine) Register(processor CommandHandler) {
	e.processors[processor.CommandName()] = processor
}

func (e *engine) Execute(rawCommand *RawCommand) interface{} {
	processor := e.processors[rawCommand.Name]
	if processor == nil {
		return ErrorResult{Value: ErrUnknownCommand}
	}

	command, err := processor.Decode(rawCommand)
	if err != nil {
		return ErrorResult{Value: ErrMalformedCommand}
	}

	return processor.Process(e.cache, command)
}
