package maputil_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMaputil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Maputil Suite")
}
