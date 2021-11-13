package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func GetParentDir(t *testing.T) string {
	// Set working directory to parent folder
	wd, err := os.Getwd()
	assert.NoError(t, err)
	pd := filepath.Dir(wd)
	return pd
}
