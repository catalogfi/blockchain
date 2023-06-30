package btc

import (
	"time"

	"github.com/btcsuite/btcd/chaincfg"
)

const (
	// DefaultClientTimeout used by the Client.
	DefaultClientTimeout = time.Minute
	// DefaultClientTimeoutRetry used by the Client.
	DefaultClientTimeoutRetry = 5 * time.Second
	// DefaultClientHost used by the Client. This should only be used for local
	// deployments of the regression testnet.
	DefaultClientHost = "http://0.0.0.0:18443"
	// DefaultClientUser used by the Client. This is insecure, and should only
	// be used for local — or publicly accessible — deployments of the
	// multichain.
	DefaultClientUser = ""
	// DefaultClientPassword used by the Client. This is insecure, and should
	// only be used for local — or publicly accessible — deployments of the
	// multichain.
	DefaultClientPassword = ""
)

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions struct {
	Net *chaincfg.Params

	// For RPC client
	Timeout      time.Duration
	TimeoutRetry time.Duration
	Host         string
	User         string
	Password     string
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Net:          &chaincfg.RegressionNetParams,
		Timeout:      DefaultClientTimeout,
		TimeoutRetry: DefaultClientTimeoutRetry,
		Host:         DefaultClientHost,
		User:         DefaultClientUser,
		Password:     DefaultClientPassword,
	}
}

// WithNet sets the network parameter for the client.
func (opts ClientOptions) WithNet(net *chaincfg.Params) ClientOptions {
	opts.Net = net
	return opts
}

// WithTimeout sets http client timeout
func (opts ClientOptions) WithTimeout(timeout time.Duration) ClientOptions {
	opts.Timeout = timeout
	return opts
}

// WithTimeoutRetry sets http client timeout
func (opts ClientOptions) WithTimeoutRetry(t time.Duration) ClientOptions {
	opts.TimeoutRetry = t
	return opts
}

// WithHost sets the URL of the Bitcoin node.
func (opts ClientOptions) WithHost(host string) ClientOptions {
	opts.Host = host
	return opts
}

// WithUser sets the username that will be used to authenticate with the Bitcoin
// node.
func (opts ClientOptions) WithUser(user string) ClientOptions {
	opts.User = user
	return opts
}

// WithPassword sets the password that will be used to authenticate with the
// Bitcoin node.
func (opts ClientOptions) WithPassword(password string) ClientOptions {
	opts.Password = password
	return opts
}
