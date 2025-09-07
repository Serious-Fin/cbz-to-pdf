package main

import (
	"fmt"
	"image"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/gen2brain/go-unarr"
	"github.com/signintech/gopdf"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exePath)

	files, err := getArchiveFileNames(exeDir)
	if err != nil {
		log.Fatalf("Error reading files from current directory: %s", err.Error())
	}

	for _, archiveFile := range files {
		pdfName := archiveFile[:len(archiveFile)-len(filepath.Ext(archiveFile))] + ".pdf"
		archivePath := filepath.Join(exeDir, archiveFile)
		pdfPath := filepath.Join(exeDir, pdfName)
		err = createPdfFromArchive(archivePath, pdfPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func createPdfFromArchive(archivePath, pdfName string) error {
	extractedArchivePath, err := extractArchive(archivePath)
	if err != nil {
		return fmt.Errorf("error extracting archive: %w", err)
	}
	defer os.RemoveAll(extractedArchivePath)

	imageFiles := getAllImagesInDir(extractedArchivePath)
	sortFileNamesAsc(imageFiles)
	err = createPdfFromImages(imageFiles, pdfName)
	if err != nil {
		return fmt.Errorf("error creating PDF: %w", err)
	}
	return nil
}

func getArchiveFileNames(dir string) ([]string, error) {
	cbzFiles := make([]string, 0)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory contents: %w", err)
	}

	for _, file := range files {
		filename := file.Name()
		if isCbzFile(filename) || isCbrFile(filename) {
			cbzFiles = append(cbzFiles, filename)
		}
	}
	return cbzFiles, nil
}

var cbzFileRegex = regexp.MustCompile(`.*\.[Cc][Bb][Zz]$`)

func isCbzFile(filename string) bool {
	matches := cbzFileRegex.FindStringSubmatch(filename)
	return len(matches) == 1
}

var cbrFileRegex = regexp.MustCompile(`.*\.[Cc][Bb][Rr]$`)

func isCbrFile(filename string) bool {
	matches := cbrFileRegex.FindStringSubmatch(filename)
	return len(matches) == 1
}

func extractArchive(filename string) (string, error) {
	a, err := unarr.NewArchive(filename)
	if err != nil {
		return "", fmt.Errorf("error creating unarr Archine object from \"%s\": %w", filename, err)
	}
	defer a.Close()

	extractLocation := fmt.Sprintf("./folder_%d", random.Int())
	_, err = a.Extract(extractLocation)
	if err != nil {
		return "", fmt.Errorf("can not extract archive from \"%s\": %w", filename, err)
	}
	return extractLocation, nil
}

func getAllImagesInDir(dir string) []string {
	imageFiles := make([]string, 0)
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if isJpgFile(path) || isPngFile(path) {
			imageFiles = append(imageFiles, path)
		}
		return nil
	})
	return imageFiles
}

var jpgRegex = regexp.MustCompile(`.*\.[Jj][Pp][Gg]$`)

func isJpgFile(filename string) bool {
	matches := jpgRegex.FindStringSubmatch(filename)
	return len(matches) == 1
}

var pngRegex = regexp.MustCompile(`.*\.[Pp][Nn][Gg]$`)

func isPngFile(filename string) bool {
	matches := pngRegex.FindStringSubmatch(filename)
	return len(matches) == 1
}

func sortFileNamesAsc(filenames []string) {
	isSortedInsensitive := slices.IsSortedFunc(filenames, func(x, y string) int {
		return strings.Compare(strings.ToLower(x), strings.ToLower(y))
	})
	if !isSortedInsensitive {
		sort.Slice(filenames, func(i, j int) bool {
			return strings.ToLower(filenames[i]) < strings.ToLower(filenames[j])
		})
	}
}

func createPdfFromImages(imageFiles []string, pdfPath string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, imageFile := range imageFiles {
		w, h, err := imageSize(imageFile)
		if err != nil {
			return fmt.Errorf("could not get image dimensions from file: %w", err)
		}

		pdf.AddPageWithOption(gopdf.PageOption{
			PageSize: &gopdf.Rect{W: w, H: h},
		})

		err = pdf.Image(imageFile, 0, 0, &gopdf.Rect{W: w, H: h})
		if err != nil {
			return fmt.Errorf("could not add image to PDF: %w", err)
		}
	}

	err := pdf.WritePdf(pdfPath)
	if err != nil {
		return fmt.Errorf("could not write to PDF: %w", err)
	}
	return nil
}

func imageSize(path string) (w, h float64, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return float64(cfg.Width), float64(cfg.Height), nil
}
