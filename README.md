# mango

*Generate manual pages from the source code of your Go commands*

*mango* is a small command line utility that allows you to create manual
pages from the source code of your Go commands.

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

### License

mango is released under MIT license.
You can find a copy of the MIT License in the [LICENSE](./LICENSE) file.

