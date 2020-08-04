package acceptance_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance", func() {
	It("exits successfully", func() {
		Expect(exec.Command(binaryPath).Run()).To(Succeed())
	})
})
