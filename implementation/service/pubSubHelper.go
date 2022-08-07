package service

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
)

func Publish(projectId string, topicID string, msg []byte) error {
	ctx := context.Background()
	//client, err := pubsub.NewClient(ctx, projectID)
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %v", err)
	}
	defer client.Close()
	t := client.Topic(topicID)

	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		//panic(err)
		log.Error().Err(err).Msg("")
		return err
	}
	//fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	log.Info().Msg("Published a message; msg ID: " + id)
	return nil
}
