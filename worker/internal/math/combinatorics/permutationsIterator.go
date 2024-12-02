package combinatorics

// PermutationsIterator - permutation iterator
type PermutationsIterator struct {
	// Alphabet - alphabet
	Alphabet string
	// AlphabetLen - length of alphabet
	AlphabetLen int
	// Last - last element of alphabet
	Last byte
	// PermutationLength - length of permutation
	PermutationLength int
	// CurrentPermutation - CurrentPermutation permutation
	CurrentPermutation []byte
	// LastPermutation - last permutation
	LastPermutation []byte
}

// GetNextLetter - get next letter
func (pi *PermutationsIterator) GetNextLetter(elem byte) byte {
	// get index of CurrentPermutation letter in alphabet
	index := GetIndex(elem, pi.Alphabet)
	// return next letter
	return pi.Alphabet[(index+1)%pi.AlphabetLen]
}

// NextPermutation - get next permutation
// There is an error here, related to the fact that a new character is not added if this is the last element of this length
func (pi *PermutationsIterator) NextPermutation() string {
	for i := pi.PermutationLength - 1; i >= 0; i-- {
		if pi.CurrentPermutation[i] != pi.Last {
			pi.CurrentPermutation[i] = pi.GetNextLetter(pi.CurrentPermutation[i])
			break
		}
		pi.CurrentPermutation[i] = pi.Alphabet[0]
	}
	return string(pi.CurrentPermutation)
}

// GetIndex - get index
func GetIndex(elem byte, a string) int {
	l, r := 0, len(a)
	m := 0
	for l < r {
		m = l + (r-l)/2
		if a[m] < elem {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}
