package stravaminhashlsh

import (
	"hash/fnv"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

const hashPrime = 2038074743

type stravaMinHash struct {
	zoom         int32
	shingleSize  int
	coefficients [][]uint64
}

type stravaMinHashBuilder interface {
	Zoom(int32) stravaMinHashBuilder
	ShingleSize(int) stravaMinHashBuilder
	Build() *stravaMinHash
}

type basicStravaMinHashBuilder struct {
	zoom        int32
	shingleSize int
}

func NewStravaMinHash() stravaMinHashBuilder {
	return &basicStravaMinHashBuilder{}
}

func (b *basicStravaMinHashBuilder) Zoom(zoom int32) stravaMinHashBuilder {
	b.zoom = zoom
	return b
}

func (b *basicStravaMinHashBuilder) ShingleSize(size int) stravaMinHashBuilder {
	b.shingleSize = size
	return b
}

func (b *basicStravaMinHashBuilder) Build() *stravaMinHash {
	return &stravaMinHash{
		zoom:        b.zoom,
		shingleSize: b.shingleSize,
	}
}

func (a *stravaMinHash) initialize(coefficients [][]uint64) {
	a.coefficients = coefficients
}

func (a *stravaMinHash) getSignature(data [][]int) []uint64 {
	shingles := a.buildShingles(data)
	uniqueShingles := a.getUniqueShingles(shingles)
	hashedShingles := a.getHashedShingles(uniqueShingles)

	signature := make([]uint64, len(a.coefficients))

	for i := 0; i < len(a.coefficients); i++ {
		min := uint64(math.MaxUint64)
		coefficient := a.coefficients[i]
		for _, shingle := range hashedShingles {
			result := (coefficient[0]*uint64(shingle) + coefficient[1]) % hashPrime
			if result < min {
				min = result
			}
		}
		signature[i] = min
	}

	return signature
}

func (a *stravaMinHash) GetRandomCoefficients(count uint) [][]uint64 {
	var randomCoefficients [][]uint64
	for i := 0; i < int(count); i++ {
		a := math.Floor(float64(1 + rand.Int63n(int64(hashPrime-1))))
		b := math.Floor(float64(1 + rand.Int63n(int64(hashPrime-1))))
		randomCoefficients = append(randomCoefficients, []uint64{uint64(a), uint64(b)})
	}
	return randomCoefficients
}

func (a *stravaMinHash) getHashedShingles(shingles [][][]int) []uint32 {
	var hashedShingles []uint32
	// 32-bit FNV-1a hash
	hash32 := fnv.New32a()
	for _, shingle := range shingles {
		key := a.getShingleKey(shingle)
		_, _ = hash32.Write([]byte(key))
		hashedShingles = append(hashedShingles, hash32.Sum32())
		hash32.Reset()
	}
	return hashedShingles
}

func (a *stravaMinHash) getShingleKey(shingle [][]int) string {
	var key strings.Builder

	for _, shinglePart := range shingle {
		key.WriteString(strconv.Itoa(shinglePart[0]))
		key.WriteString(strconv.Itoa(shinglePart[1]))
	}

	return key.String()
}

func (a *stravaMinHash) getUniqueShingles(shingles [][][]int) [][][]int {
	var uniqueShingles [][][]int
	set := newSet()
	for _, shingle := range shingles {
		key := a.getShingleKey(shingle)
		if !set.Contains(key) {
			set.Add(key)
			uniqueShingles = append(uniqueShingles, shingle)
		}
	}
	return uniqueShingles
}

func (a *stravaMinHash) buildShingles(path [][]int) [][][]int {
	var shingles [][][]int
	for i := int(math.Min(float64(len(path)), float64(a.shingleSize))); i <= len(path); i++ {
		shingles = append(shingles, path[i-a.shingleSize:i])
	}
	return shingles
}
