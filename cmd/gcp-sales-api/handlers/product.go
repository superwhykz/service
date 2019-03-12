package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/ardanlabs/service/internal/gcp/product"
	"github.com/ardanlabs/service/internal/platform/gcp/ds"
	"github.com/ardanlabs/service/internal/platform/web"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Product represents the Product API method handler set.
type Product struct {
	ClientDS *ds.DS
}

// List returns all the existing products in the system.
func (p *Product) List(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Product.List")
	defer span.End()

	products, err := product.List(ctx, p.ClientDS)
	if err = translate(err); err != nil {
		return errors.Wrap(err, "")
	}

	web.Respond(ctx, log, w, products, http.StatusOK)
	return nil
}

// Retrieve returns the specified product from the system.
func (p *Product) Retrieve(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Product.Retrieve")
	defer span.End()

	prod, err := product.Retrieve(ctx, p.ClientDS, params["id"])
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "ID: %s", params["id"])
	}

	web.Respond(ctx, log, w, prod, http.StatusOK)
	return nil
}

// Create inserts a new product into the system.
func (p *Product) Create(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Product.Create")
	defer span.End()

	v := ctx.Value(web.KeyValues).(*web.Values)

	var np product.NewProduct
	if err := web.Unmarshal(r.Body, &np); err != nil {
		return errors.Wrap(err, "")
	}

	nUsr, err := product.Create(ctx, p.ClientDS, &np, v.Now)
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "Product: %+v", &np)
	}

	web.Respond(ctx, log, w, nUsr, http.StatusCreated)
	return nil
}

// Update updates the specified product in the system.
func (p *Product) Update(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Product.Update")
	defer span.End()

	v := ctx.Value(web.KeyValues).(*web.Values)

	var up product.UpdateProduct
	if err := web.Unmarshal(r.Body, &up); err != nil {
		return errors.Wrap(err, "")
	}

	err := product.Update(ctx, p.ClientDS, params["id"], up, v.Now)
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "ID: %s Update: %+v", params["id"], up)
	}

	web.Respond(ctx, log, w, nil, http.StatusNoContent)
	return nil
}

// Delete removes the specified product from the system.
func (p *Product) Delete(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Product.Delete")
	defer span.End()

	err := product.Delete(ctx, p.ClientDS, params["id"])
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "Id: %s", params["id"])
	}

	web.Respond(ctx, log, w, nil, http.StatusNoContent)
	return nil
}
