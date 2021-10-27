package ezdb

import (
	"fmt"
	"net/url"
)

type ConnectionData struct {
	User     string
	Password string
	Host     string
	Db       string
	Port     string
}

func (cd ConnectionData) DSN() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s sslmode=disable",
		cd.User, cd.Password, cd.Host, cd.Port,
	)
}
func (cd ConnectionData) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cd.User,
		url.PathEscape(cd.Password),
		cd.Host,
		cd.Port,
		cd.Db,
	)
}
