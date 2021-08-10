package journal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/filecoin-project/lotus/node/repo"
	"github.com/stretchr/testify/require"
)

func TestRollingRemovesOldFiles(t *testing.T) {
	r := repo.NewMemory(nil)
	req := require.New(t)
	lr, err := r.Lock(repo.FullNode)
	req.NoError(err)
	j := newtestfsjournal(t, lr, 0, 3)
	dir := filepath.Join(lr.Path(), "journal")
	for i := 0; i <= j.keep; i++ {
		files, _ := os.ReadDir(dir)
		req.Lenf(files, i, "add one file for every roll before max keep")
		j.rollJournalFile()
		// there is a lag between when the file is Create()'d and when it appears.
		// it's actually a pretty long delay.
		time.Sleep(time.Second)
	}
	// on the last iteration, one of the files should have been pruned,
	// so we should still have only the maximum kept files.
	files, _ := os.ReadDir(dir)
	req.Lenf(files, j.keep, "files are not being pruned from the journal directory")
}
