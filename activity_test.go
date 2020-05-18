package sample

import (
	"io/ioutil"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {
	td, err := ioutil.ReadFile("token.txt")
	assert.Nil(t, err)
	settings := &Settings{URI: "api2.rocos.io:443", Token: string(td)}
	ictx := test.NewActivityInitContext(settings, nil)

	act, err := New(ictx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	in := &Input{ProjectID: "takahe", Callsign: "t-1"}
	err = tc.SetInputObject(in)
	assert.Nil(t, err)

	t.Logf("input value: %v", tc.GetInput("projectID"))

	assert.Nil(t, err)
	t.Logf("input object: %+v", in)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "final", output.AnOutput)
}
