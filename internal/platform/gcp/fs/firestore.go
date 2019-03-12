package fs

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ErrInvalidClientProvided is returned in the event that an uninitialized client is
// used to perform actions against.
var ErrInvalidClientProvided = errors.New("invalid client provided")

// FS represents firestore client
type FS struct {
	client *firestore.Client
}

// New initialize client Google Firestore.
func New(projectID string) (*FS, error) {
	ctx := context.Background()
	clientFS, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Errorf("Failed to create client: %v", projectID)
	}

	fsOut := FS{
		client: clientFS,
	}

	return &fsOut, nil
}

// Execute is used to execute Firestore commands
func (fs *FS) Execute(ctx context.Context, f func(*firestore.Client) error) error {
	ctx, span := trace.StartSpan(ctx, "platform.FS.Execute")
	defer span.End()

	if fs == nil {
		return errors.Wrap(ErrInvalidClientProvided, "fs == nil")
	}

	return f(fs.client)
}

// Close closes the Client.
func (fs *FS) Close() {
	fs.client.Close()
}
