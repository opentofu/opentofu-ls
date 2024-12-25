// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"

	lsctx "github.com/opentofu/opentofu-ls/internal/context"
	"github.com/opentofu/opentofu-ls/internal/langserver/cmd"
	"github.com/opentofu/opentofu-ls/internal/langserver/handlers/command"
	ilsp "github.com/opentofu/opentofu-ls/internal/lsp"
	lsp "github.com/opentofu/opentofu-ls/internal/protocol"
)

func (svc *service) TextDocumentDidSave(ctx context.Context, params lsp.DidSaveTextDocumentParams) error {
	expFeatures, err := lsctx.ExperimentalFeatures(ctx)
	if err != nil {
		return err
	}
	if !expFeatures.ValidateOnSave {
		return nil
	}

	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)

	cmdHandler := &command.CmdHandler{
		StateStore: svc.stateStore,
	}
	_, err = cmdHandler.TerraformValidateHandler(ctx, cmd.CommandArgs{
		"uri": dh.Dir.URI,
	})

	return err
}
