package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("delete", func() {
		It("deletes an element", func() {
			s := []int{1, 2, 3, 4}
			t := del(s, 2)

			Expect(t).To(ConsistOf([]int{1, 3, 4}))

			t = del(s, 1)
			Expect(t).To(ConsistOf([]int{2, 3, 4}))

			t = del(s, 4)
			Expect(t).To(ConsistOf([]int{1, 2, 3}))
		})
	})

})
