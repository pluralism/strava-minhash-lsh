package stravaminhashlsh

import (
	"errors"
	"fmt"
)

type StravaSignatureCalculator struct {
	MinHash *stravaMinHash
	LSH     *stravaLSH
}

type stravaSignatureCalculatorBuilder interface {
	MinHash(minHash *stravaMinHash) stravaSignatureCalculatorBuilder
	LSH(lsh *stravaLSH) stravaSignatureCalculatorBuilder
	Build() *StravaSignatureCalculator
}

type basicStravaSignatureCalculatorBuilder struct {
	minHash *stravaMinHash
	lsh     *stravaLSH
}

func (b *basicStravaSignatureCalculatorBuilder) MinHash(minHash *stravaMinHash) stravaSignatureCalculatorBuilder {
	b.minHash = minHash
	return b
}

func (b *basicStravaSignatureCalculatorBuilder) LSH(lsh *stravaLSH) stravaSignatureCalculatorBuilder {
	b.lsh = lsh
	return b
}

func (b *basicStravaSignatureCalculatorBuilder) Build() *StravaSignatureCalculator {
	return &StravaSignatureCalculator{
		MinHash: b.minHash,
		LSH:     b.lsh,
	}
}

func NewStravaSignatureCalculator() stravaSignatureCalculatorBuilder {
	return &basicStravaSignatureCalculatorBuilder{}
}

func (b *StravaSignatureCalculator) Setup(coefficients [][]uint64) error {
	signatureSize := b.LSH.GetSignatureSize()
	if int(signatureSize) != len(coefficients) {
		return errors.New(fmt.Sprintf("coefficients size %d is different from signature size %d.",
			len(coefficients), signatureSize))
	}
	b.LSH.initialize(coefficients)
	return nil
}

func (b *StravaSignatureCalculator) GetSignature(data [][]float64) []uint64 {
	activityTiles := getActivityTiles(data, b.MinHash.zoom)
	bresenhamPath := getBresenhamPathForTiles(activityTiles)
	return b.LSH.hash(bresenhamPath)
}

func getActivityTiles(data [][]float64, zoom int32) [][]int {
	activityTiles := make([][]int, 0, len(data))
	index := 0
	lastX := -1
	lastY := -1

	for _, latLngPair := range data {
		x, y := deg2Num(latLngPair[0], latLngPair[1], zoom)
		if x == lastX && y == lastY {
			continue
		}
		activityTiles = append(activityTiles, []int{x, y})
		index++
		lastX = x
		lastY = y
	}

	return activityTiles[:index]
}

func getBresenhamPathForTiles(tiles [][]int) [][]int {
	var bresenhamPath [][]int

	for i := 1; i < len(tiles); i++ {
		result := bresenham(tiles[i-1][0], tiles[i-1][1], tiles[i][0], tiles[i][1])
		if i > 1 {
			result.array = result.array[1:]
		}
		bresenhamPath = append(bresenhamPath, result.array...)
	}

	return bresenhamPath
}
