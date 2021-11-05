package random_delay

import (
	"direwolf/internal/pkg/helpers"
)

var randomDelayRanges = map[int][]int{
	1: []int{312, 717},
	2: []int{717, 1123},
	3: []int{1037, 1878},
	4: []int{3086, 2473},
}

type RandomDelayRange int

func NewRandomDelayRange(randomDelayRangeName string) RandomDelayRange {
	var (
		rdRange = unknownRandomDelayRange
	)

	for num, rangeName := range randomDelayRangeNames[1:] {
		if rangeName == randomDelayRangeName {
			rdRange = RandomDelayRange(num)
		}
	}

	return rdRange
}

func (rdr RandomDelayRange) Int() int {
	return int(rdr)
}

const (
	unknownRandomDelayRange RandomDelayRange = iota
	defaultRandomDelayRange
	shortRandomDelayRange
	mediumRandomDelayRange
	longRandomDelayRange
)

var randomDelayRangeNames = []string{
	"defaultRandomDelayRange",
	"shortRandomDelayRange",
	"mediumRandomDelayRange",
	"longRandomDelayRange",
}

type RandomDelayGenerator struct {
	delayRange RandomDelayRange
}

func NewRandomDelayGenerator(rangeName string) *RandomDelayGenerator {
	//randomDelayRangeName := os.Getenv("DW_DEFAULT_TOR_CRAWLER_RANDOM_DELAY_RANGE")
	return &RandomDelayGenerator{
		delayRange: NewRandomDelayRange(rangeName),
	}
}

func (rdg *RandomDelayGenerator) GenerateRandomDelay() int {
	nums := randomDelayRanges[rdg.delayRange.Int()]

	random := nums[0] + helpers.RandomInt(nums[1])
	return random
}
