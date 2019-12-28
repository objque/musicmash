package cron

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutils/suite"
)

type testCronSuite struct {
	suite.Suite
}

func (t *testCronSuite) SetupTest() {
	t.Suite.SetupTest()
}

func (t *testCronSuite) TearDownTest() {
	t.Suite.TearDownTest()
}

func TestCronSuite(t *testing.T) {
	suite.Run(t, new(testCronSuite))
}
