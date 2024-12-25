// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package progress

import (
	"context"

	"github.com/creachadair/jrpc2"
	lsctx "github.com/opentofu/opentofu-ls/internal/context"
	lsp "github.com/opentofu/opentofu-ls/internal/protocol"
)

func Begin(ctx context.Context, title string) error {
	token, ok := lsctx.ProgressToken(ctx)
	if !ok {
		return nil
	}

	return jrpc2.ServerFromContext(ctx).Notify(ctx, "$/progress", lsp.ProgressParams{
		Token: token,
		Value: lsp.WorkDoneProgressBegin{
			Kind:  "begin",
			Title: title,
		},
	})
}

func Report(ctx context.Context, message string) error {
	token, ok := lsctx.ProgressToken(ctx)
	if !ok {
		return nil
	}

	return jrpc2.ServerFromContext(ctx).Notify(ctx, "$/progress", lsp.ProgressParams{
		Token: token,
		Value: lsp.WorkDoneProgressReport{
			Kind:    "report",
			Message: message,
		},
	})
}

func End(ctx context.Context, message string) error {
	token, ok := lsctx.ProgressToken(ctx)
	if !ok {
		return nil
	}

	return jrpc2.ServerFromContext(ctx).Notify(ctx, "$/progress", lsp.ProgressParams{
		Token: token,
		Value: lsp.WorkDoneProgressEnd{
			Kind:    "end",
			Message: message,
		},
	})
}
