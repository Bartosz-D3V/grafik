package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSnakeCaseToPascalCase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		exp   string
	}{
		{"num_of_missions", "numOfMissions"},
		{"total_num_of_missions_in_the_us_zone_1_c", "totalNumOfMissionsInTheUsZone1C"},
		{"i_phone", "iPhone"},
		{"total_count", "totalCount"},
		{"asci", "asci"},
		{"x", "x"},
		{"A_1", "A1"},
		{"1_X", "1X"},
		{"1_abc", "1Abc"},
		{"landing_zone_", "landingZone"},
		{"", ""},
	}
	for _, test := range tests {
		res := SnakeCaseToCamelCase(test.input)
		assert.Equal(t, test.exp, res)
	}
}
