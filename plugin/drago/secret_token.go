package drago

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	SecretTokenType = "token"
)

func secretToken(b *dragoBackend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Request token",
			},
		},

		Renew:  b.RenewSecretToken,
		Revoke: b.RevokeSecretToken,
	}
}

func (b *dragoBackend) RenewSecretToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lease, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}
	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = lease.TTL
	resp.Secret.MaxTTL = lease.MaxTTL
	return resp, nil
}

func (b *dragoBackend) RevokeSecretToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := b.client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, fmt.Errorf("error getting Drago client")
	}

	accessorIDRaw, ok := req.Secret.InternalData["accessor_id"]
	if !ok {
		return nil, fmt.Errorf("accessor_id is missing on the lease")
	}
	accessorID, ok := accessorIDRaw.(string)
	if !ok {
		return nil, errors.New("unable to convert accessor_id")
	}
	_, err = c.ACLTokens().Delete(accessorID, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
