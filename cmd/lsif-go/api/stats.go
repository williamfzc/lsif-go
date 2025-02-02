package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/sourcegraph/lsif-go/internal/indexer"
	"github.com/sourcegraph/lsif-go/internal/util"
)

func displayStats(indexerStats indexer.IndexerStats, packageDataCacheStats indexer.PackageDataCacheStats, start time.Time) {
	stats := []struct {
		name  string
		value string
	}{
		{"Wall time elapsed", fmt.Sprintf("%s", util.HumanElapsed(start))},
		{"Packages indexed", fmt.Sprintf("%d", indexerStats.NumPkgs)},
		{"Files indexed", fmt.Sprintf("%d", indexerStats.NumFiles)},
		{"Definitions indexed", fmt.Sprintf("%d", indexerStats.NumDefs)},
		{"Elements emitted", fmt.Sprintf("%d", indexerStats.NumElements)},
		{"Packages traversed", fmt.Sprintf("%d", packageDataCacheStats.NumPks)},
	}

	n := 0
	for _, stat := range stats {
		if n < len(stat.name) {
			n = len(stat.name)
		}
	}

	fmt.Printf("\nStats:\n")

	for _, stat := range stats {
		fmt.Printf("\t%s: %s%s\n", stat.name, strings.Repeat(" ", n-len(stat.name)), stat.value)
	}
}
