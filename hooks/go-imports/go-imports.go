package main

import (
	"bytes"
	"errors"
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

type lintError struct {
	err  error
	path string
}

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dirs := helpers.DirsWith(workDir, "\\.go$")

	errc := make(chan lintError, 10)
	wg := &sync.WaitGroup{}

	go func() {
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
						sortFileImports(filepath.Join(d, name), errc)
						wg.Done()
					}(dir, f.Name())
				}
			}
		}
		wg.Wait()
		close(errc)
	}()

	var le []lintError
	for lintErr := range errc {
		log.Println(lintErr.path, lintErr.err)
		le = append(le, lintErr)
	}
	if len(le) != 0 {
		log.Println("files changed:", len(le))
		os.Exit(1)
	}
}

func sortFileImports(path string, errc chan<- lintError) {
	fSet := token.NewFileSet()

	f, err := parser.ParseFile(fSet, path, nil, parser.ParseComments)
	if err != nil {
		errc <- lintError{
			err:  err,
			path: path,
		}
		return
	}

	sortImports(f)

	buf := &bytes.Buffer{}
	if err := format.Node(buf, fSet, f); err != nil {
		errc <- lintError{
			err:  err,
			path: path,
		}
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		errc <- lintError{
			err:  err,
			path: path,
		}
		return
	}

	if buf.String() == string(data) {
		return
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0664); err != nil {
		errc <- lintError{
			err:  err,
			path: path,
		}
		return
	}
	errc <- lintError{
		err:  errors.New("file has changed"),
		path: path,
	}
}

func sortImports(f *ast.File) {
	if len(f.Imports) <= 1 {
		return
	}

	imp1 := make(impSlice, 0)
	imp2 := make(impSlice, 0)

	for _, imp := range f.Imports {
		impData := importData{}

		if imp.Doc != nil && imp.Name != nil && imp.Name.Name == "_" {
			impData.comment = imp.Doc.Text()
		}

		if imp.Name != nil {
			impData.name = imp.Name.Name
		}
		impData.value = imp.Path.Value

		// determine third-party package import
		if strings.Contains(imp.Path.Value, ".") {
			imp2 = append(imp2, impData)
			continue
		}

		imp1 = append(imp1, impData)
	}

	nonImportComment := f.Comments[:0]
	startPos := f.Imports[0].Pos()
	lastPos := f.Imports[len(f.Imports)-1].End()

	for _, c := range f.Comments {
		if c.Pos() > lastPos || c.Pos() < startPos {
			nonImportComment = append(nonImportComment, c)
		}
	}

	f.Comments = nonImportComment

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
			addISpec(imp, d)
		}

		if len(imp2) != 0 {
			// Add empty line between groups.
			d.Specs = append(d.Specs, &ast.ImportSpec{Path: &ast.BasicLit{}})

			for _, imp := range imp2 {
				addISpec(imp, d)

			}
		}
	}
}

func addISpec(imp importData, d *ast.GenDecl) {
	if imp.name == "_" {
		comm := imp.comment
		if comm == "" {
			comm = "todo comment here why do you use blank import"
		}
		d.Specs = append(d.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{Value: "// " + strings.TrimSpace(comm)},
		})
	}
	iSpec := ast.ImportSpec{
		Path: &ast.BasicLit{Value: imp.value},
		Name: &ast.Ident{Name: imp.name},
	}
	d.Specs = append(d.Specs, &iSpec)
}

type impSlice []importData

type importData struct {
	value   string
	name    string
	comment string
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
