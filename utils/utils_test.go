package utils

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transform", func() {
	It("Transforms each element", func() {
		count := 10
		slice := make([]int, count)
		expected := make([]string, count)

		for i := 0; i < count; i++ {
			slice[i] = i
			expected[i] = fmt.Sprint(i)
		}

		result := Transform(slice, func(i int) string { return fmt.Sprint(i) })

		Expect(result).To(ContainElements(expected))
	})
})

var _ = Describe("TransformValues", func() {
	It("Transforms each value", func() {
		count := 10
		m := make(map[int]int, count)
		expected := make([]string, count)

		for i := 0; i < count; i++ {
			m[i] = i
			expected[i] = fmt.Sprint(i)
		}

		result := TransformValues(m, func(i int) string { return fmt.Sprint(i) })

		Expect(result).To(ContainElements(expected))
	})
})

var _ = Describe("Where", func() {
	It("Filters slice", func() {
		count := 10
		slice := make([]int, count)

		for i := 0; i < count; i++ {
			slice[i] = i
		}

		result := Where(slice, func(i int) bool { return i < 5 })

		Expect(result).To(HaveLen(5))
		Expect(result).To(ContainElements(BeNumerically("<", 5)))
	})
})

var _ = Describe("ValuesWhere", func() {
	It("Filters values", func() {
		count := 10
		m := make(map[int]int, count)

		for i := 0; i < count; i++ {
			m[i] = i
		}

		result := ValuesWhere(m, func(i int) bool { return i < 5 })

		Expect(result).To(HaveLen(5))
		Expect(result).To(ContainElements(BeNumerically("<", 5)))
	})
})

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

var _ = Describe("OrderedValues", func() {
	It("returns the values in order", func() {
		count := 10
		m := make(map[int]int, count)
		expected := make([]int, count)

		for i := count - 1; i >= 0; i-- {
			m[i] = i
		}
		for i := 0; i < count; i++ {
			expected[i] = i
		}

		result := OrderedValues(m)

		Expect(result).To(Equal(expected))
	})
})
