package datalayer

import (
	"os"

	commonV1 "github.com/binuud/ai-green-field/gen/go/v1/common"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListFilesNonRecursive lists ONLY files in the specified directory (no subdirs)
func ListFilesNonRecursive(dirPath string) ([]*commonV1.FileMetadata, error) {
	// var files []*commonV1.FileMetadata

	files := make([]*commonV1.FileMetadata, 0)

	// Open directory directly (no recursion)
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	logrus.Infof("Listing files from folder %s", dirPath)
	// Read directory entries
	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// Skip directories entirely
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip hidden files
		if isHiddenFile(name) {
			continue
		}

		//fullPath := filepath.Join(dirPath, name)
		files = append(files, &commonV1.FileMetadata{
			Name:        name,
			Path:        "",
			Size:        entry.Size(),
			ModifiedAt:  timestamppb.New(entry.ModTime()),
			IsDirectory: false, // We already filtered directories
		})
	}

	return files, nil
}

func isHiddenFile(name string) bool {
	return len(name) > 0 && name[0] == '.'
}
