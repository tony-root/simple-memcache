package domain

type Encoder interface {
	EncodeResult(result interface{}) ([]byte, error)

	EncodeArray(value []string) []byte
	EncodeInt(value int) []byte
	EncodeBulkString(value string) []byte
	EncodeString(value string) []byte
	EncodeError(error error) []byte
	EncodeNil() []byte
}

type Decoder interface {
	Decode(rawCommand []byte) (*RawCommand, error)
}
