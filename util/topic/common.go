package topic

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
)

func PublishMsg(ctx context.Context, topic string, msg interface{}) error {
	event := micro.NewEvent(topic, client.DefaultClient)
	if err := event.Publish(ctx, msg); err != nil {
		return err
	}
	return nil
}
