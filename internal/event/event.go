package event

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApplicationEventPublisher,
)

// alias
var (
	NewMessage = message.NewMessage
	NewUUID    = watermill.NewUUID
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, message *message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
}

func Marshal(v any) []byte {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return buf
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
