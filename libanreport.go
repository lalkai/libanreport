package libanreport

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/oneplus1000/errord"
	"github.com/signintech/pdft/render"
)

//go:embed "customtextbreak/thaidict/lexitron.txt"
var fdLexitron embed.FS

// TypeText text
const TypeText = 1

// TypeIMG img
const TypeIMG = 2

var ErrFileNotFound = errors.New("file not found")

// สร้าง PDF จากเทมเพลตและข้อมูล
func GenPdf(
	tmpl Tmpl,
	datas []DataJSON,
	fontoverrides []FontOverrideJSON,
	w io.Writer,
) error {

	var finfos render.FieldInfos
	for _, f := range tmpl.tmplJSON.Fields {
		var finfo render.FieldInfo
		var styleDefine *FieldJSON
		if f.Style != nil {
			index := tmpl.tmplJSON.stylesIndexByName(*f.Style)
			if index < 0 {
				return errors.New("Fields.Style not found")
			}
			styleDefine = &tmpl.tmplJSON.Styles[index].Define
		}

		if styleDefine != nil {
			bindFieldInfo(styleDefine, &finfo)
		}
		bindFieldInfo(&f, &finfo)
		finfos = append(finfos, finfo)
	}

	rdr, err := newRender(tmpl.tmplPDFPath, finfos)
	if err != nil {
		return errord.Errorf("newRender error: %w", err)
	}

	//add font
	for _, f := range tmpl.tmplJSON.Fonts {
		fontpath := filepath.Join(tmpl.tmplFolderPath, f.TTF)
		err = rdr.AddFont(f.Name, fontpath)
		if err != nil {
			return errord.Errorf("AddFont error: %w", err)
		}
	}

	//start setup override
	for _, fontoverride := range fontoverrides {
		rdr.TextriseOverride(fontoverride.Name, func(
			leftRune rune,
			rightRune rune,
			fontSize int,
			allText string,
			currTextIndex int,
		) float32 {
			gap := float32(0)
			runes := []rune(allText)
			var nextRune rune
			if len(runes) > currTextIndex+1 {
				nextRune = runes[currTextIndex+1]
			}
			for _, textrise := range fontoverride.Textrises {
				if isInRune(leftRune, textrise.Lefts) && isInRune(rightRune, textrise.Rights) && isInRune(nextRune, textrise.Nexts) {
					gap = textrise.Val
					break
				}
			}
			return gap * float32(fontSize) / 100
		})

		k := int64(0)
		lastK := int64(0)
		lastPairVal := int16(0)
		rdr.KernOverride(fontoverride.Name, func(
			leftRune rune,
			rightRune rune,
			leftPair uint,
			rightPair uint,
			pairVal int16,
		) int16 {

			for _, kern := range fontoverride.Kernings {
				if isInRune(leftRune, kern.Lefts) && isInRune(rightRune, kern.Rights) {
					pairVal = kern.Val
					lastPairVal = kern.Val
					lastK = k
				}
			}
			if k-1 == lastK { //ทำเพื่อโยกตัวหลังจากทีี่ถูกปรับ  kerning ให้ไปวางที่เดิมไม่ถูกบีบ
				pairVal = (-1) * lastPairVal
			}
			k++
			return pairVal
		})
	}
	//end setup override

	//bind data
	for _, d := range datas {
		if d.Type == TypeText {
			err = rdr.Text(d.Key, d.Val)
			if err != nil {
				return errord.Errorf("Error on key %s: %s", d.Key, err)
			}
		} else if d.Type == TypeIMG {
			err = rdr.ImgBase64(d.Key, d.Val)
			if err != nil {
				return errord.Errorf("Error on key %s: %s", d.Key, err)
			}
		}
	}

	//save pdf
	err = rdr.SaveTo(w)
	if err != nil {
		return errord.Errorf("rd.SaveTo error: %v", err)
	}

	return nil
}

// path ไปยัง fontoverride.json
func ReadfontsJSON(path string) ([]FontOverrideJSON, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrFileNotFound
	}
	var fontoverrides []FontOverrideJSON
	err = json.Unmarshal(d, &fontoverrides)
	if err != nil {
		return nil, errord.Errorf("error unmarshalling font overrides: %w", err)
	}
	return fontoverrides, nil
}

func ReadTmplDir(path string) (Tmpl, error) {

	folderName := filepath.Base(path)
	tmplJsonPath := filepath.Join(path, "tmpl.json")
	tmplPdfPath := filepath.Join(path, "tmpl.pdf")

	tmplJson, err := parseTmplJSON(tmplJsonPath)
	if err != nil {
		return Tmpl{}, errord.Errorf("error reading tmpl.json: %w", err)
	}
	return Tmpl{
		tmplFolderPath: path,
		folderName:     folderName,
		tmplJSON:       &tmplJson,
		tmplPDFPath:    tmplPdfPath,
	}, nil
}

func RemoveSpecialRune(txt string) string {
	const noBreakSpace = '\u00A0'
	var buff bytes.Buffer
	for _, r := range txt {
		if r == noBreakSpace {
			continue
		}
		buff.WriteRune(r)
	}
	return buff.String()
}

// --- private ---

func parseTmplJSON(path string) (TmplJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TmplJSON{}, err
	}
	var objs TmplJSON
	err = json.Unmarshal(data, &objs)
	if err != nil {
		return TmplJSON{}, err
	}
	return objs, nil
}

func isInRune(r rune, vals []string) bool {
	if vals == nil || len(vals) <= 0 {
		return true
	}
	str := string(r)
	for _, val := range vals {
		if str == val {
			return true
		}
	}
	return false
}
