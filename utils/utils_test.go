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
