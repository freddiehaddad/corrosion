# Corrosion

A compiler and interpreter project written in Go and inspired by the books
[Writing A Compiler In Go] and [Writing An Interpreter In Go] by Thorsten Ball.

## Building

Go (see [go.mod] for minimum version) is required for building. In general, any
recent version should work.

```bash
git clone https://github.com/freddiehaddad/corrosion
cd corrosion
go build -o bin ./...
```

## Project Layout

```text
├── bin
│   └── corrosion
├── cmd
│   └── corrosion
│       └── corrosion.go
├── go.mod
├── LICENSE
├── pkg
│   ├── lexer
│   │   ├── lexer.go
│   │   └── lexer_test.go
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
