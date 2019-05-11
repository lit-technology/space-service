package s3

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestS3(t *testing.T) {
	suite.Run(t, new(S3Suite))
}

type S3Suite struct {
	suite.Suite
}

func (s *S3Suite) SetupSuite() {
	// TODO:
}
