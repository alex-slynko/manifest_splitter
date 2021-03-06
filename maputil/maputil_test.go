package maputil_test

import (
	"github.com/alex-slynko/manifest_splitter/maputil"
	"github.com/alex-slynko/manifest_splitter/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Maputil", func() {

	var newValue, oldValue map[string]interface{}
	BeforeEach(func() {
		newValue = map[string]interface{}{
			"property": "value",
		}
	})

	It("detect missing elements", func() {
		oldValue = map[string]interface{}{}
		operations, err := maputil.ExtractOperations(newValue, oldValue)
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
		newValue = map[string]interface{}{}
		oldValue = map[string]interface{}{
			"property": "value",
		}
		operations, err := maputil.ExtractOperations(newValue, oldValue)
		Expect(err).NotTo(HaveOccurred())
		Expect(operations).To(Equal([]types.Operation{
			types.Operation{
				Path: "/property",
				Type: "remove",
			},
		}))
	})

	It("detect different elements", func() {
		oldValue = map[string]interface{}{
			"property": "wrongvalue",
		}
		operations, err := maputil.ExtractOperations(newValue, oldValue)
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
			newValue = map[string]interface{}{
				"property": map[interface{}]interface{}{},
			}

		})
		Context("when old value is a map", func() {
			BeforeEach(func() {
				oldValue = map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "wrongvalue",
					},
				}
			})

			It("detect extra elements", func() {
				operations, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(Equal([]types.Operation{
					types.Operation{
						Path: "/property/nested",
						Type: "remove",
					},
				}))
			})

			It("detects different subelements", func() {
				newValue = map[string]interface{}{
					"property": map[interface{}]interface{}{
						"nested": "value",
					},
				}
				operations, err := maputil.ExtractOperations(newValue, oldValue)
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

		Context("when old value is not a map of maps", func() {
			BeforeEach(func() {
				oldValue = map[string]interface{}{
					"property": "value",
				}
			})

			It("it returns an error", func() {
				_, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Context("when old value is a map", func() {
		BeforeEach(func() {
			oldValue = map[string]interface{}{
				"property": map[interface{}]interface{}{
					"nested": "wrongvalue",
				},
			}
		})

		Context("when new value is not a map", func() {
			It("it returns an error", func() {
				_, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when new value is a slice", func() {
		BeforeEach(func() {
			newValue = map[string]interface{}{
				"property": []interface{}{"value"},
			}

		})
		Context("when old value is a slice", func() {
			BeforeEach(func() {
				oldValue = map[string]interface{}{
					"property": []interface{}{
						"originalvalue",
					},
				}
			})

			It("detect extra elements", func() {
				operations, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(ContainElement(
					types.Operation{
						Path: "/property/0",
						Type: "remove",
					}))
			})

			It("detects missing subelements", func() {
				operations, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(ContainElement(types.Operation{
					Path:  "/property/-",
					Type:  "replace",
					Value: "value",
				}))
			})
		})

		Context("when old value is not a slice", func() {
			BeforeEach(func() {
				oldValue = map[string]interface{}{
					"property": "value",
				}
			})

			It("it returns an error", func() {
				_, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Context("when old value is a slice", func() {
		BeforeEach(func() {
			oldValue = map[string]interface{}{
				"property": []interface{}{
					"wrongvalue",
				},
			}
		})

		Context("when new value is not a slice", func() {
			It("it returns an error", func() {
				_, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when values are slices of maps", func() {
		var (
			nestedOldValue map[interface{}]interface{}
			nestedNewValue map[interface{}]interface{}
		)
		BeforeEach(func() {
			nestedOldValue = map[interface{}]interface{}{}
			nestedNewValue = map[interface{}]interface{}{}
		})
		JustBeforeEach(func() {

			newValue = map[string]interface{}{
				"property": []interface{}{
					nestedNewValue,
				},
			}
			oldValue = map[string]interface{}{
				"property": []interface{}{
					nestedOldValue,
				},
			}
		})

		It("does not raise error", func() {
			_, err := maputil.ExtractOperations(newValue, oldValue)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when old value is slice of simple values", func() {
			It("returns error", func() {
				oldValue = map[string]interface{}{
					"property": []interface{}{1, 2, 3},
				}
				_, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when new value is slice of simple values", func() {
			It("returns error", func() {
				newValue = map[string]interface{}{
					"property": []interface{}{1, 2, 3},
				}
				_, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when new map contains extra elements", func() {
			It("adds operation for new map", func() {
				oldValue = map[string]interface{}{
					"property": []interface{}{},
				}
				operations, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(HaveLen(1))
				Expect(operations[0]).To(Equal(types.Operation{
					Path:  "/property/-",
					Type:  "replace",
					Value: nestedNewValue,
				}))

			})
		})

		Context("when old map contains extra elements", func() {
			It("adds operation to remove old map", func() {
				newValue = map[string]interface{}{
					"property": []interface{}{},
				}
				operations, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(HaveLen(1))
				Expect(operations[0]).To(Equal(types.Operation{
					Path: "/property/0",
					Type: "remove",
				}))

			})
		})

		Context("when map contains name", func() {
			BeforeEach(func() {
				nestedNewValue["name"] = "Test"
				nestedNewValue["value"] = "New"
				nestedOldValue["name"] = "Test"
				nestedOldValue["value"] = "Old"
			})

			It("compares map that has the same name", func() {
				operations, err := maputil.ExtractOperations(newValue, oldValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(operations).To(HaveLen(1))
				Expect(operations[0]).To(Equal(types.Operation{
					Path:  "/property/name=Test/value",
					Type:  "replace",
					Value: "New",
				}))

			})
		})
	})
})
