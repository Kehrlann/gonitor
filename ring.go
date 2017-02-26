package main

import "container/ring"

// RingToIntSlice takes a ring.Ring, and dumps an slice of ints reprensenting the ring.
// Non-int values and nils will be zero-ed.
func RingToIntSlice(r *ring.Ring) []int {
	ret := make([]int, r.Len())
	j := 0
	for i := 0; i < r.Len(); i++ {
		intVal, ok := r.Value.(int)
		if ok {
			ret[j] = intVal
			j++
		}
		r = r.Next()
	}

	return ret
}