package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAggregateError_NoErrors(t *testing.T) {
	t.Parallel()

	err := AggregateError(nil)
	require.NoError(t, err)

	err = AggregateError([]error{})
	require.NoError(t, err)
}

func TestAggregateError_IgnoresNilAndEmptyErrors(t *testing.T) {
	t.Parallel()

	errs := []error{
		nil,
		errors.New(""),
		errors.New("first"),
		nil,
		errors.New("second"),
	}

	err := AggregateError(errs)
	require.Error(t, err)
	require.Equal(t, "first; second", err.Error())
}
