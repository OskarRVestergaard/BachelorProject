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
