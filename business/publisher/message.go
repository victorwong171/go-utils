package pubsub

var (
	// DefaultExpire is the default message expiration time in seconds
	DefaultExpire = 300
)

// Message represents a message in the pub/sub system.
// It contains metadata about the event and the actual data payload.
type Message struct {
	// Event is the type or name of the event
	Event string

	// Data contains the actual message payload
	Data any

	// Source identifies where the message originated from
	Source string

	// TimeStamp is the RFC3339 formatted timestamp when the message was created
	TimeStamp string

	// Expire is the message expiration time in seconds
	// Messages that cannot be delivered within this time will be discarded
	Expire int
}

//
//type MessageBuilder struct {
//	options Message
//}
//
//func NewMessageBuilder() *MessageBuilder {
//	return &MessageBuilder{
//		options: Message{
//			Expire: DefaultExpire,
//		},
//	}
//}
//
//func (b *MessageBuilder) WithEvent(event string) *MessageBuilder {
//	b.options.Event = event
//	return b
//}
//
//func (b *MessageBuilder) WithData(data any) *MessageBuilder {
//	b.options.Data = data
//	return b
//}
//
//func (b *MessageBuilder) WithSource(source string) *MessageBuilder {
//	b.options.Source = source
//	return b
//}
//
//func (b *MessageBuilder) WithTimeStamp(timeStamp string) *MessageBuilder {
//	b.options.TimeStamp = timeStamp
//	return b
//}
//
//func (b *MessageBuilder) WithExpire(expire int) *MessageBuilder {
//	b.options.Expire = expire
//	return b
//}
//
//func (b *MessageBuilder) Build() *Message {
//	return &b.options
//}
