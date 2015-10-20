package assertion_test

import (
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/ginkgo"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega"

	"testing"
)

func TestAssertion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Assertion Suite")
}
