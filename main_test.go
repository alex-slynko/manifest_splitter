package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	It("fails when manifests are not provided", func() {
		command := exec.Command(pathToPackage)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
	})

	It("succeed when two manifests are provided", func() {
		command := exec.Command(pathToPackage, "manifest.yml", "small_manifest.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	})

	It("outputs operation-file", func() {
		command := exec.Command(pathToPackage, "manifest.yml", "small_manifest.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit())
		output := string(session.Out.Contents())
		Expect(output).To(Equal(`- type: replace
  path: /properties/a
  value: key
		`))
	})
})
