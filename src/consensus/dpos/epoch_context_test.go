package dpos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalSlot(t *testing.T) {
	testCases := []struct {
		now  int64
		slot int64
		err  error
	}{
		{
			now:  86411,
			slot: 1,
			err:  nil,
		},
		{
			now:  86421,
			slot: 2,
			err:  nil,
		},
		{
			now:  86439,
			slot: 3,
			err:  nil,
		},
		{
			now:  86499,
			slot: 9,
			err:  nil,
		},
		{
			now:  8640199,
			slot: 19,
			err:  nil,
		},
	}
	for _, cs := range testCases {
		slot, err := calSolt(PrevSlot(cs.now))
		if err != nil {
			assert.Equal(t, cs.err, err)
		}
		assert.Equal(t, cs.slot, slot)
	}
}
