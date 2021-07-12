package asyncapi

// Document is an object representing an AsyncAPI document.
// It's API implements https://github.com/asyncapi/parser-api/blob/master/docs/v1.md.
type Document interface {
	Extendable
	Servers()
}

// Server is an object representing a message broker, a server or any other kind of computer program capable of
// sending and/or receiving data.
type Server interface {
	Extendable
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
