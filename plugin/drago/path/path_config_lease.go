package path

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/seashell/vault-plugin-secrets-drago/plugin/drago"
)

func pathConfigLease(b drago.DragoBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config/lease",
		Fields: map[string]*framework.FieldSchema{
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "Duration before which the issued token needs renewal",
			},
			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: `Duration after which the issued token should not be allowed to be renewed`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.HandleLeaseRead,
			logical.UpdateOperation: b.HandleLeaseUpdate,
			logical.DeleteOperation: b.HandleLeaseDelete,
		},

		HelpSynopsis: "Configure the lease parameters for generated tokens",
		HelpDescription: `Sets the ttl and max_ttl values for the secrets to be issued by this backend.
		Both ttl and max_ttl takes in an integer number of seconds as input as well as
		inputs like "1h".
		`,
	}
}
