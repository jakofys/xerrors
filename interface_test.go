package xerrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jakofys/xerrors"
	"github.com/stretchr/testify/assert"
)

type (
	errUnimplement struct{}
	errImplement   struct {
		msg string
	}
)

func (errUnimplement) Error() string        { return "error 2" }
func (errImplement) Error() string          { return "error 1" }
func (e errImplement) IsImplements() string { return e.msg }

type Implementer interface {
	IsImplements() string
}

func TestAsInterface(t *testing.T) {
	const implementMsg = "implements"
	cases := []struct {
		name, msg      string
		expectedNil    bool
		expectedString string
		input          error
	}{
		{
			name:        "no implementation",
			msg:         "should return nil interface",
			input:       errUnimplement{},
			expectedNil: true,
		},
		{
			name:           "direct implementation",
			msg:            "should return interface implementation",
			input:          errImplement{msg: implementMsg},
			expectedString: implementMsg,
		},
		{
			name:           "wrapped error implementation",
			msg:            "should return interface implementation within a wrapped error",
			input:          fmt.Errorf("simple: %w", errImplement{msg: implementMsg}),
			expectedString: implementMsg,
		},
		{
			name:        "joined errors implementation",
			msg:         "should return nil implementation due to unwrap fail",
			input:       errors.Join(errors.New("error"), errImplement{msg: implementMsg}),
			expectedNil: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := xerrors.AsInterface[Implementer](c.input)
			if c.expectedNil {
				assert.Nil(t, result, c.msg)
			} else {
				assert.Equal(t, result.IsImplements(), c.expectedString)
			}
		})
	}
}

func TestAsInterfacePanicOnNoneInterface(t *testing.T) {
	assert.Panics(t, func() {
		e := xerrors.AsInterface[errUnimplement](errUnimplement{})
		assert.Nil(t, e)
	})
	assert.Panics(t, func() {
		e := xerrors.AsInterface[*errUnimplement](errUnimplement{})
		assert.Nil(t, e)
	})
}
