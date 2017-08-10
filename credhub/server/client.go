package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Provides an unauthenticated http.Client to the CredHub server
func (c *Config) Client() (*http.Client, error) {
	parsedUrl, err := url.Parse(c.ApiUrl)
	if err != nil {
		return nil, err
	}

	if parsedUrl.Scheme == "https" {
		return c.httpsClient()
	} else {
		return c.httpClient()
	}
}

func (c *Config) httpClient() (*http.Client, error) {
	return &http.Client{
		Timeout: time.Second * 45,
	}, nil
}

func (c *Config) httpsClient() (*http.Client, error) {
	client, _ := c.httpClient()

	rootCAs, err := c.certPool()
	if err != nil {
		return nil, err
	}

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:       c.InsecureSkipVerify,
			PreferServerCipherSuites: true,
			RootCAs:                  rootCAs,
		},
	}

	return client, nil
}

func (c *Config) certPool() (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, cert := range c.CaCerts {
		ok := certPool.AppendCertsFromPEM([]byte(cert))
		if !ok {
			return nil, errors.New("Invalid certificate")
		}
	}

	return certPool, nil
}