package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat_Shift(t *testing.T) {
	f0 := new(Float).SetUint64(22)
	f1 := f0.Shift(2)
	assert.Equal(t, f1.String(), "0.22")
}
