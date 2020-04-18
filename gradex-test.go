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

		for _, pageCount := range N {

			fmt.Printf("\n-------------------------\n%s(%d)\n--------------------\n", getDocName(), pageCount)
			for n := 0; n < pageCount; n = n + 1 {
				fmt.Println(files[ImageIndex%NumSourceImages])
				ImageIndex = ImageIndex + 1
			}

		}

	}

}

func getDocName() string {
	num := fmt.Sprintf("%d", rand.Int())
	return fmt.Sprintf("ENGI01020-B%s", num[0:5])
}

/*


	suffix := filepath.Ext(inputPath)

	// sanity check
	if suffix != ".pdf" {
		fmt.Printf("Error: input path must be a .pdf\n")
		os.Exit(1)
	}

	// need page count to find the jpeg files again later
	numPages, err := countPages(inputPath)

	// render to images
	jpegPath := "./jpg"
	err = ensureDir(jpegPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	basename := strings.TrimSuffix(inputPath, suffix)
	jpegFileOption := fmt.Sprintf("%s/%s%%04d.jpg", jpegPath, basename)

	err = convertPDFToJPEGs(inputPath, jpegPath, jpegFileOption)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// convert images to individual pdfs, with form overlay

	pagePath := "./pdf"
	err = ensureDir(pagePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)

	mergePaths := []string{}

	// gs starts indexing at 1
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		// construct image name
		previousImagePath := fmt.Sprintf(jpegFileOption, imgIdx)
		pageFilename := fmt.Sprintf(pageFileOption, imgIdx)

		//TODO select Layout to suit landscape or portrait
		svgLayoutPath := "./test/layout-312pt-static-mark-dynamic-moderate-comment-static-check.svg"

		err := parsesvg.RenderSpread(svgLayoutPath, spreadName, previousImagePath, imgIdx, pageFilename)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		//save the pdf filename for the merge at the end
		mergePaths = append(mergePaths, pageFilename)
	}

	outputPath := fmt.Sprintf("%s-%s.pdf", basename, spreadName)
	err = mergePdf(mergePaths, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
*/
