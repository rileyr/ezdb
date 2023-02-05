package ezdb

import "os"

// Connector defines the interface for obtaining database connection information.
type Connector interface {
	Data() (ConnectionData, error)
}

// ConnectorFunc is a helper type for function-only implementations of Connector
type ConnectorFunc func() (ConnectionData, error)

func (cf ConnectorFunc) Data() (ConnectionData, error) {
	return cf()
}

var _ Connector = &EnvConnector{}

// EnvConnector implements Connector by pulling values from environment variables.
type EnvConnector struct {
	UserKey     string
	PasswordKey string
	HostKey     string
	DbKey       string
	PortKey     string
}

func (e EnvConnector) Data() (ConnectionData, error) {
	return ConnectionData{
		User:     os.Getenv(e.UserKey),
		Password: os.Getenv(e.PasswordKey),
		Host:     os.Getenv(e.HostKey),
		Db:       os.Getenv(e.DbKey),
		Port:     os.Getenv(e.PortKey),
	}, nil
}

func NewDefaultEnvConnector() EnvConnector {
	return EnvConnector{
		UserKey:     "PGUSER",
		PasswordKey: "PGPASSWORD",
		HostKey:     "PGHOST",
		DbKey:       "PGDATABASE",
		PortKey:     "PGPORT",
	}
}
