package osexitlinter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// OSExitAnalyzer - инициализация переменной кастомного анализатора.
var OSExitAnalyzer = &analysis.Analyzer{
	Name:     "osexit",
	Doc:      "Check if os.Exit method exists in main function at main package",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	isMainFile := func(x *ast.File) bool {
		return x.Name.Name == "main"
	}

	isMainFunc := func(x *ast.FuncDecl) bool {
		return x.Name.Name == "main"
	}

	isOsExit := func(x *ast.SelectorExpr, isMain bool) bool {
		if !isMain || x.X == nil {
			return false
		}

		ident, ok := x.X.(*ast.Ident)
		if !ok {
			return false
		}

		if ident.Name == "os" && x.Sel.Name == "Exit" {
			pass.Reportf(ident.NamePos, "os.Exit called in main func in main package")
			return true
		}

		return false
	}

	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.SelectorExpr)(nil),
	}

	mainInspecting := false
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch x := n.(type) {
		case *ast.File:
			if !isMainFile(x) {
				return
			}
		case *ast.FuncDecl:
			f := isMainFunc(x)
			if mainInspecting && !f {
				mainInspecting = false
				return
			}
			mainInspecting = f
		case *ast.SelectorExpr:
			if isOsExit(x, mainInspecting) {
				return
			}
		}
	})

	return nil, nil
}
