package drago

import (
	"context"

	"github.com/coreos/etcd/client"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/seashell/vault-plugin-secrets-drago/plugin/drago/path"
	"github.com/seashell/vault-plugin-secrets-drago/plugin/drago/client"
)

type DragoBackend interface {
	HandleLeaseRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error)
	HandleLeaseUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error)
	HandleLeaseDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error)
}

type dragoBackend struct {
	*framework.Backend
}

// Factory returns a Drago backend that satisfies the logical.Backend interface
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// FactoryType is a wrapper func that allows the Factory func to specify
// the backend type for the mock backend plugin instance.
func FactoryType(backendType logical.BackendType) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b := Backend()
		b.BackendType = backendType
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// Backend returns the configured Drago backend
func Backend() *dragoBackend {
	var b dragoBackend
	b.Backend = &framework.Backend{
		Help: "Drago backend",
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"special",
			},
		},
		Paths: []*framework.Path{
			path.CreateCreds(&b),
			path.ConfigAccess(&b),
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}
	return &b
}

func (b *dragoBackend) Client(ctx context.Context, s logical.Storage) (*client.Client, error) {

	conf, err := b.readAccessConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	nomadConf := client.DefaultConfig()
	if conf != nil {
		if conf.Address != "" {
			nomadConf.Address = conf.Address
		}
		if conf.Token != "" {
			nomadConf.SecretID = conf.Token
		}
		if conf.CACert != "" {
			nomadConf.TLSConfig.CACertPEM = []byte(conf.CACert)
		}
		if conf.ClientCert != "" {
			nomadConf.TLSConfig.ClientCertPEM = []byte(conf.ClientCert)
		}
		if conf.ClientKey != "" {
			nomadConf.TLSConfig.ClientKeyPEM = []byte(conf.ClientKey)
		}
	}

	client, err := drago.NewClient(nomadConf)
	if err != nil {
		return nil, err
	}

	return client, nil
}
