package customtextbreak

import (
	"strings"
	"testing"
)

func TestThaiTextBreak(t *testing.T) {
	tbk := NewThaiTextBreak()
	err := tbk.Load("./thaidict/lexitron.txt")
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := tbk.BreakTextToToken("สวัสดีครับ")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 2 {
		t.Fatal("tokens is empty")
	}
}

func TestThaiTextBreak2(t *testing.T) {
	tbk := NewThaiTextBreak()
	err := tbk.Load("./thaidict/lexitron.txt")
	if err != nil {
		t.Fatal(err)
	}
	text := "มานน มา มา"
	tokens, err := tbk.BreakTextToToken(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 5 {
		t.Fatal("tokens is not 5")
	}
	if strings.Join(tokens, "") != text {
		t.Fatal("tokens not match")
	}
}

func TestThaiTextBreak3(t *testing.T) {
	tbk := NewThaiTextBreak()
	err := tbk.Load("./thaidict/lexitron.txt")
	if err != nil {
		t.Fatal(err)
	}
	text := "ผค 5555 นครราชสีมา"
	tokens, err := tbk.BreakTextToToken(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 4 {
		t.Fatal("tokens is not 4")
	}
	if strings.Join(tokens, "") != text {
		t.Fatal("tokens not match")
	}
}

func TestThaiTextBreak4(t *testing.T) {
	tbk := NewThaiTextBreak()
	err := tbk.Load("./thaidict/lexitron.txt")
	if err != nil {
		t.Fatal(err)
	}
	text := "ผค 555 กท"
	tokens, err := tbk.BreakTextToToken(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 3 {
		t.Fatal("tokens is not 3")
	}
	if strings.Join(tokens, "") != text {
		t.Fatal("tokens not match")
	}
}

func TestThaiTextBreak5(t *testing.T) {
	tbk := NewThaiTextBreak()
	err := tbk.Load("./thaidict/lexitron.txt")
	if err != nil {
		t.Fatal(err)
	}
	text := "ผค 555 กท xx"
	tokens, err := tbk.BreakTextToToken(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 4 {
		t.Fatal("tokens is not 4")
	}
	if strings.Join(tokens, "") != text {
		t.Fatal("tokens not match")
	}
}
