# mango

*Generate manual pages from the source code of your Go commands*

mango is a small command line utility that allows you to create manual
pages from the source code of your Go commands. It builds manual pages from
the comments and *flag* function calls found in your .go files.

#### Building

Execute
```bash
go build ./mango.go
```
to build mango.

#### Running

Pass one or more .go files as command line arguments to mango.
mango treats them as a list of independent Go commands and creates a
manual page for each argument.

```bash
mango file1.go file2.go ...
```

## Usage

```go
// TODO:...
```

### License

mango is released under MIT license.
You can find a copy of the MIT License in the [LICENSE](./LICENSE) file.
