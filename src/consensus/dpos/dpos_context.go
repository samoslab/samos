package dpos

import (
	"errors"

	"github.com/samoslab/samos/src/cipher"
)

type DposContext struct {
	candidate []cipher.PubKey
}

func NewDposContext() *DposContext {
	//address := []string{"2fxav8p7QFkKk8TBwmE6wvu8S8VVEyvpX8C", "CB2tqSePaPBrMiBh2513njfUtev8GfMjEX", "jFAUc1AUeAgVjc4Br5mv3baaQkuiKZ7maw"}
	//candidate := []cipher.PubKey{}
	//for _, v := range address {
	//	candidate = append(candidate, cipher.MustDecodeBase58PubKey(v))
	//}
	return &DposContext{}
}

func (dc *DposContext) GetValidators() ([]cipher.PubKey, error) {
	if len(dc.candidate) == 0 {
		return dc.candidate, errors.New("zero validator")
	}
	return dc.candidate, nil
}

func (dc *DposContext) SetValidators(validators []cipher.PubKey) error {
	dc.candidate = validators
	return nil
}
