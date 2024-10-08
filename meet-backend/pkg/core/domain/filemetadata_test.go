package domain_test

import (
	"testing"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetFileName(t *testing.T) {
	fileMetaData := domain.FileMetaData{
		Path: "dir1/dir2/file.foo",
	}

	expected := "file.foo"
	actual := fileMetaData.GetFileName()

	assert.Equal(t, expected, actual)
}

func TestGetFolders(t *testing.T) {
	fileMetaData := domain.FileMetaData{
		Path: "dir1/dir2/file.foo",
	}

	expected := []string{"dir1", "dir2"}
	actual := fileMetaData.GetFolders()

	assert.Equal(t, expected, actual)
}

func TestGetFoldersNoFolder(t *testing.T) {
	fileMetaData := domain.FileMetaData{
		Path: "file.foo",
	}

	expected := []string{}
	actual := fileMetaData.GetFolders()

	assert.Equal(t, expected, actual)
}

func TestGetFoldersPath(t *testing.T) {
	fileMetaData := domain.FileMetaData{
		Path: "dir1/dir2/file.foo",
	}

	expected := "dir1/dir2"
	actual := fileMetaData.GetFoldersPath()

	assert.Equal(t, expected, actual)
}
