package acceptance_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance", func() {
	var cmd *exec.Cmd

	BeforeEach(func() {
		cmd = exec.Command(binaryPath)
		Expect(cmd.Start()).To(Succeed())
	})

	AfterEach(func() {
		Expect(cmd.Process.Signal(os.Interrupt)).To(Succeed())
		_, _ = cmd.Process.Wait()
	})

	It("responds to a health GET request", func() {
		ping := func() error {
			response, err := http.Get("http://127.0.0.1:9092/ping")
			if err != nil {
				return err
			}
			_ = response.Body.Close()

			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected response code %d", response.StatusCode)
			}

			return nil
			}

		Eventually(ping).Should(Succeed())
	})
})
