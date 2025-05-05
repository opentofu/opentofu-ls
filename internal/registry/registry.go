// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	defaultBaseURL = "https://registry.opentofu.org"
	defaultTimeout = 5 * time.Second
	tracerName     = "github.com/opentofu/opentofu-ls/internal/registry"
)

type Client struct {
	BaseURL          string
	Timeout          time.Duration
	ProviderPageSize int
	httpClient       *http.Client
}

func NewClient() Client {
	client := cleanhttp.DefaultClient()
	client.Timeout = defaultTimeout
	client.Transport = otelhttp.NewTransport(client.Transport)

	return Client{
		BaseURL:          defaultBaseURL,
		Timeout:          defaultTimeout,
		ProviderPageSize: 100,
		httpClient:       client,
	}
}
