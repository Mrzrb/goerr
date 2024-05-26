package aop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicTest(t *testing.T) {
	_, err := Run("../testsdata/basic")
	assert.Nil(t, err, "run basic testcase error, notNil")
}

func TestSubTest(t *testing.T) {
	_, err := Run("../testsdata/subdir/")
	assert.Nil(t, err, "run basic testcase error, notNil")
}
