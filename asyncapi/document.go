package asyncapi

type Document interface {
	Extendable
	Servers()
}

type Server interface {
	Extendable
	Name() string
	HasName() bool
	URL() string
	HasURL() bool
	Protocol() string
	Variables() []ServerVariable
}

type ServerVariable interface {
	Extendable
	Name() string
	HasName() string
	DefaultValue() string
	AllowedValues() []string // Parser API spec says any[], but AsyncAPI mentions is []string
}

type Extendable interface {
	HasExtension(name string) bool
	Extension(name string) interface{}
}

type Identifiable interface {
	IDField() string
	ID() string
}
