# Function Examples

The goal of the `fdk`'s are to make functions easy to write.

This is an example of a function using the fdk-go bindings. The [function
documentation](https://github.com/fnproject/fn/blob/master/docs/developers/fn-format.md)
contains details of how this example works under the hood. With any of the
examples provided here, you may use any format to configure your functions in
`fn` itself.

### How to run the example

Install the CLI tool, start a Fn server and run `docker login` to login to
DockerHub. See the [front page](https://github.com/fnproject/fn) for
instructions.

Initialize the example with an image name you can access:

```sh
fn init --runtime docker --name hello
```

Build and deploy the function to the Fn server (default localhost:8080)

```sh
fn deploy --app myapp
```

Now call your function (may take a sec to pull image):

```sh
echo '{"name":"Clarice"}' | fn invoke myapp hello
```

**Note** that this expects you were in a directory named 'hello' (where this
example lives), if this doesn't work, replace 'hello' with your `$PWD` from
the `deploy` command.

### Details

If you poke around in the Dockerfile you'll see that we're simply adding the
file found in this directory, getting the `fdk-go` package with `go mod`
and then building a binary and building an image with that binary. That then
gets deployed to dockerhub and fn.

Scoping out `func.go` you can see that the handler code only deals with input
and output, and doesn't have to deal with decoding the formatting from
functions (i.e. i/o is presented through `io.Writer` and `io.Reader`). This
makes it much easier to write functions.
