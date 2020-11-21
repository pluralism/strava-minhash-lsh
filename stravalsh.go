package stravaminhashlsh

import "math"

const largePrime = 433494437
const threshold = 0.5

type stravaLSH struct {
	bands         uint
	buckets       uint64
	stravaMinHash *stravaMinHash
}

type stravaLSHBuilder interface {
	Bands(uint) stravaLSHBuilder
	Buckets(uint64) stravaLSHBuilder
	StravaMinHash(hash *stravaMinHash) stravaLSHBuilder
	Build() *stravaLSH
}

type basicStravaLSHBuilder struct {
	bands         uint
	buckets       uint64
	stravaMinHash *stravaMinHash
}

func NewStravaLSH() stravaLSHBuilder {
	return &basicStravaLSHBuilder{}
}

func (b *basicStravaLSHBuilder) Bands(bands uint) stravaLSHBuilder {
	b.bands = bands
	return b
}

func (b *basicStravaLSHBuilder) Buckets(buckets uint64) stravaLSHBuilder {
	b.buckets = buckets
	return b
}

func (b *basicStravaLSHBuilder) StravaMinHash(hash *stravaMinHash) stravaLSHBuilder {
	b.stravaMinHash = hash
	return b
}

func (b *basicStravaLSHBuilder) Build() *stravaLSH {
	return &stravaLSH{
		bands:         b.bands,
		buckets:       b.buckets,
		stravaMinHash: b.stravaMinHash,
	}
}

func (s *stravaLSH) GetSignatureSize() uint {
	// Threshold is where the rise is considered to be the steepest.
	// Pairs with similarity above the threshold are very likely to become candidate pairs.
	//
	// An approximation to the threshold is given by (1/b)^(1/r), where "b" is the desired number of bands
	// and "r" is the number of rows by band.
	//
	// Assuming that the threshold is 0.5, the number of rows can be calculate by the following expression:
	// (1/b)^(1/r) = 0.5 <=> ln((1/b)^(1/r)) = ln(0.5) <=> (1/r) * ln(1/b) = ln(0.5) <=> r = ln(1/b) / ln(0.5)
	rows := math.Ceil(math.Log(1/float64(s.bands)) / math.Log(threshold))
	return uint(rows) * s.bands
}

func (s *stravaLSH) initialize(coefficients [][]uint64) {
	s.stravaMinHash.initialize(coefficients)
}

func (s *stravaLSH) hash(data [][]int) []uint64 {
	return s.hashSignature(s.stravaMinHash.getSignature(data))
}

func (s *stravaLSH) hashSignature(signature []uint64) []uint64 {
	rows := len(signature) / int(s.bands)
	result := make([]uint64, s.bands)
	for i := 0; i < len(result); i++ {
		result[i] = 0
	}

	// Use the same hash function for all bands, but use a separate bucket array for each band, so columns
	// with the same vector in different bands will not hash to the same bucket!
	//
	// We shall normally assume that two vectors hash to the same bucket if and only if they are identical.
	for i := 0; i < len(signature); i++ {
		bandIndex := uint(math.Floor(float64(i / rows)))
		result[bandIndex] = (result[bandIndex] + signature[i]*largePrime) % s.buckets
	}

	return result
}
