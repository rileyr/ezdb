package ezdb

import (
	"fmt"
	"net/url"
)

type ConnectionData struct {
	User        string
	Password    string
	Host        string
	Db          string
	Port        string
	SslMode     string
	SslRootCert string
	// SearchPath, when set, pins the connection's search_path. This matters for
	// migrations: golang-migrate issues unqualified statements against the
	// schema_migrations table, so the table it targets is whatever the
	// connection's search_path resolves to first. Setting it explicitly avoids
	// depending on PGOPTIONS/role defaults reaching the process.
	SearchPath string
}

func (cd ConnectionData) DSN() string {
	ssl := cd.SslMode
	if ssl == "" {
		ssl = "disable"
	}
	v := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s sslmode=%s",
		cd.User, cd.Password, cd.Host, cd.Port, ssl,
	)

	if cd.SslRootCert != "" {
		v = v + " sslrootcert=" + cd.SslRootCert
	}

	if cd.SearchPath != "" {
		v = v + " search_path=" + cd.SearchPath
	}

	return v
}
func (cd ConnectionData) ConnString() string {
	ssl := cd.SslMode
	if ssl == "" {
		ssl = "disable"
	}

	v := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cd.User,
		url.PathEscape(cd.Password),
		cd.Host,
		cd.Port,
		cd.Db,
		ssl,
	)

	// lib/pq forwards unknown query params as run-time parameters in the startup
	// packet; search_path is a valid startup GUC, so this pins it for every
	// connection without relying on the ambient PGOPTIONS env var.
	if cd.SearchPath != "" {
		v = v + "&search_path=" + url.QueryEscape(cd.SearchPath)
	}

	return v
}
