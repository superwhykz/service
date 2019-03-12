package ps

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ErrInvalidClientProvided is returned in the event that an uninitialized client is
// used to perform actions against.
var ErrInvalidClientProvided = errors.New("invalid client provided")

// PS represents Pub/Sub client
type PS struct {
	client *pubsub.Client
}

// New initialize client Google PubSub.
func New(projectID string) (*PS, error) {
	ctx := context.Background()
	clientPS, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Errorf("Failed to create client: %v", projectID)
	}

	psOut := PS{
		client: clientPS,
	}

	return &psOut, nil

}

// Execute is used to execute Publish/Subscribe commands
func (ps *PS) Execute(ctx context.Context, f func(*pubsub.Client) error) error {
	ctx, span := trace.StartSpan(ctx, "platform.PS.Execute")
	defer span.End()

	if ps == nil {
		return errors.Wrap(ErrInvalidClientProvided, "ps == nil")
	}

	return f(ps.client)

}

// Close closes the Client.
func (ps *PS) Close() {
	ps.client.Close()
}
