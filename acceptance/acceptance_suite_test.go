package acceptance_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

const binaryPath = "../build/crud"

func tempDir(dir, pattern string) string {
	name, err := ioutil.TempDir(dir, pattern)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return name
}
