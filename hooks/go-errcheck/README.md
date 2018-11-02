# errcheck

errcheck is a program for checking for unchecked errors in go programs.

## Use

For basic usage, just give the package path of interest as the first argument.

To check all packages beneath the current directory:

    errcheck ./...

Or check all packages in your $GOPATH and $GOROOT:

    errcheck all

errcheck also recognizes the following command-line options:

The `-tags` flag takes a space-separated list of build tags, just like `go
build`. If you are using any custom build tags in your code base, you may need
to specify the relevant tags here.

The `-asserts` flag enables checking for ignored type assertion results. It
takes no arguments.

The `-blank` flag enables checking for assignments of errors to the
blank identifier. It takes no arguments.


## Excluding functions

Use the `-exclude` flag to specify a path to a file containing a list of functions to
be excluded.

    errcheck -exclude errcheck_excludes.txt path/to/package

The file should contain one function signature per line. The format for function signatures is
`package.FunctionName` while for methods it's `(package.Receiver).MethodName` for value receivers
and `(*package.Receiver).MethodName` for pointer receivers.

An example of an exclude file is:

    io/ioutil.ReadFile
    (*net/http.Client).Do

The exclude list is combined with an internal list for functions in the Go standard library that
have an error return type but are documented to never return an error.

## Cgo

Currently errcheck is unable to check packages that import "C" due to limitations in the importer when used with versions earlier than Go 1.11.

However, you can use errcheck on packages that depend on those which use cgo. In order for this to work you need to go install the cgo dependencies before running errcheck on the dependent packages.

See https://github.com/kisielk/errcheck/issues/16 for more details.

## Exit Codes

errcheck returns 1 if any problems were found in the checked files.
It returns 2 if there were any other failures.