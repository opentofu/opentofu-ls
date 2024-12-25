// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package decoder

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfmodule "github.com/hashicorp/terraform-schema/module"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/opentofu/opentofu-ls/internal/features/modules/state"
)

func functionsForModule(mod *state.ModuleRecord, stateReader CombinedReader) (map[string]schema.FunctionSignature, error) {
	resolvedVersion := tfschema.ResolveVersion(stateReader.TerraformVersion(mod.Path()), mod.Meta.CoreRequirements)
	sm := tfschema.NewFunctionsMerger(mustFunctionsForVersion(resolvedVersion))
	sm.SetTerraformVersion(resolvedVersion)
	sm.SetStateReader(stateReader)

	meta := &tfmodule.Meta{
		Path:                 mod.Path(),
		ProviderRequirements: mod.Meta.ProviderRequirements,
		ProviderReferences:   mod.Meta.ProviderReferences,
	}

	return sm.FunctionsForModule(meta)
}

func mustFunctionsForVersion(v *version.Version) map[string]schema.FunctionSignature {
	s, err := tfschema.FunctionsForVersion(v)
	if err != nil {
		// this should never happen
		panic(err)
	}
	return s
}
