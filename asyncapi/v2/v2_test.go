package v2

import (
	"testing"

	"github.com/asyncapi/event-gateway/asyncapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestOperationShortcuts(t *testing.T) {
	doc := &Document{
		ChannelsField: map[string]Channel{
			"test-channel-1": {
				Subscribe: NewSubscribeOperation(generateTestMessageWithName("test-channel-1_subscribe_message")),
			},
			"test-channel-2": {
				Publish: NewPublishOperation(generateTestMessageWithName("test-channel-2_publish_message")),
			},
			"test-channel-3": {
				Subscribe: NewSubscribeOperation(generateTestMessageWithName("test-channel-3_subscribe_message")),
				Publish:   NewPublishOperation(generateTestMessageWithName("test-channel-3_publish_message")),
			},
		},
	}

	// Application publishes messages in channels where the Client subscribes to. Subscribe is the right type of Operation.
	assert.ElementsMatch(t, []asyncapi.Channel{doc.ChannelsField["test-channel-1"], doc.ChannelsField["test-channel-3"]}, doc.ApplicationPublishableChannels())

	// Application subscribes to messages in channels where the Client publishes to. Publish is the right type of Operation.
	assert.ElementsMatch(t, []asyncapi.Channel{doc.ChannelsField["test-channel-2"], doc.ChannelsField["test-channel-3"]}, doc.ApplicationSubscribableChannels())

	// Client publishes messages in channels where the Application subscribes to. Publish is the right type of Operation.
	assert.ElementsMatch(t, []asyncapi.Channel{doc.ChannelsField["test-channel-2"], doc.ChannelsField["test-channel-3"]}, doc.ClientPublishableChannels())

	// Client subscribe to messages in channels where the Application publishes to. Subscribe is the right type of Operation.
	assert.ElementsMatch(t, []asyncapi.Channel{doc.ChannelsField["test-channel-1"], doc.ChannelsField["test-channel-3"]}, doc.ClientSubscribableChannels())

	// Application publishes messages through operations of type Subscribe.
	assert.ElementsMatch(t, []asyncapi.Operation{doc.ChannelsField["test-channel-1"].Subscribe, doc.ChannelsField["test-channel-3"].Subscribe}, doc.ApplicationPublishOperations())

	// Application publishes messages through operations of type Subscribe.
	assert.ElementsMatch(t, []asyncapi.Operation{doc.ChannelsField["test-channel-1"].Subscribe, doc.ChannelsField["test-channel-3"].Subscribe}, doc.ApplicationPublishOperations())

	// Application subscribes to messages through operations of type Publish.
	assert.ElementsMatch(t, []asyncapi.Operation{doc.ChannelsField["test-channel-2"].Publish, doc.ChannelsField["test-channel-3"].Publish}, doc.ApplicationSubscribeOperations())

	// Client publishes messages through operations of type Publish.
	assert.ElementsMatch(t, []asyncapi.Operation{doc.ChannelsField["test-channel-2"].Publish, doc.ChannelsField["test-channel-3"].Publish}, doc.ClientPublishOperations())

	// Client subscribes to messages through operations of type Subscribe.
	assert.ElementsMatch(t, []asyncapi.Operation{doc.ChannelsField["test-channel-1"].Subscribe, doc.ChannelsField["test-channel-3"].Subscribe}, doc.ClientSubscribeOperations())

	// Application publishes messages to channels by using Subscribe operations.
	assert.ElementsMatch(t, []asyncapi.Message{doc.ChannelsField["test-channel-1"].Subscribe.Messages()[0], doc.ChannelsField["test-channel-3"].Subscribe.Messages()[0]}, doc.ApplicationPublishableMessages())

	// Application subscribes to messages by using Publish operations.
	assert.ElementsMatch(t, []asyncapi.Message{doc.ChannelsField["test-channel-2"].Publish.Messages()[0], doc.ChannelsField["test-channel-3"].Publish.Messages()[0]}, doc.ApplicationSubscribableMessages())

	// Client publishes messages to channels by using Publish operations.
	assert.ElementsMatch(t, []asyncapi.Message{doc.ChannelsField["test-channel-2"].Publish.Messages()[0], doc.ChannelsField["test-channel-3"].Publish.Messages()[0]}, doc.ClientPublishableMessages())

	// Client subscribes to messages by using Subscribe operations.
	assert.ElementsMatch(t, []asyncapi.Message{doc.ChannelsField["test-channel-1"].Subscribe.Messages()[0], doc.ChannelsField["test-channel-3"].Subscribe.Messages()[0]}, doc.ClientSubscribableMessages())
}

func generateTestMessageWithName(name string) *Message {
	m := generateTestMessage()
	m.NameField = name

	return m
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
