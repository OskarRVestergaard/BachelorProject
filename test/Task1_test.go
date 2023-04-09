package test

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTask1Initialization(t *testing.T) {
	_, _, result := Task1.SimulateInitialization()
	assert.True(t, result)
}

func TestTask1InitializationAndExecution(t *testing.T) {
	p, v, initOkay := Task1.SimulateInitialization()
	assert.True(t, initOkay)
	execOkay := Task1.SimulateExecution(p, v)
	assert.True(t, execOkay)
}

//TODO ADD Negative tests, such as, prover not sending all information needed or sending wrong values
