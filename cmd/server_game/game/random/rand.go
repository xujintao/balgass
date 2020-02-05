package random

import "math/rand"

type pair struct {
	value  interface{}
	weight int
}

type RandManager struct {
	values    []pair
	weightSum int
}

func (o *RandManager) Put(value interface{}, weight int) {
	o.values = append(o.values, pair{value, weight})
	o.weightSum += weight
}

func (o *RandManager) Get() interface{} {
	key := rand.Int() % len(o.values)
	return o.values[key].value
}

func (o *RandManager) GetWithWeight() interface{} {
	// randWeight := rand.Int() % o.weightSum
	randWeight := rand.Intn(o.weightSum)
	curWeight := 0
	for _, v := range o.values {
		curWeight += v.weight
		if randWeight < curWeight {
			return v.value
		}
	}
	return 0
}
