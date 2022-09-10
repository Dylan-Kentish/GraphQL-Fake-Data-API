package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TryLimitIfPresent", func() {
	var size int
	var hasLimit bool
	var limit int

	var args map[string]interface{}
	var slice []int

	JustBeforeEach(func() {
		args = make(map[string]interface{})
		if hasLimit {
			args["limit"] = limit
		}

		slice = make([]int, size)
		for i := 0; i < size; i++ {
			slice[i] = i
		}
	})

	Context("limit is present", func() {
		BeforeEach(func() {
			hasLimit = true
		})

		When("limit is less than len", func() {
			BeforeEach(func() {
				limit = 5
				size = 10
			})

			It("Returns limited items", func() {
				Expect(TryLimitIfPresent(slice, args)).To(Equal(slice[:limit]))
			})
		})

		When("limit is greater than len", func() {
			BeforeEach(func() {
				limit = 10
				size = 5
			})

			It("Returns all items", func() {
				Expect(TryLimitIfPresent(slice, args)).To(Equal(slice))
			})
		})
	})

	Context("limit is not present", func() {
		BeforeEach(func() {
			hasLimit = false
		})

		It("Returns all items", func() {
			Expect(TryLimitIfPresent(slice, args)).To(Equal(slice))
		})
	})
})
