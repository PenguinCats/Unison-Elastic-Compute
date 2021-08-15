package data

type MessageDataType int

const (
	MessageDataTypeContainerCreate MessageDataType = iota
	MessageDataTypeContainerStart
	MessageDataTypeContainerStop
	MessageDataTypeContainerRemove
	MessageDataTypeContainerProfile
	MessageDataTypeContainerStatus
)

type Message struct {
	MessageType MessageDataType `json:"message_type"`
	Value       []byte          `json:"value"`
}
