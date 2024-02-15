// # Пакет кастмоного линтера для приложения сервиса сокращения ссылок.
//
// Включает в себя следующие анализаторы для вызова проверок:
//
//   - Все анализаторы пакета analysis/passes
//
//   - Два публичных анализатора "github.com/fatih/errwrap/errwrap" и "github.com/gordonklaus/ineffassign/pkg/ineffassign"
//
//   - Кастомный анализатор анализатора для проверки вызова метода os.Exit из функции main приложения
//
// # Для вызова аналзатор, необходимо собрать пакет staticlint (cd cmd/staticlint/)
//
//	cd cmd/staticlint/ && go build
//
// # Вызвать проверку анализатором командой
//
//	./staticlint ../../...
//
// # Подключенные и используемые для проверки анализаторы пакета analysis/passes
//
//   - appends - Package appends defines an Analyzer that detects if there is only one variable in append.
//
//   - asmdecl - Package asmdecl defines an Analyzer that reports mismatches between assembly files and Go declarations.
//
//   - assign - Package assign defines an Analyzer that detects useless assignments.
//
//   - atomic -	Package atomic defines an Analyzer that checks for common mistakes using the sync/atomic package.
//
//   - atomicalign - Package atomicalign defines an Analyzer that checks for non-64-bit-aligned arguments to sync/atomic functions.
//
//   - bools - Package bools defines an Analyzer that detects common mistakes involving boolean operators.
//
//   - buildssa - Package buildssa defines an Analyzer that constructs the SSA representation of an error-free package and returns the set of all functions within it.
//
//   - buildtag - Package buildtag defines an Analyzer that checks build tags.
//
//   - cgocall - Package cgocall defines an Analyzer that detects some violations of the cgo pointer passing rules.
//
//   - composite - Package composite defines an Analyzer that checks for unkeyed composite literals.
//
//   - copylock - Package copylock defines an Analyzer that checks for locks erroneously passed by value.
//
//   - ctrlflow - Package ctrlflow is an analysis that provides a syntactic control-flow graph (CFG) for the body of a function.
//
//   - deepequalerrors - Package deepequalerrors defines an Analyzer that checks for the use of reflect.DeepEqual with error values.
//
//   - defers -	Package defers defines an Analyze
//
//   - directive - Package directive defines an Analyzer that checks known Go toolchain directives.
//
//   - errorsas - The errorsas package defines an Analyzer that checks that the second argument to errors.As is a pointer to a type implementing error.
//
//   - fieldalignment - Package fieldalignment defines an Analyzer that detects structs that would use less memory if their fields were sorted.
//
//   - findcall - Package findcall defines an Analyzer that serves as a trivial example and test of the Analysis API.
//
//   - framepointer - Package framepointer defines an Analyzer that reports assembly code that clobbers the frame pointer before saving it.
//
//   - httpmux - The httpmux command runs the httpmux analyzer.
//
//   - httpresponse - Package httpresponse defines an Analyzer that checks for mistakes using HTTP responses.
//
//   - ifaceassert - Package ifaceassert defines an Analyzer that flags impossible interface-interface type assertions.
//
//   - inspect - Package inspect defines an Analyzer that provides an AST inspector (golang.org/x/tools/go/ast/inspector.Inspector) for the syntax trees of a package.
//
//   - loopclosure - Package loopclosure defines an Analyzer that checks for references to enclosing loop variables from within nested functions.
//
//   - lostcancel - Package lostcancel defines an Analyzer that checks for failure to call a context cancellation function.
//
//   - nilfunc - Package nilfunc defines an Analyzer that checks for useless comparisons against nil.
//
//   - nilness - Package nilness inspects the control-flow graph of an SSA function and reports errors such as nil pointer dereferences and degenerate nil pointer comparisons.
//
//   - pkgfact - The pkgfact package is a demonstration and test of the package fact mechanism.
//
//   - printf - Package printf defines an Analyzer that checks consistency of Printf format strings and arguments.
//
//   - reflectvaluecompare - Package reflectvaluecompare defines an Analyzer that checks for accidentally using == or reflect.DeepEqual to compare reflect.Value values.
//
//   - shadow - Package shadow defines an Analyzer that checks for shadowed variables.
//
//   - shift - Package shift defines an Analyzer that checks for shifts that exceed the width of an integer.
//
//   - sigchanyzer - Package sigchanyzer defines an Analyzer that detects misuse of unbuffered signal as argument to signal.Notify.
//
//   - slog - Package slog defines an Analyzer that checks for mismatched key-value pairs in log/slog calls.
//
//   - sortslice - Package sortslice defines an Analyzer that checks for calls to sort.Slice that do not use a slice type as first argument.
//
//   - stdmethods - Package stdmethods defines an Analyzer that checks for misspellings in the signatures of methods similar to well-known interfaces.
//
//   - stringintconv - Package stringintconv defines an Analyzer that flags type conversions from integers to strings.
//
//   - structtag - Package structtag defines an Analyzer that checks struct field tags are well formed.
//
//   - testinggoroutine - Package testinggoroutine defines an Analyzerfor detecting calls to Fatal from a test goroutine.
//
//   - tests - Package tests defines an Analyzer that checks for common mistaken usages of tests and examples.
//
//   - timeformat - Package timeformat defines an Analyzer that checks for the use of time.Format or time.Parse calls with a bad format.
//
//   - unmarshal - The unmarshal package defines an Analyzer that checks for passing non-pointer or non-interface types to unmarshal and decode functions.
//
//   - unreachable -	Package unreachable defines an Analyzer that checks for unreachable code.
//
//   - unsafeptr - Package unsafeptr defines an Analyzer that checks for invalid conversions of uintptr to unsafe.Pointer.
//
//   - unusedresult - Package unusedresult defines an analyzer that checks for unused results of calls to certain pure functions.
//
//   - unusedwrite -	Package unusedwrite checks for unused writes to the elements of a struct or array object.
//
//   - usesgenerics - Package usesgenerics defines an Analyzer that checks for usage of generic features added in Go 1.18.
package main
