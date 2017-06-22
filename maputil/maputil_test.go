package maputil_test

import (
	"github.com/alex-slynko/manifest_splitter/maputil"
	"github.com/alex-slynko/manifest_splitter/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Maputil", func() {

	var first, second map[string]interface{}
	BeforeEach(func() {
		first = map[string]interface{}{
			"property": "value",
		}
	})

	It("detect missing elements", func() {
		second = map[string]interface{}{}
		operations, err := maputil.ExtractOperations(first, second)
		Expect(err).NotTo(HaveOccurred())
		Expect(operations).To(Equal([]types.Operation{
			types.Operation{
				Path:  "/property?",
				Type:  "replace",
				Value: "value",
			},
		}))
	})
	It("detect extra elements", func() {
		first = map[string]interface{}{}
		second = map[string]interface{}{
			"property": "value",
		}
		operations, err := maputil.ExtractOperations(first, second)
		Expect(err).NotTo(HaveOccurred())
		Expect(operations).To(Equal([]types.Operation{
			types.Operation{
				Path: "/property",
				Type: "remove",
			},
		}))
	})

	It("detect different elements", func() {
		second = map[string]interface{}{
			"property": "wrongvalue",
		}
		operations, err := maputil.ExtractOperations(first, second)
		Expect(err).NotTo(HaveOccurred())
		Expect(operations).To(Equal([]types.Operation{
			types.Operation{
				Path:  "/property",
				Type:  "replace",
				Value: "value",
			},
		}))
	})

	Context("when new value is a map", func() {
		BeforeEach(func() {
			first = map[string]interface{}{
				"property": map[interface{}]interface{}{},
			}

		})
		Context("when old value is a map", func() {
			BeforeEach(func() {
				second = map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "wrongvalue",
					},
				}
			})

			It("detect extra elements", func() {
				operations, err := maputil.ExtractOperations(first, second)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(Equal([]types.Operation{
					types.Operation{
						Path: "/property/nested",
						Type: "remove",
					},
				}))
			})

			It("detects different subelements", func() {
				first = map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "value",
					},
				}
				operations, err := maputil.ExtractOperations(first, second)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(Equal([]types.Operation{
					types.Operation{
						Path:  "/property/nested",
						Type:  "replace",
						Value: "value",
					},
				}))
			})
		})

		Context("when old value is not a map", func() {
			BeforeEach(func() {
				second = map[string]interface{}{
					"property": "value",
				}
			})

			It("it returns an error", func() {
				_, err := maputil.ExtractOperations(first, second)
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Context("when old value is a map", func() {
		BeforeEach(func() {
			second = map[string]interface{}{
				"property": map[interface{}]interface{}{
					"nested": "wrongvalue",
				},
			}
		})

		Context("when new value is not a map", func() {
			It("it returns an error", func() {
				_, err := maputil.ExtractOperations(first, second)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
