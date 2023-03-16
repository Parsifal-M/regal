package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/open-policy-agent/opa/loader"
	rio "github.com/styrainc/regal/internal/io"
	"github.com/styrainc/regal/pkg/linter"
)

// Note: this will bundle the tests as well, but since that has negligible impact on the size of the binary,
// it's preferable to filter them out from the bundle than to e.g. create a separate directory for tests
//go:embed all:bundle
var bundle embed.FS

func main() {
	// Remove date and time from any `log.*` calls, as that doesn't add much of value here
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatal("At least one file or directory must be provided for linting")
	}

	// Create new fs from root of bundle, do avoid having to deal with
	// "bundle" in paths (i.e. `data.bundle.rega`)
	bfs, err := fs.Sub(bundle, "bundle")
	if err != nil {
		log.Fatal(err)
	}

	regalRules := rio.MustLoadRegalBundle(bfs)
	policies, err := loader.AllRegos(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Make timeout configurable via command line flag
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	regal := linter.NewLinter().WithAddedBundle(regalRules)

	regal.Lint(ctx, policies)
}