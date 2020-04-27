package drago

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const configAccessKey = "config/access"

type accessConfig struct {
	Address            string `json:"address"`
	Token              string `json:"token"`
	MaxTokenNameLength int    `json:"max_token_name_length"`
	CACert             string `json:"ca_cert"`
	ClientCert         string `json:"client_cert"`
	ClientKey          string `json:"client_key"`
}

func (b *dragoBackend) CheckAccessConfigExistence(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.readAccessConfig(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

func (b *dragoBackend) readAccessConfig(ctx context.Context, storage logical.Storage) (*accessConfig, error) {
	entry, err := storage.Get(ctx, configAccessKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	conf := &accessConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, errwrap.Wrapf("error reading drago access configuration: {{err}}", err)
	}

	return conf, nil
}

func (b *dragoBackend) HandleAccessConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readAccessConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address":               conf.Address,
			"max_token_name_length": conf.MaxTokenNameLength,
		},
	}, nil
}

func (b *dragoBackend) HandleAccessConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readAccessConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		conf = &accessConfig{}
	}

	address, ok := data.GetOk("address")
	if ok {
		conf.Address = address.(string)
	}
	token, ok := data.GetOk("token")
	if ok {
		conf.Token = token.(string)
	}
	caCert, ok := data.GetOk("ca_cert")
	if ok {
		conf.CACert = caCert.(string)
	}
	clientCert, ok := data.GetOk("client_cert")
	if ok {
		conf.ClientCert = clientCert.(string)
	}
	clientKey, ok := data.GetOk("client_key")
	if ok {
		conf.ClientKey = clientKey.(string)
	}

	conf.MaxTokenNameLength = data.Get("max_token_name_length").(int)

	entry, err := logical.StorageEntryJSON("config/access", conf)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *dragoBackend) HandleAccessConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configAccessKey); err != nil {
		return nil, err
	}
	return nil, nil
}
