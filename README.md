[![GoDoc](https://godoc.org/github.com/fnproject/fdk-go?status.svg)](https://godoc.org/github.com/fnproject/fdk-go)

# Go FDK Documentation
This is documentation for the Go function development kit (FDK). The Go FDK  provides convenience functions for writing Go Fn code.

## User Information
* See the Fn [Quickstart](https://github.com/fnproject/fn/blob/master/README.md) for sample commands.
* [Detailed installation instructions](http://fnproject.io/tutorials/install/).
* [Configure your CLI Context](http://fnproject.io/tutorials/install/#ConfigureyourContext).
* For a list of commands see [Fn CLI Command Guide and Reference](https://github.com/fnproject/docs/blob/master/cli/README.md).
* For general information see Fn [docs](https://github.com/fnproject/docs) and [tutorials](https://fnproject.io/tutorials/).

## Go FDK Development
See [CONTRIBUTING](https://github.com/fnproject/fn/blob/master/CONTRIBUTING.md) for information on contributing to the project.

### Notes
If you poke around in the Dockerfile you'll see that we simply add the `.go` source file and the `fdk-go` package to our workspace, then build a binary.  We then build an image with that binary that gets deployed to dockerhub and Fn.

For more robust projects, it's recommended to use a tool like `dep` or `glide` to get dependencies such as the `fdk-go` into your functions.


