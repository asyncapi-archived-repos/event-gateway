package asyncapi

// Document is an object representing an AsyncAPI document.
// It's API implements https://github.com/asyncapi/parser-api/blob/master/docs/v1.md.
// NOTE: this interface is not completed yet.
type Document interface {
	Extendable
	Channels() []Channel
	HasChannels() bool
	ApplicationPublishableChannels() []Channel
	ApplicationPublishableMessages() []Message
	ApplicationPublishOperations() []Operation
	ApplicationSubscribableChannels() []Channel
	ApplicationSubscribableMessages() []Message
	ApplicationSubscribeOperations() []Operation
	ClientPublishableChannels() []Channel
	ClientPublishableMessages() []Message
	ClientPublishOperations() []Operation
	ClientSubscribableChannels() []Channel
	ClientSubscribableMessages() []Message
	ClientSubscribeOperations() []Operation
	Messages() []Message
	Server(name string) (Server, bool)
	Servers() []Server
	HasServers() bool
}

// Channel is an addressable component, made available by the server, for the organization of messages.
// Producer applications send messages to channels and consumer applications consume messages from channels.
type Channel interface {
	Extendable
	Identifiable
	Path() string // Path is the identifier
	Description() string
	HasDescription() bool
	Parameters() []ChannelParameter
	HasParameters() bool
	Operations() []Operation
	Messages() []Message
}

// ChannelParameter describes a parameter included in a channel name.
type ChannelParameter interface {
	Extendable
	Identifiable
	Name() string
	Description() string
	Schema() Schema
}

// OperationType is the type of an operation.
type OperationType string

// Operation describes a publish or a subscribe operation.
// This provides a place to document how and why messages are sent and received.
type Operation interface {
	Extendable
	ID() string
	Description() string
	HasDescription() bool
	IsApplicationPublishing() bool
	IsApplicationSubscribing() bool
	IsClientPublishing() bool
	IsClientSubscribing() bool
	Messages() []Message
	Summary() string
	HasSummary() bool
	Type() OperationType
}

// Message describes a message received on a given channel and operation.
type Message interface {
	Extendable
	UID() string
	Name() string
	Title() string
	HasTitle() bool
	Description() string
	HasDescription() bool
	Summary() string
	HasSummary() bool
	ContentType() string
	Payload() Schema
}

// Schema is an object that allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
// This object is a superset of the JSON Schema Specification Draft 07.
type Schema interface {
	Extendable
	ID() string
	AdditionalItems() Schema
	AdditionalProperties() Schema // TODO (boolean | Schema)
	AllOf() []Schema
	AnyOf() []Schema
	CircularProps() []string
	Const() interface{}
	Contains() Schema
	ContentEncoding() string
	ContentMediaType() string
	Default() interface{}
	Definitions() map[string]Schema
	Dependencies() map[string]Schema // TODO Map[string, Schema|string[]]
	Deprecated() bool
	Description() string
	HasDescription() bool
	Discriminator() string
	Else() Schema
	Enum() []interface{}
	Examples() []interface{}
	ExclusiveMaximum() *float64
	ExclusiveMinimum() *float64
	Format() string
	HasCircularProps() bool
	If() Schema
	IsCircular() bool
	Items() []Schema // TODO Schema | Schema[]
	Maximum() *float64
	MaxItems() *float64
	MaxLength() *float64
	MaxProperties() *float64
	Minimum() *float64
	MinItems() *float64
	MinLength() *float64
	MinProperties() *float64
	MultipleOf() *float64
	Not() Schema
	OneOf() Schema
	Pattern() string
	PatternProperties() map[string]Schema
	Properties() map[string]Schema
	Property(name string) Schema
	PropertyNames() Schema
	ReadOnly() bool
	Required() string // TODO string[]
	Then() Schema
	Title() string
	Type() []string // TODO // string | string[]
	UID() string
	UniqueItems() bool
	WriteOnly() bool
}

// Server is an object representing a message broker, a server or any other kind of computer program capable of
// sending and/or receiving data.
type Server interface {
	Extendable
	Identifiable
	Name() string
	HasName() bool
	Description() string
	HasDescription() bool
	URL() string
	HasURL() bool
	Protocol() string
	HasProtocol() bool
	Variables() []ServerVariable
}

// ServerVariable is an object representing a Server Variable for server URL template substitution.
type ServerVariable interface {
	Extendable
	Identifiable
	Name() string
	HasName() bool
	DefaultValue() string
	AllowedValues() []string // Parser API spec says any[], but AsyncAPI mentions is []string
}

// Extendable means the object can have extensions.
// The extensions properties are implemented as patterned fields that are always prefixed by "x-".
// See https://www.asyncapi.com/docs/specifications/v2.0.0#specificationExtensions.
type Extendable interface {
	HasExtension(name string) bool
	Extension(name string) interface{}
}

// Identifiable identifies objects. Some objects can have fields that identify themselves as unique resources.
// For example: `id` and `name` fields.
type Identifiable interface {
	IDField() string
	ID() string
}
