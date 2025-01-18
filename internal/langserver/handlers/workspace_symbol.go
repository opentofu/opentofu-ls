// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"

	ilsp "github.com/opentofu/opentofu-ls/internal/lsp"
	lsp "github.com/opentofu/opentofu-ls/internal/protocol"
)

func (svc *service) WorkspaceSymbol(ctx context.Context, params lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	cc, err := ilsp.ClientCapabilities(ctx)
	if err != nil {
		return nil, err
	}

	// TODO? maybe kick off indexing of the whole workspace here, use ProgressToken
	symbols, err := svc.decoder.Symbols(ctx, params.Query)
	if err != nil {
		return nil, err
	}

	return ilsp.WorkspaceSymbols(symbols, cc.Workspace.Symbol), nil
}
