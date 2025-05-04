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

	return v
}
func (cd ConnectionData) ConnString() string {
	ssl := cd.SslMode
	if ssl == "" {
		ssl = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cd.User,
		url.PathEscape(cd.Password),
		cd.Host,
		cd.Port,
		cd.Db,
		ssl,
	)
}
