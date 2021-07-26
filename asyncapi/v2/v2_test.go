package v2

import (
	"testing"

	"github.com/asyncapi/event-gateway/asyncapi"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestMessage_OneOf(t *testing.T) {
	expectedMessages := []*Message{
		generateTestMessage(),
		generateTestMessage(),
		generateTestMessage(),
	}

	op := NewSubscribeOperation(expectedMessages...)
	assert.IsType(t, Messages{}, op.MessageField)
	assert.EqualValues(t, []asyncapi.Message{expectedMessages[0], expectedMessages[1], expectedMessages[2]}, op.Messages()) // Go lack of covariance :/
}

func TestMessage_PlainPayload_OneMessage(t *testing.T) {
	expectedMsg := generateTestMessage()

	op := NewSubscribeOperation(expectedMsg)
	assert.IsType(t, Messages{}, op.MessageField)
	msgs := op.Messages()
	require.Len(t, msgs, 1)
	assert.Equal(t, expectedMsg, msgs[0])
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

func generateTestMessage() *Message {
	return &Message{
		PayloadField: &Schema{
			TypeField: "object",
			PropertiesField: Schemas{
				"fieldOne": &Schema{
					TypeField: "string",
				},
			},
		},
	}
}
