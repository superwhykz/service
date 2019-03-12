package ds

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ErrInvalidClientProvided is returned in the event that an uninitialized client is
// used to perform actions against.
var ErrInvalidClientProvided = errors.New("invalid client provided")

// DS represents datastore session.
type DS struct {
	client *datastore.Client
}

// New initialize client Google Datastore.
func New(projectID string) (*DS, error) {
	ctx := context.Background()
	clientDS, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Errorf("Failed to create client: %v", projectID)
	}

	dsOut := DS{
		client: clientDS,
	}

	return &dsOut, nil

}

// Execute is used to execute Datastore commands
func (ds *DS) Execute(ctx context.Context, f func(*datastore.Client) error) error {
	ctx, span := trace.StartSpan(ctx, "platform.DS.Execute")
	defer span.End()

	if ds == nil {
		return errors.Wrap(ErrInvalidClientProvided, "ds == nil")
	}

	return f(ds.client)

}

// Close closes the Client.
func (ds *DS) Close() {
	ds.client.Close()
}
