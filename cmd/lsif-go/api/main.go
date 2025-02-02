package api

import (
	"fmt"
	"log"
	"os"

	"github.com/sourcegraph/lsif-go/internal/git"
	"github.com/sourcegraph/lsif-go/internal/gomod"
	"github.com/sourcegraph/lsif-go/internal/indexer"
	"github.com/sourcegraph/lsif-go/internal/output"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(os.Stdout)
}

func MainArgs(args []string) (err error) {
	if err := parseArgs(args); err != nil {
		return err
	}

	if !git.Check(moduleRoot) {
		return fmt.Errorf("module root is not a git repository")
	}

	defer func() {
		if err != nil {
			// Add a new line to all errors except for ones that
			// come from parsing invalid command line arguments
			// and basic environment sanity checks.
			//
			// We will print progress unconditionally after this
			// point and we want the error text to be clearly
			// visible.
			fmt.Fprintf(os.Stderr, "\n")
		}
	}()

	outputOptions := output.Options{
		Verbosity:      getVerbosity(),
		ShowAnimations: animation,
	}

	moduleName, isStdLib, err := gomod.ModuleName(moduleRoot, repositoryRemote, outputOptions)
	if err != nil {
		return fmt.Errorf("failed to infer module name: %v", err)
	}

	dependencies, err := gomod.ListDependencies(moduleRoot, moduleName, moduleVersion, outputOptions)
	if err != nil {
		return fmt.Errorf("failed to list dependencies: %v", err)
	}

	var projectDependencies []string
	if !isStdLib {
		projectDependencies, err = gomod.ListProjectDependencies(moduleRoot)
		if err != nil {
			return fmt.Errorf("failed to list project dependencies: %v", err)
		}
	}

	generationOptions := indexer.NewGenerationOptions()
	if enableApiDocs {
		return fmt.Errorf("API Docs are no longer supported. To fix this problem, remove the -enable-api-docs flag.")
	}
	generationOptions.EnableImplementations = enableImplementations
	generationOptions.DepBatchSize = depBatchSize

	if err := writeIndex(
		repositoryRoot,
		repositoryRemote,
		projectRoot,
		moduleName,
		moduleVersion,
		dependencies,
		projectDependencies,
		outFile,
		outputOptions,
		generationOptions,
	); err != nil {
		return fmt.Errorf("failed to index: %v", err)
	}

	return nil
}
