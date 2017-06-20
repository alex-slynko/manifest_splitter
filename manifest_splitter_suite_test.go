package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestManifestSplitter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ManifestSplitter Suite")
}

var pathToPackage string

var _ = BeforeSuite(func() {
	var err error
	pathToPackage, err = gexec.Build("github.com/alex-slynko/manifest_splitter")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
