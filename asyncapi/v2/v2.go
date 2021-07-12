package v2

import "github.com/asyncapi/event-gateway/asyncapi"

type Document struct {
	Extendable
	ServersField map[string]Server `mapstructure:"servers"`
}

func (d Document) Servers() []asyncapi.Server {
	var servers []asyncapi.Server
	for _, s := range d.ServersField {
		servers = append(servers, s)
	}

	return servers
}

func (d Document) HasServers() bool {
	return len(d.ServersField) > 0
}

type Server struct {
	Extendable
	NameField        string                    `mapstructure:"name"`
	DescriptionField string                    `mapstructure:"description"`
	ProtocolField    string                    `mapstructure:"protocol"`
	URLField         string                    `mapstructure:"url"`
	VariablesField   map[string]ServerVariable `mapstructure:"variables"`
}

func (s Server) Variables() []asyncapi.ServerVariable {
	var vars []asyncapi.ServerVariable
	for _, v := range s.VariablesField {
		vars = append(vars, v)
	}

	return vars
}

func (s Server) IDField() string {
	return "name"
}

func (s Server) ID() string {
	return s.NameField
}

func (s Server) Name() string {
	return s.NameField
}

func (s Server) HasName() bool {
	return s.NameField != ""
}

func (s Server) Description() string {
	return s.DescriptionField
}

func (s Server) HasDescription() bool {
	return s.DescriptionField != ""
}

func (s Server) URL() string {
	return s.URLField
}

func (s Server) HasURL() bool {
	return s.URLField != ""
}

func (s Server) Protocol() string {
	return s.ProtocolField
}

func (s Server) HasProtocol() bool {
	return s.ProtocolField != ""
}

type ServerVariable struct {
	Extendable
	NameField string   `mapstructure:"name"`
	Default   string   `mapstructure:"default"`
	Enum      []string `mapstructure:"enum"`
}

func (s ServerVariable) IDField() string {
	return "name"
}

func (s ServerVariable) ID() string {
	return s.NameField
}

func (s ServerVariable) Name() string {
	return s.NameField
}

func (s ServerVariable) HasName() bool {
	return s.NameField != ""
}

func (s ServerVariable) DefaultValue() string {
	return s.Default
}

func (s ServerVariable) AllowedValues() []string {
	return s.Enum
}

type Extendable struct {
	Raw map[string]interface{} `mapstructure:",remain"`
}

func (e Extendable) HasExtension(name string) bool {
	_, ok := e.Raw[name]
	return ok
}

func (e Extendable) Extension(name string) interface{} {
	return e.Raw[name]
}
