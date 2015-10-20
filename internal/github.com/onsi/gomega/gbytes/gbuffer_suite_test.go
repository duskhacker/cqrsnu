package gbytes_test

import (
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/ginkgo"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega"

	"testing"
)

func TestGbytes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gbytes Suite")
}
