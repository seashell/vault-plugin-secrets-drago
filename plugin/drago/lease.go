package drago

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const leaseConfigKey = "config/lease"

type configLease struct {
	TTL    time.Duration `json:"ttl" mapstructure:"ttl"`
	MaxTTL time.Duration `json:"max_ttl" mapstructure:"max_ttl"`
}

// Sets the lease configuration parameters
func (b *dragoBackend) HandleLeaseUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/lease", &configLease{
		TTL:    time.Second * time.Duration(d.Get("ttl").(int)),
		MaxTTL: time.Second * time.Duration(d.Get("max_ttl").(int)),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *dragoBackend) HandleLeaseDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, leaseConfigKey); err != nil {
		return nil, err
	}

	return nil, nil
}

// Returns the lease configuration parameters
func (b *dragoBackend) HandleLeaseRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	lease, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"ttl":     int64(lease.TTL.Seconds()),
			"max_ttl": int64(lease.MaxTTL.Seconds()),
		},
	}, nil
}

// Lease returns the lease information
func (b *dragoBackend) LeaseConfig(ctx context.Context, s logical.Storage) (*configLease, error) {
	entry, err := s.Get(ctx, leaseConfigKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result configLease
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
