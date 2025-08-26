package event

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-kratos/kratos/v2/transport"
)

var _ transport.Server = (*ApplicationEventPublisher)(nil)

type ApplicationEventPublisher struct {
	publisher *gochannel.GoChannel
}

func NewApplicationEventPublisher() *ApplicationEventPublisher {
	return &ApplicationEventPublisher{
		publisher: gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{}),
	}
}

// Start implements transport.Server.
func (pub *ApplicationEventPublisher) Start(context.Context) error {
	return nil
}

// Stop implements transport.Server.
func (pub *ApplicationEventPublisher) Stop(context.Context) error {
	return pub.publisher.Close()
}

func (pub *ApplicationEventPublisher) Publish(ctx context.Context, topic string, payload *message.Message) {
	_ = pub.publisher.Publish(topic, payload)
}

func (pub *ApplicationEventPublisher) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return pub.publisher.Subscribe(ctx, topic)
}
