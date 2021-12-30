package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSentenceCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		exp   string
	}{
		{"Person", "person"},
		{"Total_count", "total_count"},
		{"1Number", "1Number"},
		{"xRay", "xRay"},
		{"car", "car"},
		{"", ""},
	}
	for _, test := range tests {
		res := SentenceCase(test.input)
		assert.Equal(t, test.exp, res)
	}
}

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
