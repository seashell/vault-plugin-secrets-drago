package path

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type AccessConfigBackend interface {
	CheckAccessConfigExistence(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error)
	HandleAccessConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error)
	HandleAccessConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error)
	HandleAccessConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error)
}

func ConfigAccess(b AccessConfigBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",
		Fields: map[string]*framework.FieldSchema{
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Drago server address",
			},
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},
			"max_token_name_length": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Max length for name of generated Drago tokens",
			},
			"ca_cert": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `CA certificate to use when verifying Drago server certificate,
must be x509 PEM encoded.`,
			},
			"client_cert": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Client certificate used for Drago's TLS communication,
must be x509 PEM encoded and if this is set you need to also set client_key.`,
			},
			"client_key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Client key used for Drago's TLS communication,
must be x509 PEM encoded and if this is set you need to also set client_cert.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.HandleAccessConfigRead,
			logical.CreateOperation: b.HandleAccessConfigWrite,
			logical.UpdateOperation: b.HandleAccessConfigWrite,
			logical.DeleteOperation: b.HandleAccessConfigDelete,
		},

		ExistenceCheck: b.CheckAccessConfigExistence,
	}
}
