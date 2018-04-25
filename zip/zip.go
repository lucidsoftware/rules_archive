package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("zip", "Create a zip archive from files, directories, or other archives. Strips timestamps. Preserves Unix permissions.")
	archives := parser.List("a", "archive", &argparse.Options{
		Help: "Zip archive to merge. Add to the root by default; use name=path to add files to name instead.",
	})
	files := parser.List("f", "file", &argparse.Options{
		Help: "File or directory of files to add. Adds using the specified path; use name=path to add files to name instead.",
	})
	output := parser.String("o", "output", &argparse.Options{Default: "-", Help: "Archive output"})
	compress := parser.Flag("x", "compress", &argparse.Options{Help: "Deflate contents"})
	err := parser.Parse(os.Args)
	if err != nil {
		println(parser.Usage(err))
		os.Exit(2)
	}
	err = run(*archives, *files, *output, *compress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(archives []string, files []string, output string, compress bool) error {
	method := zip.Store
	if compress {
		method = zip.Deflate
	}
	entries, err := createEntries(archives, files)
	if err != nil {
		return err
	}
	return write(output, entries, method)
}

type entry interface {
	Name() string
	Write(*zip.Writer, uint16) error
}

func createEntries(archives []string, files []string) ([]entry, error) {
	var err error

	var entries []entry
	for _, path := range archives {
		var name string
		parts := strings.SplitN(path, "=", 2)
		if len(parts) == 1 {
			name = ""
		} else {
			name = parts[0]
			path = parts[1]
		}
		entries, err = appendArchive(entries, name, path)

		if err != nil {
			return entries, err
		}
	}

	for _, path := range files {
		var name string
		parts := strings.SplitN(path, "=", 2)
		if len(parts) == 1 {
			name = parts[0]
		} else {
			name = parts[0]
			path = parts[1]
		}
		entries, err = appendFiles(entries, name, path)

		if err != nil {
			return entries, err
		}
	}

	oldEntries := entries
	entries = nil
	names := map[string]bool{}
	for _, entry := range oldEntries {
		if _, ok := names[entry.Name()]; !ok {
			entries = append(entries, entry)
			names[entry.Name()] = true
		}
	}
	names = nil

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return entries, nil
}

func write(output string, entries []entry, method uint16) error {
	var err error
	var outputFile io.WriteCloser
	if output == "-" {
		outputFile = os.Stdout
	} else {
		outputFile, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("Could not open output %v: %v", output, err)
		}
	}
	defer outputFile.Close()

	writer := zip.NewWriter(outputFile)
	defer writer.Close()
	for _, entry := range entries {
		err := entry.Write(writer, method)
		if err != nil {
			return err
		}
	}

	return nil
}

func appendArchive(entries []entry, name string, path string) ([]entry, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return entries, err
	}
	defer reader.Close()
	for _, file := range reader.File {
		entries = append(entries, archiveEntry{name: name + file.Name, archivePath: path, path: file.Name})
	}
	return entries, nil
}

type archiveEntry struct {
	name        string
	archivePath string
	path        string
}

func (entry archiveEntry) Name() string {
	return entry.name
}

func (entry archiveEntry) Write(writer *zip.Writer, method uint16) error {
	reader, err := zip.OpenReader(entry.archivePath)

	var sourceFile *zip.File
	for _, file := range reader.File {
		if file.Name == entry.path {
			sourceFile = file
			break
		}
	}
	if err != nil {
		return err
	}

	header := sourceFile.FileHeader
	header.Name = entry.Name()
	header.Method = method
	header.Modified = time.Unix(1262329200, 0)

	fileWriter, err := writer.CreateHeader(&header)
	if err != nil {
		return err
	}

	file, err := sourceFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return fmt.Errorf("Could not copy %v in %v: %v", entry.path, entry.archivePath, err)
	}

	return nil
}

func appendFiles(entries []entry, name string, path string) ([]entry, error) {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to traverse %v: %v", filePath, err)
		}
		if info.IsDir() {
			return nil
		}
		entries = append(entries, fileEntry{name: name + filePath[len(path):len(filePath)], Path: filePath})
		return err
	})
	return entries, err
}

type fileEntry struct {
	name string
	Path string
}

func (entry fileEntry) Name() string {
	return entry.name
}

func (entry fileEntry) Write(writer *zip.Writer, method uint16) error {
	file, err := os.Open(entry.Path)
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Could not open file %v: %v", entry.Path, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = entry.Name()
	header.Method = method
	header.Modified = time.Unix(1262329200, 0)

	fileWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return fmt.Errorf("Could not copy file %v: %v", entry.Path, err)
	}

	return nil
}
