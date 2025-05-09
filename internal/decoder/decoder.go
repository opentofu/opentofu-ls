// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package decoder

import (
	"context"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/opentofu/tofu-ls/internal/codelens"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

func DecoderContext(ctx context.Context) decoder.DecoderContext {
	dCtx := decoder.NewDecoderContext()

	cc, err := ilsp.ClientCapabilities(ctx)
	if err == nil {
		cmdId, ok := lsp.ExperimentalClientCapabilities(cc.Experimental).ShowReferencesCommandId()
		if ok {
			dCtx.CodeLenses = append(dCtx.CodeLenses, codelens.ReferenceCount(cmdId))
		}
	}

	return dCtx
}
