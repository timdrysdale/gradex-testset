package main

/*
 * Add a cover page to a PDF file
 * Generates cover page then merges, including form field data (AcroForms).
 *
 * Run as: gradex-coverpage <barefile>.pdf
 *
 * outputs: <barefile>-covered.pdf (using internally generated cover page)
 *
 * Adapted from github.com/unidoc/unipdf-examples/pages/pdf_merge_advanced.go
 *
 *
 */

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bsipos/thist"
	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelInfo))
}

func getPageCount(PageCountMean, PageCountStdDev float64) int {
	return int(math.Round(rand.NormFloat64()*PageCountStdDev + PageCountMean))
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Requires three arguments: NumScripts PageCountMean PageCountStdDev\n")
		fmt.Printf("Example usage: gradex-testset 10 15 2\n")
		os.Exit(0)
	}

	NumScriptsInt64, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing NumScripts %s\n", err)
		os.Exit(1)
	}
	NumScripts := int(NumScriptsInt64)

	PageCountMean, err := strconv.ParseFloat(os.Args[2], 64)
	PageCountStdDev, err := strconv.ParseFloat(os.Args[3], 64)
	rand.Seed(time.Now().UnixNano())
	// initialise random number generator

	h := thist.NewHist(nil, "Page count", "fixed", 10, false)

	N := []int{}

	for n := 0; n < NumScripts; n = n + 1 {
		pc := getPageCount(PageCountMean, PageCountStdDev)
		h.Update(float64(pc))
		fmt.Println(h.Draw())
		N = append(N, pc)
	}

	fmt.Println("pagecounts")
	for _, nn := range N {
		fmt.Printf(" %d", nn)
	}
	fmt.Println("")

	doit := confirm("Make it so, Picard?", 1)

	if doit {
		fmt.Println("cooking now ...")

		var files []string

		root := "./iam/data/forms/"
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, "-clean.jpg") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}

		NumSourceImages := len(files)
		ImageIndex := 0

		pagePath := "./pdf"
		err = ensureDir(pagePath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		for _, pageCount := range N {

			// set up doc
			c := creator.New()
			c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
			c.SetPageSize(creator.PageSizeA4)

			// add image pages
			for n := 0; n < pageCount; n = n + 1 {

				imgPath := (files[ImageIndex%NumSourceImages])
				ImageIndex = ImageIndex + 1

				c.NewPage()
				img, err := c.NewImageFromFile(imgPath)
				if err != nil {
					fmt.Printf("Error opening image file %s: %v", imgPath, err)
					os.Exit(1)
				}
				//img.SetHeight(creator.PageSizeA4.Height)
				img.ScaleToWidth(210 * creator.PPMM)
				img.SetPos(0, 0)
				c.Draw(img) //draw previous image

			}
			/*		if true {
						fmt.Printf("\n-------------------------\n%s(%d)\n--------------------\n", getDocName(), pageCount)
						for n := 0; n < pageCount; n = n + 1 {
							fmt.Println(files[ImageIndex%NumSourceImages])
							ImageIndex = ImageIndex + 1
						}
					}
			*/

			c.SetOptimizer(optimize.New(optimize.Options{
				CombineDuplicateDirectObjects:   true,
				CombineIdenticalIndirectObjects: true,
				CombineDuplicateStreams:         true,
				CompressStreams:                 true,
				UseObjectStreams:                true,
				ImageQuality:                    90,
				ImageUpperPPI:                   175,
			}))

			// write to memory
			outPath := fmt.Sprintf("%s/%s", pagePath, getDocName())
			c.WriteToFile(outPath)

		}
	}
}

func getDocName() string {
	num := fmt.Sprintf("%d", rand.Int())
	return fmt.Sprintf("ENGI01020-B%s.pdf", num[0:5])
}
