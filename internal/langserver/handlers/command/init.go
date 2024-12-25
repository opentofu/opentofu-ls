// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2"
	"github.com/opentofu/opentofu-ls/internal/document"
	"github.com/opentofu/opentofu-ls/internal/langserver/cmd"
	"github.com/opentofu/opentofu-ls/internal/langserver/errors"
	"github.com/opentofu/opentofu-ls/internal/langserver/progress"
	"github.com/opentofu/opentofu-ls/internal/terraform/module"
	"github.com/opentofu/opentofu-ls/internal/uri"
)

func (h *CmdHandler) TerraformInitHandler(ctx context.Context, args cmd.CommandArgs) (interface{}, error) {
	dirUri, ok := args.GetString("uri")
	if !ok || dirUri == "" {
		return nil, fmt.Errorf("%w: expected module uri argument to be set", jrpc2.InvalidParams.Err())
	}

	if !uri.IsURIValid(dirUri) {
		return nil, fmt.Errorf("URI %q is not valid", dirUri)
	}

	dirHandle := document.DirHandleFromURI(dirUri)
	tfExec, err := module.TerraformExecutorForModule(ctx, dirHandle.Path())
	if err != nil {
		return nil, errors.EnrichTfExecError(err)
	}

	progress.Begin(ctx, "Initializing")
	defer func() {
		progress.End(ctx, "Finished")
	}()

	progress.Report(ctx, "Running terraform init ...")
	err = tfExec.Init(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
