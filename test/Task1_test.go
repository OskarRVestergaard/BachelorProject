package test

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTask1Initialization(t *testing.T) {
	prm := Task1.GenerateTestingParameters()
	_, _, result := Task1.SimulateInitialization(prm)
	assert.True(t, result)
}

func TestTask1InitializationAndExecution(t *testing.T) {
	prm := Task1.GenerateTestingParameters()
	p, v, initOkay := Task1.SimulateInitialization(prm)
	assert.True(t, initOkay)
	execOkay := Task1.SimulateExecution(p, v)
	assert.True(t, execOkay)
}

func TestTask1WithExpander(t *testing.T) {
	prm := Task1.GenerateParameters(5, 16, 16)
	p, v, initOkay := Task1.SimulateInitialization(prm)
	assert.True(t, initOkay)
	execOkay := Task1.SimulateExecution(p, v)
	assert.True(t, execOkay)
}

//TODO ADD Negative tests, such as, prover not sending all information needed or sending wrong values
