package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var tRunNamePattern = regexp.MustCompile(`^[\p{Cyrillic}A-Za-z0-9]+(?:_[\p{Cyrillic}A-Za-z0-9]+)*$`)
var formatVerbPattern = regexp.MustCompile(`%[-+#0-9.*\[\]]*[a-zA-Z]`)

type finding struct {
	path    string
	line    int
	message string
}

func main() {
	findings, err := run(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка projectlint: %v\n", err)
		os.Exit(1)
	}
	if len(findings) == 0 {
		return
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].path != findings[j].path {
			return findings[i].path < findings[j].path
		}
		if findings[i].line != findings[j].line {
			return findings[i].line < findings[j].line
		}
		return findings[i].message < findings[j].message
	})

	for _, f := range findings {
		fmt.Printf("%s:%d: %s\n", filepath.ToSlash(f.path), f.line, f.message)
	}
	os.Exit(1)
}

func run(root string) ([]finding, error) {
	var findings []finding

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".github", ".agents", ".kilo", ".kilocode", "vendor":
				if path != root {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			fileFindings, err := lintTestFile(path)
			if err != nil {
				return err
			}
			findings = append(findings, fileFindings...)
			return nil
		}

		fileFindings, err := lintSourceFile(path)
		if err != nil {
			return err
		}
		findings = append(findings, fileFindings...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return findings, nil
}

func lintSourceFile(path string) ([]finding, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var findings []finding

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				if !d.Name.IsExported() {
					continue
				}
				if !hasDocForName(d.Doc, d.Name.Name) {
					findings = append(findings, newFinding(fset, path, d.Pos(), "у экспортируемой функции должен быть комментарий, начинающийся с её имени"))
				}
				continue
			}

			if !d.Name.IsExported() {
				continue
			}
			if !hasExportedReceiver(d.Recv) {
				continue
			}
			if !hasDocForName(d.Doc, d.Name.Name) {
				findings = append(findings, newFinding(fset, path, d.Pos(), "у экспортируемого метода должен быть комментарий, начинающийся с его имени"))
			}

		case *ast.GenDecl:
			if d.Tok != token.TYPE {
				continue
			}
			for _, spec := range d.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || !typeSpec.Name.IsExported() {
					continue
				}

				doc := typeSpec.Doc
				if doc == nil {
					doc = d.Doc
				}
				if !hasDocForName(doc, typeSpec.Name.Name) {
					findings = append(findings, newFinding(fset, path, typeSpec.Pos(), "у экспортируемого типа должен быть комментарий, начинающийся с его имени"))
				}
			}
		}
	}

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		message, ok := extractMessageLiteral(call)
		if !ok {
			return true
		}
		if !requiresRussianMessage(message) {
			return true
		}
		if !containsCyrillic(message) {
			findings = append(findings, newFinding(fset, path, call.Pos(), "сообщение ошибки или лога должно быть на русском языке"))
		}
		return true
	})

	return findings, nil
}

func lintTestFile(path string) ([]finding, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var findings []finding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if !isTRunCall(call) || len(call.Args) < 1 {
			return true
		}

		nameLit, ok := call.Args[0].(*ast.BasicLit)
		if !ok {
			return true
		}
		if nameLit.Kind != token.STRING {
			return true
		}

		name, err := parseStringLiteral(nameLit.Value)
		if err != nil {
			findings = append(findings, newFinding(fset, path, call.Pos(), "не удалось разобрать имя под-теста t.Run"))
			return true
		}

		if !tRunNamePattern.MatchString(name) || strings.Contains(name, ".") || strings.Contains(name, " ") {
			findings = append(findings, newFinding(fset, path, call.Pos(), "имя под-теста t.Run должно быть на русском или с допустимыми идентификаторами, слова разделяются символом '_' без пробелов и точек"))
		}
		return true
	})

	return findings, nil
}

func hasDocForName(doc *ast.CommentGroup, name string) bool {
	if doc == nil {
		return false
	}
	text := strings.TrimSpace(doc.Text())
	return strings.HasPrefix(text, name+" ")
}

func hasExportedReceiver(recv *ast.FieldList) bool {
	if recv == nil || len(recv.List) == 0 {
		return false
	}

	expr := recv.List[0].Type
	for {
		starExpr, ok := expr.(*ast.StarExpr)
		if !ok {
			break
		}
		expr = starExpr.X
	}

	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.IsExported()
}

func extractMessageLiteral(call *ast.CallExpr) (string, bool) {
	if len(call.Args) == 0 {
		return "", false
	}
	if !isErrorOrLogCall(call) {
		return "", false
	}

	arg, ok := call.Args[0].(*ast.BasicLit)
	if !ok || arg.Kind != token.STRING {
		return "", false
	}

	value, err := strconv.Unquote(arg.Value)
	if err != nil {
		return "", false
	}
	return value, true
}

func isErrorOrLogCall(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		if pkg, ok := fun.X.(*ast.Ident); ok {
			if pkg.Name == "errors" && fun.Sel.Name == "New" {
				return true
			}
			if pkg.Name == "fmt" && fun.Sel.Name == "Errorf" {
				return true
			}
		}

		switch fun.Sel.Name {
		case "Debug", "Info", "Warn", "Error", "Fatal", "Panic", "Debugf", "Infof", "Warnf", "Errorf", "Fatalf", "Panicf", "Print", "Printf", "Println":
			return true
		}
	}
	return false
}

func requiresRussianMessage(message string) bool {
	message = formatVerbPattern.ReplaceAllString(message, "")
	for _, r := range message {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func containsCyrillic(message string) bool {
	for _, r := range message {
		if unicode.In(r, unicode.Cyrillic) {
			return true
		}
	}
	return false
}

func isTRunCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return selector.Sel != nil && selector.Sel.Name == "Run"
}

func parseStringLiteral(lit string) (string, error) {
	if len(lit) < 2 {
		return "", errors.New("некорректный литерал")
	}
	if lit[0] != '"' || lit[len(lit)-1] != '"' {
		return "", errors.New("неподдерживаемый литерал")
	}
	return lit[1 : len(lit)-1], nil
}

func newFinding(fset *token.FileSet, path string, pos token.Pos, message string) finding {
	position := fset.Position(pos)
	return finding{
		path:    path,
		line:    position.Line,
		message: message,
	}
}
