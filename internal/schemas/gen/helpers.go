package main

import (
	"github.com/hashicorp/go-version"
	"github.com/opentofu/tofu-ls/internal/registry"
)

func providerVersionSupportsOsAndArch(pVersion version.Version, providerVersions []registry.ProviderVersion, os, arch string) bool {
	for _, version := range providerVersions {
		if version.Version != pVersion.String() {
			continue
		}
		for _, platform := range version.Platforms {
			if platform.OS == os &&
				platform.Arch == arch {
				return true
			}
		}
	}

	return false
}
