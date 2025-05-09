package main

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/opentofu/tofu-ls/internal/registry"
)

func TestProviderSupport(t *testing.T) {
	ver, err := version.NewVersion("6.0.0-beta1")

	if err != nil {
		t.Fatal(err)
	}

	pVersions := []registry.ProviderVersion{
		{
			Version: "6.0.0-beta1",
			Platforms: []registry.ProviderVersionPlatform{
				{
					OS:   "darwin",
					Arch: "amd64",
				},
				{
					OS:   "linux",
					Arch: "386",
				},
				{
					OS:   "windows",
					Arch: "386",
				},
			},
		},
	}

	result := providerVersionSupportsOsAndArch(*ver, pVersions, "linux", "386")
	if !result {
		t.Fatalf("expecting linux 386 to be a supported version")
	}

}
