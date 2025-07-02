package main

import (
	"os"
	"path/filepath"
)

func Generate(outputPath string, structs []StructInfo) error {
	for _, s := range structs {
		outputPath := filepath.Join(outputPath, s.Name+"_deltagen.go")
		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := templates.ExecuteTemplate(f, "file", s); err != nil {
			return err
		}
	}

	return nil
}
