package v2

type Document struct {
	Extendable
	ServersField map[string]Server `mapstructure:"servers"`
}

func (d Document) Servers() []Server {
	var servers []Server
	for _, s := range d.ServersField {
		servers = append(servers, s)
	}

	return servers
}

type Server struct {
	Extendable
	NameField      string                    `mapstructure:"name"`
	ProtocolField  string                    `mapstructure:"protocol"`
	URLField       string                    `mapstructure:"url"`
	VariablesField map[string]ServerVariable `mapstructure:"variables"`
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

func (s Server) URL() string {
	return s.URLField
}

func (s Server) HasURL() bool {
	return s.URLField != ""
}

func (s Server) Protocol() string {
	return s.ProtocolField
}

func (s Server) Variables() []ServerVariable {
	var vars []ServerVariable
	for _, v := range s.VariablesField {
		vars = append(vars, v)
	}

	return vars
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
