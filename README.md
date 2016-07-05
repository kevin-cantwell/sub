A simple tool that will perform pattern substitutions using Go's regex parser.

### Install

```
go get -u github.com/kevin-cantwell/sub/cmd/sub
```

### Usage

You may use parameters in the replacement argument to specify submatches:

```
echo "foobarbaz" | sub '^foo(bar|biz)baz$' 'bar or biz: $1'
```
