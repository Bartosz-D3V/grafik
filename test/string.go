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
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return string(p)
}
