// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tfaddr "github.com/opentofu/registry-address"
)

type Provider struct {
	Addr    string `json:"addr"`
	Version string `json:"version"`
}

type ProviderVersion struct {
	Version   string                    `json:"version"`
	Platforms []ProviderVersionPlatform `json:"platforms"`
}

type ProviderVersionPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type providerVersionResponse struct {
	Versions []ProviderVersion `json:"versions"`
}

func (c Client) ListProviders() ([]Provider, error) {
	var providers []Provider
	url := fmt.Sprintf("%s/top/providers?limit=500", c.BaseURL)
	fmt.Printf("using URL %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected response: %s: %s", resp.Status, string(bodyBytes))
	}

	var response []Provider
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}
	providers = append(providers, response...)

	return providers, nil
}

func (c Client) CheckProviderVersionSupported(pAddr tfaddr.Provider) (*providerVersionResponse, error) {
	url := fmt.Sprintf("%s/v1/providers/%s/%s/versions", c.BaseURL, pAddr.Namespace, pAddr.Type)
	fmt.Printf("using URL %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected response: %s: %s", resp.Status, string(bodyBytes))
	}

	var response providerVersionResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	return &response, nil
}
