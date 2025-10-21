package ezdb

import (
	"fmt"
	"os"
	"strings"
)

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
	SslKey      string
	SslCertKey  string
}

func (e EnvConnector) Data() (ConnectionData, error) {
	// Be explicit, so we dont use other(default) configs in shared envs:

	// Validate required keys are configured
	var missingKeys []string
	if e.UserKey == "" {
		missingKeys = append(missingKeys, "UserKey")
	}
	if e.HostKey == "" {
		missingKeys = append(missingKeys, "HostKey")
	}
	if e.DbKey == "" {
		missingKeys = append(missingKeys, "DbKey")
	}
	if e.PortKey == "" {
		missingKeys = append(missingKeys, "PortKey")
	}

	if len(missingKeys) > 0 {
		return ConnectionData{}, fmt.Errorf("EnvConnector missing required key configuration: %s", strings.Join(missingKeys, ", "))
	}

	// Get values from environment
	user := os.Getenv(e.UserKey)
	host := os.Getenv(e.HostKey)
	db := os.Getenv(e.DbKey)
	port := os.Getenv(e.PortKey)

	// Validate required environment variables are set
	var missingEnvVars []string
	if user == "" {
		missingEnvVars = append(missingEnvVars, e.UserKey)
	}
	if host == "" {
		missingEnvVars = append(missingEnvVars, e.HostKey)
	}
	if db == "" {
		missingEnvVars = append(missingEnvVars, e.DbKey)
	}
	if port == "" {
		missingEnvVars = append(missingEnvVars, e.PortKey)
	}

	if len(missingEnvVars) > 0 {
		return ConnectionData{}, fmt.Errorf("required environment variables not set: %s", strings.Join(missingEnvVars, ", "))
	}

	// Optional fields (Password might use peer/trust auth, SSL is optional)
	return ConnectionData{
		User:        user,
		Password:    os.Getenv(e.PasswordKey),
		Host:        host,
		Db:          db,
		Port:        port,
		SslMode:     os.Getenv(e.SslKey),
		SslRootCert: os.Getenv(e.SslCertKey),
	}, nil
}

func NewDefaultEnvConnector() EnvConnector {
	return EnvConnector{
		UserKey:     "PGUSER",
		PasswordKey: "PGPASSWORD",
		HostKey:     "PGHOST",
		DbKey:       "PGDATABASE",
		PortKey:     "PGPORT",
		SslKey:      "PGSSLMODE",
		SslCertKey:  "PGSSLROOTCERT",
	}
}
