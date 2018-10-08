package cdsclient

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/facebookgo/httpcontrol"

	"github.com/ovh/cds/sdk"
)

type client struct {
	isWorker   bool
	isService  bool
	isProvider bool
	HTTPClient HTTPClient
	config     Config
	name       string
	service    *sdk.Service
}

// New returns a client from a config struct
func New(c Config) Interface {
	cli := new(client)
	cli.config = c
	cli.HTTPClient = &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerifyTLS},
		},
	}
	cli.init()
	return cli
}

// NewService returns client for a service
func NewService(endpoint string, timeout time.Duration, insecureSkipVerifyTLS bool) Interface {
	conf := Config{
		Host:  endpoint,
		Retry: 2,
		InsecureSkipVerifyTLS: insecureSkipVerifyTLS,
	}
	cli := new(client)
	cli.config = conf
	cli.HTTPClient = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			MaxIdleConns:        500,
			MaxIdleConnsPerHost: 500,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: conf.InsecureSkipVerifyTLS},
		},
	}
	cli.isService = true
	cli.init()
	return cli
}

// NewWorker returns client for a worker
func NewWorker(endpoint string, name string, c HTTPClient) Interface {
	conf := Config{
		Host:  endpoint,
		Retry: 10,
	}
	cli := new(client)
	cli.config = conf

	if c == nil {
		cli.HTTPClient = &http.Client{Timeout: time.Second * 360}
	} else {
		cli.HTTPClient = c
	}

	cli.isWorker = true
	cli.name = name
	cli.init()
	return cli
}

// NewClientFromConfig returns a client from the config file
func NewClientFromConfig(r io.Reader) (Interface, error) {
	return nil, nil
}

// NewClientFromEnv returns a client from the environment variables
func NewClientFromEnv() (Interface, error) {
	return nil, nil
}

// NewProviderClient returns an implementation for ProviderClient interface
func NewProviderClient(cfg ProviderConfig) ProviderClient {
	conf := Config{
		Host:  cfg.Host,
		Retry: 2,
		Token: cfg.Token,
		User:  cfg.Name,
	}

	if cfg.RequestSecondsTimeout == 0 {
		cfg.RequestSecondsTimeout = 60
	}

	cli := new(client)
	cli.config = conf
	cli.HTTPClient = &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout:  time.Duration(cfg.RequestSecondsTimeout) * time.Second,
			MaxTries:        5,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerifyTLS},
		},
	}
	cli.isProvider = true
	cli.name = cfg.Name
	cli.init()
	return cli
}

func (c *client) init() {
	if c.isWorker {
		c.config.userAgent = sdk.WorkerAgent
	} else if c.isService {
		c.config.userAgent = sdk.ServiceAgent
	} else {
		c.config.userAgent = sdk.SDKAgent
	}

	if os.Getenv("CDS_VERBOSE") == "true" {
		c.config.Verbose = true
	}
}

func (c *client) APIURL() string {
	return c.config.Host
}
