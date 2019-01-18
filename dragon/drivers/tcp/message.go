package tcp

// Message represents the message structure contracts.
type Message interface {

	// MessageType to get message type
	MessageType() uint16

	// Serialize serialize Message into bytes.
	Serialize() ([]byte, error)
}

// MessageCodec represents the codec contracts.
type MessageCodec interface {

	// Encode encode the message to bytes
	Encode(Message) ([]byte, error)

	// Decode decode bytes to message
	Decode([]byte, chan Message) []byte
}
