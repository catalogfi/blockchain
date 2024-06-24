package evm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEvm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Evm Suite")
}
