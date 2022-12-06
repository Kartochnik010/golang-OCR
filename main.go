package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/ledongthuc/pdf"
	"github.com/xuri/excelize/v2"
)

var (
	lg        = log.New(os.Stderr, "ERROR:\t", log.Ltime)
	checkDate = regexp.MustCompile(`\d\d\.\d\d.\d\d\d\d`)
)

func main() {
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	IDs, err := GetData(absPath + `/data`)
	if err != nil {
		lg.Fatalln(err)
	}

	f := excelize.NewFile()
	index := f.NewSheet("Удостоверения")

	for i := 0; i < len(IDs); i++ {
		for j := 0; j < len(IDs[i])-4; j++ {
			// if i == 2 && len(IDs[i][j]) != 0 && unicode.IsDigit(rune(IDs[i][j][0])) {}
			setRange := string([]rune{rune('A' + j)}) + strconv.Itoa((i + 1))
			f.SetCellValue("Удостоверения", setRange, IDs[i][j])
		}
	}

	f.SetActiveSheet(index)
	reNew(absPath + "/Book1.xlsx")

	time.Sleep((1 * time.Second) / 2)
	if err := f.SaveAs(absPath + "/Book1.xlsx"); err != nil {
		lg.Println(err)
	}
	os.Remove(absPath + "/data")
}

func GetData(dirPath string) ([][]string, error) {
	data := [][]string{}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return [][]string{}, err
	}
	i := 0
	fmt.Println(" Тип |   Время  |          Имя файла          |  Статус  |")
	os.Mkdir(path.Join(dirPath[:len(dirPath)-4]+"/failed"), 0755)
	os.Mkdir(path.Join(dirPath[:len(dirPath)-4]+"/success"), 0755)
	for _, f := range files {

		content, err := readPdf(dirPath + "/" + f.Name())
		if err != nil {
			lg.Printf(" %-30s skipped", f.Name())
			os.Rename(path.Join(dirPath, "/"+f.Name()), path.Join(dirPath[:len(dirPath)-4]+"/failed/", f.Name()))
			continue
		}
		if len(content) == 0 {
			lg.Printf(" %-30s skipped", f.Name())
			os.Rename(path.Join(dirPath, "/"+f.Name()), path.Join(dirPath[:len(dirPath)-4]+"/failed/", f.Name()))
			continue
		}
		data = append(data, []string{}, []string{})
		data[i] = append(data[i], content...)
		os.Rename(path.Join(dirPath, "/"+f.Name()), path.Join(dirPath[:len(dirPath)-4]+"/success/", f.Name()))
		i++
	}
	return data, nil
}

func reNew(fileName string) {
	_, err := os.ReadFile(fileName)
	if err == nil {
		os.Remove(fileName)
	}
}

func readPdf(path string) ([]string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Println(err)
		}
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

		rows, err := p.GetTextByRow()
		if err != nil {
			return []string{}, err
		}

		for i, row := range rows {
			// println(">>>> row: ", row.Position)
			for _, word := range row.Content {

				if i == 2 && checkDate.MatchString(word.S) {
					data = append(data, "")
				}
				data = append(data, word.S)
			}
		}
	}
	return data, nil
}
