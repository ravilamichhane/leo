package kana

type PrefixTree struct {
	children map[string]*PrefixTree
	letter   string
	values   []string
}

func newPrefixTree() *PrefixTree {
	return &PrefixTree{map[string]*PrefixTree{}, "", []string{}}
}

func (t *PrefixTree) insert(letters, value string) {
	lettersRune := []rune(letters)

	for l, letter := range lettersRune {

		letterStr := string(letter)

		if t.children[letterStr] != nil {
			t = t.children[letterStr]
		} else {
			t.children[letterStr] = &PrefixTree{map[string]*PrefixTree{}, "", []string{}}
			t = t.children[letterStr]
		}

		if l == len(lettersRune)-1 {
			t.values = append(t.values, value)
			break
		}
	}
}

func (t *PrefixTree) convert(origin string) (result string) {
	root := t
	originRune := []rune(origin)
	result = ""

	for l := 0; l < len(originRune); l++ {
		t = root
		foundVal := ""
		depth := 0
		for i := 0; i+l < len(originRune); i++ {
			letter := string(originRune[l+i])
			if t.children[letter] == nil {
				break
			}
			if len(t.children[letter].values) > 0 {
				foundVal = t.children[letter].values[0]
				depth = i
			}
			t = t.children[letter]
		}
		if foundVal != "" {
			result += foundVal
			l += depth
		} else {
			result += string(originRune[l : l+1])
		}
	}
	return result
}
