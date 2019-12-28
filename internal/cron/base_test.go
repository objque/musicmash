package cron

import (
	"testing"

	"github.com/musicmash/musicmash/internal/testutils/suite"
)

type testCronSuite struct {
	suite.Suite
}

func TestCronSuite(t *testing.T) {
	suite.Run(t, new(testCronSuite))
}
