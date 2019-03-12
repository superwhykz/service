package handlers

import (
	"github.com/ardanlabs/service/internal/gcp/product"
	"github.com/ardanlabs/service/internal/platform/web"

	"github.com/pkg/errors"
)

// translate looks for certain error types and transforms
// them into web errors. We are losing the trace when this
// error is converted. But we don't log traces for these.
func translate(err error) error {
	switch errors.Cause(err) {
	case product.ErrNotFound:
		return web.ErrNotFound
	case product.ErrInvalidID:
		return web.ErrInvalidID
	}
	return err
}
