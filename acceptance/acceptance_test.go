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

			resp, err := http.Get("http://127.0.0.1:9092/Read/bar")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			By("responding with 200 OK", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			By("responding with the content of that file", func() {
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("Content: foobar"))
			})
		})
	})

	Describe("/Update/filename", func() {
		It("updates a file", func() {
			Expect(create()).To(BeNil())

			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9092/Update/bar", bytes.NewBufferString("barfoo"))
			Expect(err).NotTo(HaveOccurred())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			By("responding with 200 OK", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			By("by storing the contents of the request body", func() {
				body, err := ioutil.ReadFile(filepath.Join(workDir, "bar"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("barfoo"))
			})
		})

		It("cannot update a non-existant file", func() {
			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9092/Update/foo", bytes.NewBufferString("barfoo"))
			Expect(err).NotTo(HaveOccurred())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			By("by responding with a 404 error code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("/Delete/filename", func() {
		It("deletes an existing file", func() {
			Expect(create()).To(BeNil())

			req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:9092/Delete/bar", nil)
			Expect(err).NotTo(HaveOccurred())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			By("responding with 200 OK", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			By("by ensuring that the existing named file is deleted from the store", func() {
				body, err := ioutil.ReadFile(filepath.Join(workDir, "bar"))
				Expect(err).To(HaveOccurred())
				Expect(string(body)).NotTo(Equal("foobar"))
			})
		})

		It("cannot delete a non-existant file", func() {
			req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:9092/Delete/foobar", nil)
			Expect(err).NotTo(HaveOccurred())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			By("by responding with a 404 error code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
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

func create() error {
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:9092/Create/bar", bytes.NewBufferString("foobar"))
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	return nil
}
