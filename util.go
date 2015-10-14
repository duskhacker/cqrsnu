package main

import "log"

func indexOf(s []int, el int) int {
	for i, e := range s {
		if e == el {
			return i
		}
	}
	return -1
}

func del(s []int, el int) []int {
	idx := indexOf(s, el)
	if idx < 0 {
		return s
	}
	a := make([]int, len(s))
	n := copy(a, s)
	if n <= 0 {
		log.Fatalf("error copying data for del")
	}
	a = append(a[:idx], a[idx+1:]...)
	return a
}
