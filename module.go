package caddyfss3

import (
	"errors"
	"io/fs"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/geektheripper/vast-dsn/s3_dsn"
	"github.com/jszwec/s3fs/v2"
)

func init() {
	caddy.RegisterModule(FS{})
}

// Interface guards
var (
	_ fs.StatFS             = (*FS)(nil)
	_ caddyfile.Unmarshaler = (*FS)(nil)
)

// FS is a Caddy virtual filesystem module for AWS S3 (and compatible) object store.
type FS struct {
	fs.StatFS `json:"-"`

	// The dsn of the S3 bucket. if set, the bucket and region will be ignored.
	DSN string `json:"dsn,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (FS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.fs.s3",
		New: func() caddy.Module { return new(FS) },
	}
}

func (fs *FS) Provision(ctx caddy.Context) error {
	client, bucket, err := s3_dsn.NewS3Bucket(fs.DSN)
	if err != nil {
		return err
	}

	fs.StatFS = s3fs.New(client, bucket, s3fs.WithReadSeeker)
	return nil
}

// UnmarshalCaddyfile unmarshals a caddyfile.
func (fs *FS) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() { // skip block beginning
		return d.ArgErr()
	}

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "dsn":
			if !d.AllArgs(&fs.DSN) {
				return d.ArgErr()
			}
		default:
			return d.Errf("%s not a valid caddy.fs.s3 option", d.Val())
		}
	}

	if fs.DSN == "" {
		return errors.New("dsn must be set")
	}

	_, _, err := s3_dsn.NewS3Bucket(fs.DSN)
	if err != nil {
		return err
	}

	return nil
}
