package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asyncapi/event-gateway/asyncapi"
)

func TestMessage_OneOfPayload_MultipleMessages(t *testing.T) {
	oneOfSchemas := []asyncapi.Schema{
		&Schema{
			PropertiesField: map[string]*Schema{
				"schemaOnefieldOne": {},
				"schemaOnefieldTwo": {},
			},
		},
		&Schema{
			PropertiesField: map[string]*Schema{
				"schemaTwofieldOne": {},
				"schemaTWofieldTwo": {},
			},
		},
	}
	msg := &Message{
		PayloadField: &Schema{
			OneOfField: oneOfSchemas,
		},
	}

	msgs := NewSubscribeOperation(msg).Messages()
	require.Len(t, msgs, 2)
	assert.Equal(t, oneOfSchemas[0], msgs[0].Payload())
	assert.Equal(t, oneOfSchemas[1], msgs[1].Payload())
}

func TestMessage_PlainPayload_OneMessage(t *testing.T) {
	msg := &Message{
		PayloadField: &Schema{
			PropertiesField: map[string]*Schema{
				"schemaFieldOne":      {},
				"schemaFieldfieldTwo": {},
			},
		},
	}

	msgs := NewSubscribeOperation(msg).Messages()
	require.Len(t, msgs, 1)
	assert.Equal(t, msg, msgs[0])
}

func TestSchema_AdditionalProperties(t *testing.T) {
	schema := &Schema{}
	assert.Nil(t, schema.AdditionalProperties())

	schema = &Schema{AdditionalPropertiesField: false}
	assert.True(t, schema.AdditionalProperties().IsFalse())
	assert.False(t, schema.AdditionalProperties().IsSchema())
	assert.Nil(t, schema.AdditionalProperties().Schema())

	field := &Schema{TypeField: "string"}
	schema = &Schema{AdditionalPropertiesField: field}
	assert.False(t, schema.AdditionalProperties().IsFalse())
	assert.True(t, schema.AdditionalProperties().IsSchema())
	assert.Equal(t, field, schema.AdditionalProperties().Schema())
}

func TestSchema_AdditionalItems(t *testing.T) {
	schema := &Schema{}
	assert.Nil(t, schema.AdditionalItems())

	schema = &Schema{AdditionalItemsField: false}
	assert.True(t, schema.AdditionalItems().IsFalse())
	assert.False(t, schema.AdditionalItems().IsSchema())
	assert.Nil(t, schema.AdditionalItems().Schema())

	field := &Schema{TypeField: "string"}
	schema = &Schema{AdditionalItemsField: field}
	assert.False(t, schema.AdditionalItems().IsFalse())
	assert.True(t, schema.AdditionalItems().IsSchema())
	assert.Equal(t, field, schema.AdditionalItems().Schema())
}
