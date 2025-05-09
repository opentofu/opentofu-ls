// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	hcinstall "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
)

func TestExec_timeout(t *testing.T) {
	// This test is known to fail under '-race'
	// and similar race conditions are reproducible upstream
	// See https://github.com/hashicorp/terraform-exec/issues/129
	t.Skip("upstream implementation prone to race conditions")

	e := newExecutor(t)
	timeout := 1 * time.Millisecond
	e.SetTimeout(timeout)

	expectedErr := ExecTimeoutError("Version", timeout)

	_, _, err := e.Version(context.Background())
	if err != nil {
		if errors.Is(err, expectedErr) {
			return
		}

		t.Fatalf("errors don't match.\nexpected: %#v\ngiven:    %#v\n",
			expectedErr, err)
	}

	t.Fatalf("expected timeout error: %#v, given: %#v", expectedErr, err)
}

func TestExec_cancel(t *testing.T) {
	e := newExecutor(t)

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()

	expectedErr := ExecCanceledError("Version")

	_, _, err := e.Version(ctx)
	if err != nil {
		if errors.Is(err, expectedErr) {
			return
		}

		t.Fatalf("errors don't match.\nexpected: %#v\ngiven:    %#v\n",
			expectedErr, err)
	}

	t.Fatalf("expected cancel error: %#v, given: %#v", expectedErr, err)
}

func newExecutor(t *testing.T) TerraformExecutor {
	ctx := context.Background()
	workDir := TempDir(t)
	installDir := filepath.Join(workDir, "hcinstall")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(installDir); err != nil {
			t.Fatal(err)
		}
	})

	i := hcinstall.NewInstaller()

	execPath, err := i.Ensure(ctx, []src.Source{
		&releases.LatestVersion{
			Product:    product.Terraform,
			InstallDir: installDir,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := i.Remove(ctx); err != nil {
			t.Fatal(err)
		}
	})

	e, err := NewExecutor(workDir, execPath)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TempDir(t *testing.T) string {
	tmpDir := filepath.Join(os.TempDir(), "tofu-ls", t.Name())

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
	})
	return tmpDir
}
