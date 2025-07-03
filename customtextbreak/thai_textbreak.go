package customtextbreak

import (
	"bufio"
	"io"
	"os"
	"regexp"

	"github.com/armon/go-radix"
	"github.com/oneplus1000/errord"
)

type ThaiTextBreak struct {
	tree                      *radix.Tree
	latinRegexp, numberRegexp *regexp.Regexp
}

func NewThaiTextBreak() *ThaiTextBreak {
	return &ThaiTextBreak{
		tree: radix.New(),
	}
}

func (t ThaiTextBreak) BreakTextToToken(text string) ([]string, error) {
	return t.breakTextToToken(text)
}

func (t ThaiTextBreak) breakTextToToken(text string) ([]string, error) {
	i := 0
	last := len(text)
	tokens := make([]string, 0)
	var notfound []byte
	for i < last {
		found := t.search(text[i:last])
		if found == "" {
			notfound = append(notfound, text[i])
			i++
			continue
		}
		if len(notfound) > 0 {
			tokens = append(tokens, string(notfound))
			notfound = nil
		}
		tokens = append(tokens, found)
		i += len(found)
	}
	if len(notfound) > 0 {
		tokens = append(tokens, string(notfound))
		notfound = nil
	}
	return tokens, nil
}

func (t ThaiTextBreak) search(sub string) string {
	str := t.latinRegexp.FindString(sub)
	if str != "" {
		return str
	}
	str = t.numberRegexp.FindString(sub)
	if str != "" {
		return str
	}

	str, _, _ = t.tree.LongestPrefix(sub)
	return str
}

func (t *ThaiTextBreak) Load(dictPath string) error {
	fd, err := os.Open(dictPath)
	if err != nil {
		return errord.Errorf("Error open dict file: %s", err)
	}
	defer fd.Close()

	err = t.LoadFromReader(fd)
	if err != nil {
		return errord.Errorf("Error load dict file: %s", err)
	}
	return nil
}

func (t *ThaiTextBreak) LoadFromReader(fd io.Reader) error {
	scr := bufio.NewScanner(fd)
	for scr.Scan() {
		t.tree.Insert(scr.Text(), 1)
	}

	latinRegexp, err := regexp.Compile(`[A-Za-z\d]*`)
	if err != nil {
		return errord.Errorf("Error compile latin regexp: %s", err)
	}
	t.latinRegexp = latinRegexp

	numberRegexp, err := regexp.Compile(`[\d]*`)
	if err != nil {
		return errord.Errorf("Error compile number regexp: %s", err)
	}
	t.numberRegexp = numberRegexp
	return nil
}
