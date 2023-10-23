package tfuncs

import (
	"errors"
	"fmt"
	"reflect"
)

// after slices an array to only the items after the Nth item.
func after(index any, seq any) (any, error) {
	if index == nil || seq == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	indexv, err := tointE(index)
	if err != nil {
		return nil, err
	}

	if indexv < 0 {
		s := fmt.Sprintf("%d", indexv)
		return nil, errors.New("sequence bounds out of range [" + s + ":]")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirectValue(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}

	if indexv >= seqv.Len() {
		return seqv.Slice(0, 0).Interface(), nil
	}

	return seqv.Slice(indexv, seqv.Len()).Interface(), nil
}
