package test

import (
	"github.com/stretchr/testify/assert"
	"go/format"
	"strings"
	"testing"
)

func PrepExpCode(t *testing.T, code string) string {
	code = strings.TrimPrefix(code, "\n")
	p, err := format.Source([]byte(code))
	assert.NoError(t, err)
	return string(p)
}
