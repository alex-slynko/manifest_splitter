package main_test

import (
	"os/exec"

	yaml "gopkg.in/yaml.v2"

	"github.com/alex-slynko/manifest_splitter/types"
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
		command := exec.Command(pathToPackage, "sample.yml", "sample.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	})

	It("", func() {
		command := exec.Command(pathToPackage, "manifest.yml", "small_manifest.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		output := string(session.Out.Contents())
		Expect(output).To(ContainSubstring(`- type: replace
  path: /name
  value: new-name
`))

	})

	It("outputs operation-file", func() {
		command := exec.Command(pathToPackage, "manifest.yml", "small_manifest.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		output := session.Out.Contents()
		ops := []types.Operation{}
		yaml.Unmarshal(output, &ops)
		Expect(ops).To(ContainElement(types.Operation{
			Type:  "replace",
			Path:  "/properties/a?",
			Value: "key",
		}))
		Expect(ops).To(ContainElement(types.Operation{
			Type:  "replace",
			Path:  "/properties/array/-",
			Value: "a",
		}))
	})

	It("outputs just the required operation-file", func() {
		command := exec.Command(pathToPackage, "manifest.yml", "small_manifest.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		output := session.Out.Contents()
		ops := []types.Operation{}
		yaml.Unmarshal(output, &ops)
		Expect(ops).ToNot(ContainElement(types.Operation{
			Type: "remove",
			Path: "/stemcells/0",
		}))
	})

	It("outputs nothing when manifest and operation-file are equivalent", func() {
		command := exec.Command(pathToPackage, "file.yml", "equivalent_file.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		output := session.Out.Contents()
		ops := []types.Operation{}
		yaml.Unmarshal(output, &ops)
		Expect(ops).To(HaveLen(0))
	})
})
