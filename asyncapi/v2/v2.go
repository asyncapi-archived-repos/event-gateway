//nolint:goconst
package v2

import (
	"github.com/asyncapi/event-gateway/asyncapi"
)

// Constants for Operation Types.
const (
	OperationTypePublish   asyncapi.OperationType = "publish"
	OperationTypeSubscribe asyncapi.OperationType = "subscribe"
)

type Document struct {
	Extendable
	ServersField  map[string]Server  `mapstructure:"servers"`
	ChannelsField map[string]Channel `mapstructure:"channels"`
}

func (d Document) ApplicationPublishableChannels() []asyncapi.Channel {
	return d.filterChannels(func(operation asyncapi.Operation) bool {
		return operation.IsApplicationPublishing()
	})
}

func (d Document) ApplicationPublishableMessages() []asyncapi.Message {
	return d.filterMessages(func(operation asyncapi.Operation) bool {
		return operation.IsApplicationPublishing()
	})
}

func (d Document) ApplicationPublishOperations() []asyncapi.Operation {
	return d.filterOperations(func(operation asyncapi.Operation) bool {
		return operation.IsApplicationPublishing()
	})
}

func (d Document) ApplicationSubscribableChannels() []asyncapi.Channel {
	return d.filterChannels(func(operation asyncapi.Operation) bool {
		return operation.IsApplicationSubscribing()
	})
}

func (d Document) ApplicationSubscribableMessages() []asyncapi.Message {
	return d.filterMessages(func(operation asyncapi.Operation) bool {
		return operation.IsApplicationSubscribing()
	})
}

func (d Document) ApplicationSubscribeOperations() []asyncapi.Operation {
	return d.filterOperations(func(operation asyncapi.Operation) bool {
		return operation.IsApplicationSubscribing()
	})
}

func (d Document) Channels() []asyncapi.Channel {
	var channels []asyncapi.Channel
	for _, c := range d.ChannelsField {
		channels = append(channels, c)
	}

	return channels
}

func (d Document) HasChannels() bool {
	return len(d.ChannelsField) > 0
}

func (d Document) ClientPublishableChannels() []asyncapi.Channel {
	return d.filterChannels(func(operation asyncapi.Operation) bool {
		return operation.IsClientPublishing()
	})
}

func (d Document) ClientPublishableMessages() []asyncapi.Message {
	return d.filterMessages(func(operation asyncapi.Operation) bool {
		return operation.IsClientPublishing()
	})
}

func (d Document) ClientPublishOperations() []asyncapi.Operation {
	return d.filterOperations(func(operation asyncapi.Operation) bool {
		return operation.IsClientPublishing()
	})
}

func (d Document) ClientSubscribableChannels() []asyncapi.Channel {
	return d.filterChannels(func(operation asyncapi.Operation) bool {
		return operation.IsClientSubscribing()
	})
}

func (d Document) ClientSubscribableMessages() []asyncapi.Message {
	return d.filterMessages(func(operation asyncapi.Operation) bool {
		return operation.IsClientSubscribing()
	})
}

func (d Document) ClientSubscribeOperations() []asyncapi.Operation {
	return d.filterOperations(func(operation asyncapi.Operation) bool {
		return operation.IsClientSubscribing()
	})
}

func (d Document) Messages() []asyncapi.Message {
	var messages []asyncapi.Message
	for _, c := range d.Channels() {
		messages = append(messages, c.Messages()...)
	}

	return messages
}

