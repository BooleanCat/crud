package acceptance_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance", func() {
	var (
		cmd     *exec.Cmd
		workDir string
	)

	BeforeEach(func() {
		workDir = tempDir("", "")

		cmd = exec.Command(binaryPath, "--store", workDir)
		cmd.Stdout = GinkgoWriter
		cmd.Stderr = GinkgoWriter
	})

	JustBeforeEach(func() {
		Expect(cmd.Start()).To(Succeed())
		Eventually(ping("http://127.0.0.1:9092")).Should(Succeed())
	})

	AfterEach(func() {
		Expect(cmd.Process.Signal(os.Interrupt)).To(Succeed())
		_, err := cmd.Process.Wait()
		Expect(err).NotTo(HaveOccurred())
		Expect(os.RemoveAll(workDir)).To(Succeed())
	})

	Describe("/Create/filename", func() {
		It("stores a file", func() {
			request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:9092/Create/bar", bytes.NewBufferString("foobar"))
			Expect(err).NotTo(HaveOccurred())

			response, err := http.DefaultClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()

			By("responding with 200 OK", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			By("by storing the contents of the request body", func() {
				body, err := ioutil.ReadFile(filepath.Join(workDir, "bar"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("foobar"))
			})
		})
	})

	Describe("/Read/filename", func() {
		It("reads a file", func() {
			Expect(create()).To(BeNil())

			readResp, err := http.Get("http://127.0.0.1:9092/Read/bar")
			Expect(err).NotTo(HaveOccurred())
			defer readResp.Body.Close()

			By("responding with 200 OK", func() {
				Expect(readResp.StatusCode).To(Equal(http.StatusOK))
			})

			By("responding with the content of that file", func() {
				body, err := ioutil.ReadAll(readResp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("Content: foobar"))
			})
		})
	})
})

func ping(url string) func() error {
	return func() error {
		response, err := http.Get(url + "/ping")
		if err != nil {
			return err
		}
		_ = response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected response code %d", response.StatusCode)
		}

		return nil
	}
}

func create() func() error {
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:9092/Create/bar", bytes.NewBufferString("foobar"))
	Expect(err).NotTo(HaveOccurred())

	_, err = http.DefaultClient.Do(request)
	Expect(err).NotTo(HaveOccurred())

	return nil
}
