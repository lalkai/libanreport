package libanreport

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/oneplus1000/errord"
)

func TestGenPdf(t *testing.T) {
	err := testGenPdf()
	if err != nil {
		t.Errorf("testGenPdf failed: %v", err)
	}
}

func testGenPdf() error {

	//สร้างโฟลเดอร์สำหรับเก็บผลลัพธ์ ถ้ายังไม่มี
	if _, err := os.Stat("testing_output"); os.IsNotExist(err) {
		err = os.Mkdir("testing_output", 0755)
		if err != nil {
			return errord.Errorf("error creating output directory: %w", err)
		}
	}

	//code จริงเริ่มที่นี่
	outputPath := filepath.Join("testing_output", "test04.pdf")
	tmplPath := filepath.Join("testing", "test04.tmpl")
	tmpl, err := ReadTmplDir(tmplPath)
	if err != nil {
		return errord.Errorf("error reading template directory: %w", err)
	}
	fontOverrides, err := ReadfontsJSON(filepath.Join("testing", "fontoverride.json"))
	if err != nil {
		return errord.Errorf("error reading fontoverride.json: %w", err)
	}

	data := []DataJSON{
		{Type: 1, Key: "car_model", Val: "123"},
		{Type: 1, Key: "car_register", Val: "กั้น"},
	}

	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errord.Errorf("error opening output file: %w", err)
	}

	err = GenPdf(
		tmpl,
		data,
		fontOverrides,
		f,
	)
	if err != nil {
		return fmt.Errorf("error generating PDF: %w", err)
	}
	return nil
}
