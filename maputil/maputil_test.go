package maputil_test

import (
	"github.com/alex-slynko/manifest_splitter/maputil"
	"github.com/alex-slynko/manifest_splitter/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Maputil", func() {
	It("detect missing elements", func() {
		first := map[string]interface{}{
			"property": "value",
		}
		second := map[string]interface{}{}
		something, err := maputil.SomeFunction(first, second)
		Expect(err).NotTo(HaveOccurred())
		Expect(something).To(Equal([]types.Operation{
			types.Operation{
				Path:  "/property?",
				Type:  "replace",
				Value: "value",
			},
		}))
	})
	It("detect extra elements", func() {
		first := map[string]interface{}{}
		second := map[string]interface{}{
			"property": "value",
		}
		something, err := maputil.SomeFunction(first, second)
		Expect(err).NotTo(HaveOccurred())
		Expect(something).To(Equal([]types.Operation{
			types.Operation{
				Path: "/property",
				Type: "remove",
			},
		}))
	})

	It("detect different elements", func() {
		first := map[string]interface{}{
			"property": "value",
		}
		second := map[string]interface{}{
			"property": "wrongvalue",
		}
		something, err := maputil.SomeFunction(first, second)
		Expect(err).NotTo(HaveOccurred())
		Expect(something).To(Equal([]types.Operation{
			types.Operation{
				Path:  "/property",
				Type:  "replace",
				Value: "value",
			},
		}))
	})

	Context("when new value is a map", func() {
		Context("when old value is a map", func() {
			It("detect extra elements", func() {
				first := map[string]interface{}{
					"property": map[interface{}]interface{}{},
				}
				second := map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "wrongvalue",
					},
				}
				something, err := maputil.SomeFunction(first, second)
				Expect(err).NotTo(HaveOccurred())
				Expect(something).To(Equal([]types.Operation{
					types.Operation{
						Path: "/property/nested",
						Type: "remove",
					},
				}))
			})

			It("detects different subelements", func() {
				first := map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "value",
					},
				}
				second := map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "wrongvalue",
					},
				}
				something, err := maputil.SomeFunction(first, second)
				Expect(err).NotTo(HaveOccurred())
				Expect(something).To(Equal([]types.Operation{
					types.Operation{
						Path:  "/property/nested",
						Type:  "replace",
						Value: "value",
					},
				}))
			})
		})

		Context("when old value is not a map", func() {
			It("it returns an error", func() {
				first := map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "value",
					},
				}
				second := map[string]interface{}{
					"property": "value",
				}
				_, err := maputil.SomeFunction(first, second)
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Context("when old value is a map", func() {
		Context("when new value is not a map", func() {
			It("it returns an error", func() {
				first := map[string]interface{}{
					"property": "value",
				}
				second := map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "wrongvalue",
					},
				}
				_, err := maputil.SomeFunction(first, second)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
