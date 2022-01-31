package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConst_IsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		c   Const
		exp bool
	}{
		{Const{Val: "test"}, true},
		{Const{Val: 1}, false},
		{Const{Val: 1.2}, false},
		{Const{Val: true}, false},
		{Const{Val: false}, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.c.IsString())
	}
}
