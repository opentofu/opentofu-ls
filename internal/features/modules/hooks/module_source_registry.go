// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"context"
	"strings"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/zclconf/go-cty/cty"
)

type RegistryModule struct {
	FullName    string `json:"full-name"`
	Description string `json:"description"`
}

func (h *Hooks) RegistryModuleSources(ctx context.Context, value cty.Value) ([]decoder.Candidate, error) {
	candidates := make([]decoder.Candidate, 0)
	prefix := value.AsString()

	if strings.HasPrefix(prefix, ".") {
		// We're likely dealing with a local module source here; no need to search the registry
		// A search for "." will not return any results
		return candidates, nil
	}

	// TODO: Because we disabled the old logic that was fetching modules from Algolia, we'll need to figure out our own solution for modules auto completion
	return candidates, nil
}
