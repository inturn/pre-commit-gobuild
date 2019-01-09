package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/inturn/pre-commit-gobuild/internal/helpers"
)

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dirs := helpers.DirsWith(workDir, "\\.go$")

	wg := &sync.WaitGroup{}

	for _, dir := range dirs {
		if !strings.Contains(dir, "/vendor/") {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Printf("error occured on read dir %s: %s", dir, err)
			}
			for _, f := range files {
				if !strings.HasSuffix(f.Name(), ".go") {
					continue
				}
				wg.Add(1)
				go func(d, name string) {
					sortFileImports(filepath.Join(d, name))
					wg.Done()
				}(dir, f.Name())
			}
		}
	}
	wg.Wait()
}

func sortFileImports(path string) {
	fSet := token.NewFileSet()

	f, err := parser.ParseFile(fSet, path, nil, parser.ParseComments)
	if err != nil {
		log.Println(err)
		return
	}

	sortImports(f)

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	if err := file.Truncate(0); err != nil {
		log.Println(err)
		return
	}

	if err := format.Node(file, fSet, f); err != nil {
		log.Println(err)
		return
	}
}

func sortImports(f *ast.File) {
	imp1 := make(impSlice, 0)
	imp2 := make(impSlice, 0)

	for _, imp := range f.Imports {
		var v, n string
		v = imp.Path.Value
		if imp.Name != nil {
			n = imp.Name.Name
		}

		// determine third-party package import
		if strings.Contains(v, ".") {
			imp2 = append(imp2, importData{v, n})
			continue
		}

		imp1 = append(imp1, importData{v, n})
	}

	sort.Sort(imp1)
	sort.Sort(imp2)

	for _, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Imports are always first.
			break
		}

		if !d.Lparen.IsValid() {
			// Not a block: sorted by default.
			continue
		}

		d.Specs = d.Specs[:0]

		for _, imp := range imp1 {
			iSpec := ast.ImportSpec{
				Path: &ast.BasicLit{Value: imp.value},
				Name: &ast.Ident{Name: imp.name},
			}
			d.Specs = append(d.Specs, &iSpec)
		}

		if len(imp2) != 0 {
			// Add empty line between groups.
			d.Specs = append(d.Specs, &ast.ImportSpec{Path: &ast.BasicLit{}})

			for _, imp := range imp2 {
				iSpec := ast.ImportSpec{
					Path: &ast.BasicLit{Value: imp.value},
					Name: &ast.Ident{Name: imp.name},
				}
				d.Specs = append(d.Specs, &iSpec)
			}
		}
	}
}

type impSlice []importData

type importData struct {
	value string
	name  string
}

func (s impSlice) Len() int {
	return len(s)
}

func (s impSlice) Less(i, j int) bool {
	return s[i].value < s[j].value
}

func (s impSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}