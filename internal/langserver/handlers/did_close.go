// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"

	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

func (svc *service) TextDocumentDidClose(ctx context.Context, params lsp.DidCloseTextDocumentParams) error {
	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)
	return svc.stateStore.DocumentStore.CloseDocument(dh)
}