func (d Document) Server(name string) (asyncapi.Server, bool) {
	s, ok := d.ServersField[name]
	return s, ok
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

func (d Document) filterChannels(filter func(operation asyncapi.Operation) bool) []asyncapi.Channel {
	var channels []asyncapi.Channel
	for _, c := range d.Channels() {
		for _, o := range c.Operations() {
			if filter(o) {
				channels = append(channels, c)
			}
		}
	}

	return channels
}

func (d Document) filterMessages(filter func(operation asyncapi.Operation) bool) []asyncapi.Message {
	var messages []asyncapi.Message
	for _, c := range d.Channels() {
		for _, o := range c.Operations() {
			if filter(o) {
				messages = append(messages, c.Messages()...)
			}
		}
	}

	return messages
}

func (d Document) filterOperations(filter func(operation asyncapi.Operation) bool) []asyncapi.Operation {
	var operations []asyncapi.Operation
	for _, c := range d.Channels() {
		for _, o := range c.Operations() {
			if filter(o) {
				operations = append(operations, o)
			}
		}
	}

	return operations
}

type Channel struct {
	Extendable
	PathField        string                      `mapstructure:"path"`
	DescriptionField string                      `mapstructure:"description"`
	ParametersField  map[string]ChannelParameter `mapstructure:"parameters"`
	Subscribe        *SubscribeOperation         `mapstructure:"subscribe"`
	Publish          *PublishOperation           `mapstructure:"publish"`
}

func (c Channel) IDField() string {
	return "path"
}

func (c Channel) ID() string {
	return c.PathField
}

func (c Channel) Path() string {
	return c.PathField
}

func (c Channel) Description() string {
	return c.DescriptionField
}

func (c Channel) HasDescription() bool {
	return c.DescriptionField != ""
}

func (c Channel) Parameters() []asyncapi.ChannelParameter {
	var parameters []asyncapi.ChannelParameter
	for _, p := range c.ParametersField {
		parameters = append(parameters, p)
	}

	return parameters
}

func (c Channel) HasParameters() bool {
	return len(c.ParametersField) > 0
}

func (c Channel) Operations() []asyncapi.Operation {
	var operations []asyncapi.Operation
	if c.Publish != nil {
		operations = append(operations, c.Publish)
	}
	if c.Subscribe != nil {
		operations = append(operations, c.Subscribe)
	}

	return operations
}

func (c Channel) Messages() []asyncapi.Message {
	var messages []asyncapi.Message
	for _, o := range c.Operations() {
		messages = append(messages, o.Messages()...)
	}

	return messages
}

type ChannelParameter struct {
	Extendable
	NameField        string  `mapstructure:"name"`
	DescriptionField string  `mapstructure:"description"`
	SchemaField      *Schema `mapstructure:"schema"`
}

func (c ChannelParameter) IDField() string {
	return "name"
}

func (c ChannelParameter) ID() string {
	return c.NameField
}

func (c ChannelParameter) Name() string {
	return c.NameField
}

func (c ChannelParameter) Description() string {
	return c.DescriptionField
}

func (c ChannelParameter) Schema() asyncapi.Schema {
	return c.SchemaField
}

type SubscribeOperation struct {
	Operation
}

func (o SubscribeOperation) MapStructureDefaults() map[string]interface{} {
	return map[string]interface{}{
		"operationType": OperationTypeSubscribe,
	}
}

type PublishOperation struct {
	Operation
}

func (o PublishOperation) MapStructureDefaults() map[string]interface{} {
	return map[string]interface{}{
		"operationType": OperationTypePublish,
	}
}

type Operation struct {
	Extendable
	OperationIDField string                 `mapstructure:"operationId"`
	DescriptionField string                 `mapstructure:"description"`
	MessageField     *Message               `mapstructure:"message"`
	OperationType    asyncapi.OperationType `mapstructure:"operationType"` // set by hook
	SummaryField     string                 `mapstructure:"summary"`
}

func (o Operation) ID() string {
	return o.OperationIDField
}

func (o Operation) Description() string {
	return o.DescriptionField
}

func (o Operation) HasDescription() bool {
	return o.DescriptionField != ""
}

func (o Operation) IsApplicationPublishing() bool {
	return o.Type() == OperationTypePublish
}

func (o Operation) IsApplicationSubscribing() bool {
	return o.Type() == OperationTypeSubscribe
}

func (o Operation) IsClientPublishing() bool {
	return o.Type() == OperationTypeSubscribe
}

func (o Operation) IsClientSubscribing() bool {
	return o.Type() == OperationTypePublish
}

func (o Operation) Messages() []asyncapi.Message {
	if o.MessageField != nil {
		return []asyncapi.Message{o.MessageField}
	}
	return nil
}

func (o Operation) Type() asyncapi.OperationType {
	// TODO enum{'ClientSubscribing', 'ClientPublishing', 'ApplicationSubscribing', 'ApplicationPublishing}
	return o.OperationType
}

func (o Operation) Summary() string {
	return o.SummaryField
}

func (o Operation) HasSummary() bool {
	return o.SummaryField != ""
}

type Message struct {
	Extendable
	NameField        string  `mapstructure:"name"`
	TitleField       string  `mapstructure:"title"`
	DescriptionField string  `mapstructure:"description"`
	SummaryField     string  `mapstructure:"summary"`
	ContentTypeField string  `mapstructure:"contentType"`
	PayloadField     *Schema `mapstructure:"payload"`
}

func (m Message) UID() string {
	return m.Name() // TODO What is expected here?
}

func (m Message) Name() string {
	return m.NameField
}

func (m Message) Title() string {
	return m.TitleField
}

func (m Message) HasTitle() bool {
	return m.TitleField != ""
}

func (m Message) Description() string {
	return m.DescriptionField
}

func (m Message) HasDescription() bool {
	return m.DescriptionField != ""
}

func (m Message) Summary() string {
	return m.SummaryField
}

func (m Message) HasSummary() bool {
	return m.SummaryField != ""
}

func (m Message) ContentType() string {
	return m.ContentTypeField
}

func (m Message) Payload() asyncapi.Schema {
	return m.PayloadField
}

type Schemas map[string]*Schema

func (s Schemas) ToInterface(dst map[string]asyncapi.Schema) map[string]asyncapi.Schema {
	if len(dst) > 0 {
		return dst
	}

	if dst == nil {
		dst = make(map[string]asyncapi.Schema)
	}

	for k, v := range s {
		dst[k] = v
	}

	return dst
}

type Schema struct {
	Extendable
	AdditionalItemsField      *Schema           `mapstructure:"additionalItems"`
	AdditionalPropertiesField *Schema           `mapstructure:"additionalProperties"`
	AllOfField                []asyncapi.Schema `mapstructure:"allOf"`
	AnyOfField                []asyncapi.Schema `mapstructure:"anyOf"`
	ConstField                interface{}       `mapstructure:"const"`
	ContainsField             *Schema           `mapstructure:"contains"`
	ContentEncodingField      string            `mapstructure:"contentEncoding"`
	ContentMediaTypeField     string            `mapstructure:"contentMediaType"`
	DefaultField              interface{}       `mapstructure:"default"`
	DefinitionsField          Schemas           `mapstructure:"definitions"`
	DependenciesField         Schemas           `mapstructure:"dependencies"`
	DeprecatedField           bool              `mapstructure:"deprecated"`
	DescriptionField          string            `mapstructure:"description"`
	DiscriminatorField        string            `mapstructure:"discriminator"`
	ElseField                 *Schema           `mapstructure:"else"`
	EnumField                 []interface{}     `mapstructure:"enum"`
	ExamplesField             []interface{}     `mapstructure:"examples"`
	ExclusiveMaximumField     *float64          `mapstructure:"exclusiveMaximum"`
	ExclusiveMinimumField     *float64          `mapstructure:"exclusiveMinimum"`
	FormatField               string            `mapstructure:"format"`
	IDField                   string            `mapstructure:"$id"`
	IfField                   *Schema           `mapstructure:"if"`
	ItemsField                []asyncapi.Schema `mapstructure:"items"`
	MaximumField              *float64          `mapstructure:"maximum"`
	MaxItemsField             *float64          `mapstructure:"maxItems"`
	MaxLengthField            *float64          `mapstructure:"maxLength"`
	MaxPropertiesField        *float64          `mapstructure:"maxProperties"`
	MinimumField              *float64          `mapstructure:"minimum"`
	MinItemsField             *float64          `mapstructure:"minItems"`
	MinLengthField            *float64          `mapstructure:"minLength"`
	MinPropertiesField        *float64          `mapstructure:"minProperties"`
	MultipleOfField           *float64          `mapstructure:"multipleOf"`
	NotField                  *Schema           `mapstructure:"not"`
	OneOfField                *Schema           `mapstructure:"oneOf"`
	PatternField              string            `mapstructure:"pattern"`
	PatternPropertiesField    Schemas           `mapstructure:"patternProperties"`
	PropertiesField           Schemas           `mapstructure:"properties"`
	PropertyNamesField        *Schema           `mapstructure:"propertyNames"`
	ReadOnlyField             bool              `mapstructure:"readOnly"`
	RequiredField             string            `mapstructure:"required"`
	ThenField                 *Schema           `mapstructure:"then"`
	TitleField                string            `mapstructure:"title"`
	TypeField                 interface{}       `mapstructure:"type"` // string | []string
	UniqueItemsField          bool              `mapstructure:"uniqueItems"`
	WriteOnlyField            bool              `mapstructure:"writeOnly"`

	// cached converted map[string]asyncapi.Schema from map[string]*Schema
	propertiesFieldMap        map[string]asyncapi.Schema
	patternPropertiesFieldMap map[string]asyncapi.Schema
	DefinitionsFieldMap       map[string]asyncapi.Schema
	DependenciesFieldMap      map[string]asyncapi.Schema
}

func (s *Schema) AdditionalItems() asyncapi.Schema {
	return s.AdditionalItemsField
}

func (s *Schema) AdditionalProperties() asyncapi.Schema {
	return s.AdditionalPropertiesField
}

func (s *Schema) AllOf() []asyncapi.Schema {
	return s.AllOfField
}

func (s *Schema) AnyOf() []asyncapi.Schema {
	return s.AnyOfField
}

func (s *Schema) CircularProps() []string {
	if props, ok := s.Extension("x-parser-circular-props").([]string); ok {
		return props
	}

	return nil
}

func (s *Schema) Const() interface{} {
	return s.ConstField
}

func (s *Schema) Contains() asyncapi.Schema {
	return s.ContainsField
}

func (s *Schema) ContentEncoding() string {
	return s.ContentEncodingField
}

func (s *Schema) ContentMediaType() string {
	return s.ContentMediaTypeField
}

func (s *Schema) Default() interface{} {
	return s.DefaultField
}

func (s *Schema) Definitions() map[string]asyncapi.Schema {
	s.DefinitionsFieldMap = s.DefinitionsField.ToInterface(s.DefinitionsFieldMap)
	return s.DefinitionsFieldMap
}

func (s *Schema) Dependencies() map[string]asyncapi.Schema {
	// TODO Map[string, Schema|string[]]
	s.DependenciesFieldMap = s.DependenciesField.ToInterface(s.DependenciesFieldMap)
	return s.DependenciesFieldMap
}

func (s *Schema) Deprecated() bool {
	return s.DeprecatedField
}

func (s *Schema) Description() string {
	return s.DescriptionField
}

func (s *Schema) HasDescription() bool {
	return s.DescriptionField != ""
}

func (s *Schema) Discriminator() string {
	return s.DiscriminatorField
}

func (s *Schema) Else() asyncapi.Schema {
	return s.ElseField
}

func (s *Schema) Enum() []interface{} {
	return s.EnumField
}

func (s *Schema) Examples() []interface{} {
	return s.ExamplesField
}

func (s *Schema) ExclusiveMaximum() *float64 {
	return s.ExclusiveMaximumField
}

func (s *Schema) ExclusiveMinimum() *float64 {
	return s.ExclusiveMinimumField
}

func (s *Schema) Format() string {
	return s.FormatField
}

func (s *Schema) HasCircularProps() bool {
	return len(s.CircularProps()) > 0
}

func (s *Schema) ID() string {
	return s.IDField
}

func (s *Schema) If() asyncapi.Schema {
	return s.IfField
}

func (s *Schema) IsCircular() bool {
	if isCircular, ok := s.Extension("x-parser-circular").(bool); ok {
		return isCircular
	}

	return false
}

func (s *Schema) Items() []asyncapi.Schema {
	// TODO Schema | Schema[]
	return s.ItemsField
}

func (s *Schema) Maximum() *float64 {
	return s.MaximumField
}

func (s *Schema) MaxItems() *float64 {
	return s.MaxItemsField
}

func (s *Schema) MaxLength() *float64 {
	return s.MaxLengthField
}

func (s *Schema) MaxProperties() *float64 {
	return s.MaxPropertiesField
}

func (s *Schema) Minimum() *float64 {
	return s.MinimumField
}

func (s *Schema) MinItems() *float64 {
	return s.MinItemsField
}

func (s *Schema) MinLength() *float64 {
	return s.MinLengthField
}

func (s *Schema) MinProperties() *float64 {
	return s.MinPropertiesField
}

func (s *Schema) MultipleOf() *float64 {
	return s.MultipleOfField
}

func (s *Schema) Not() asyncapi.Schema {
	return s.NotField
}

func (s *Schema) OneOf() asyncapi.Schema {
	// TODO Schema[]
	return s.OneOfField
}

func (s *Schema) Pattern() string {
	return s.PatternField
}

func (s *Schema) PatternProperties() map[string]asyncapi.Schema {
	s.patternPropertiesFieldMap = s.PatternPropertiesField.ToInterface(s.patternPropertiesFieldMap)
	return s.patternPropertiesFieldMap
}

func (s *Schema) Properties() map[string]asyncapi.Schema {
	s.propertiesFieldMap = s.PropertiesField.ToInterface(s.propertiesFieldMap)
	return s.propertiesFieldMap
}

func (s *Schema) Property(name string) asyncapi.Schema {
	return s.PropertiesField[name]
}

func (s *Schema) PropertyNames() asyncapi.Schema {
	return s.PropertyNamesField
}

func (s *Schema) ReadOnly() bool {
	return s.ReadOnlyField
}

func (s *Schema) Required() string {
	// TODO string[]
	return s.RequiredField
}

func (s *Schema) Then() asyncapi.Schema {
	return s.ThenField
}

func (s *Schema) Title() string {
	return s.TitleField
}

func (s *Schema) Type() []string {
	if stringVal, isString := s.TypeField.(string); isString {
		return []string{stringVal}
	}

	if sliceVal, isSlice := s.TypeField.([]string); isSlice {
		return sliceVal
	}

	return nil
}

func (s *Schema) UID() string {
	if id := s.ID(); id != "" {
		return id
	}

	// This is not yet supported by parser-go
	parserGeneratedID, ok := s.Extension("x-parser-id").(string)
	if ok {
		return parserGeneratedID
	}

	return ""
}

func (s *Schema) UniqueItems() bool {
	return s.UniqueItemsField
}

func (s *Schema) WriteOnly() bool {
	return s.WriteOnlyField
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
	// TODO variable substitution if applies
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
