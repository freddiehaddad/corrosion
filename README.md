# Corrosion

A compiler and interpreter project written in Go and inspired by the books
[Writing A Compiler In Go] and [Writing An Interpreter In Go] by Thorsten Ball.

## Language Style

Corrosion inherits a few language syntax styles. Some examples:

```C
// Syntax examples
var foo = 100;
func add(left, right) { return left + right; }
add(foo, foo); // 200

func conditional(check) {
    if (!check == true) {
        return true;
    } else {
        return false;
    }
}

conditional(); // false

func foo() {
    func bar() { return 2; }
    return bar;
}

foo()(); // 2
```

## Obtaining Source

```bash
git clone https://github.com/freddiehaddad/corrosion
```

## Building

```bash
go build -o bin ./...
```

## Testing

```bash
go test -v ./...
```

## Running

After building the code, the REPL can be launched with:

```bash
./bin/corrosion
```

## Dependencies

Go (see [go.mod] for minimum version) is required for building. In general, any
recent version should work.

## Project Layout

```text
.
├── bin
│   └── corrosion
├── cmd
│   └── corrosion
│       └── corrosion.go
├── go.mod
├── LICENSE
├── pkg
│   ├── ast
│   │   └── ast.go
│   ├── evaluator
│   │   ├── evaluator.go
│   │   └── evaluator_test.go
│   ├── lexer
│   │   ├── lexer.go
│   │   └── lexer_test.go
│   ├── object
│   │   └── object.go
│   ├── parser
│   │   ├── parser.go
│   │   └── parser_test.go
│   └── token
│       └── token.go
└── README.md
```

## License

Licensed under the [MIT] license.

[go.mod]: go.mod
[mit]: LICENSE
[writing a compiler in go]: https://compilerbook.com/
[writing an interpreter in go]: https://interpreterbook.com/
