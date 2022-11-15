// package main

// import (
// 	"fmt"

// 	"github.com/otiai10/gosseract/v2"
// )

// func main() {
// 	client := gosseract.NewClient()
// 	defer client.Close()
// 	client.SetImage("data/ilyas-1.png")
// 	text, _ := client.Text()
// 	fmt.Println(text)
// 	// Hello, World!
// }

// package main

// import (
// 	"fmt"
// 	"log"

// 	"code.sajari.com/docconv"
// )

// func main() {
// 	res, err := docconv.ConvertPath("data/ilyas.pdf")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(res)
// }

// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"log"

// 	"github.com/ledongthuc/pdf"
// )

// func main() {
// 	content, err := readPdf("data/ilyas.pdf") // Read local pdf file
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(content)
// 	return
// }

// func readPdf(path string) (string, error) {
// 	_, r, err := pdf.Open(path)
// 	if err != nil {
// 		return "", err
// 	}
// 	totalPage := r.NumPage()

// 	var textBuilder bytes.Buffer
// 	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
// 		p := r.Page(pageIndex)
// 		if p.V.IsNull() {
// 			continue
// 		}
// 		text, err := p.GetPlainText("\n")
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		textBuilder.WriteString(text)
// 	}
// 	return textBuilder.String(), nil
// }

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"

	"github.com/ledongthuc/pdf"
	"github.com/xuri/excelize/v2"
)

var (
	lg = log.New(os.Stderr, "ERROR:\t", log.Lshortfile)
)

func main() {
	// for i := 0; i < 5; i++ {
	// 	for j := 0; j < 3; j++ {
	// 		fmt.Print(string([]rune{rune('A' + i)}) + strconv.Itoa((j + 1)))
	// 	}
	// 	fmt.Println()
	// }

	IDs, err := GetData("data")
	if err != nil {
		lg.Fatalln(err)
	}

	f := excelize.NewFile()
	index := f.NewSheet("Удостоверения")
outerLoop:
	for i := 0; i < len(IDs); i++ {
		for j := 0; j < len(IDs[i]); j++ {
			err := f.SetCellValue("Удостоверения", string([]rune{rune('A' + j)})+strconv.Itoa((i+1)), IDs[i][j])
			if err != nil {

				continue outerLoop
			}
		}
	}

	f.SetActiveSheet(index)
	_, err = os.ReadFile("Book1.xlsx")
	if err == nil {
		os.Remove("Book1.xlsx")
	}
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func GetData(dirPath string) ([][]string, error) {
	data := [][]string{}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return [][]string{}, err
	}

	for i, f := range files {
		data = append(data, []string{})

		content, err := readPdf(dirPath + "/" + f.Name())
		if err != nil {
			// lg.Printf("File %s skipped: %s", f.Name(), err)
			continue
		}

		data[i] = append(data[i], content...)
	}
	return data, nil
}

func readPdf(path string) ([]string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return []string{}, err
	}
	totalPage := r.NumPage()
	data := []string{}
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()

		for i, row := range rows {
			// println(">>>> row: ", row.Position)
			for _, word := range row.Content {
				if i == 2 && len(word.S) != 0 && unicode.IsDigit(rune(word.S[0])) {
					data = append(data, "asodk")
				} else {
					data = append(data, word.S)
				}

			}
		}
	}
	return data, nil
}
