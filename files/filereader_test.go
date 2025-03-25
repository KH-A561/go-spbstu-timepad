package files_test

import (
	"testing"
	"universityTimepad/files"
)

func TestFileReader_ReadFaculties(t *testing.T) {
	t.Run("should fetch articles from all sources", func(t *testing.T) {
		var (
			reader = files.FileReader{}
		)

		reader.ReadFacultiesWithGroups()
	})
}
