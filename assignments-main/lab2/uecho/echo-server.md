# UDP Echo Server

In this task we will focus on the user datagram protocol (UDP), which provides unreliable datagram service.
You will find the documentation of the [UDPConn](https://golang.org/pkg/net/#UDPConn) type useful.

In the provided code under `uecho`, we have implemented a simple `SendCommand()` function that acts as a client, along with a bunch of tests.
You can run these test with `go test -v`, and as described in Lab 1, you can use the `-run` flag to run only a specific test.

You can also compile your server code into a binary using `go build`.
This will produce a file named `uecho` in the same folder as the `.go` source files.
You can run this binary in two ways:

1. `./uecho -server &` will start the server in the background.
   Note: _This will not work until you have implemented the necessary server parts._

2. `./uecho` will start the command line client, from which you may interact with the server by typing commands into the terminal window.

If you want to extend the capabilities of this runnable client and server, you can edit the files `echo.go` and `echo_client.go`.
But note that the tests executed by the quickfeed will use original `SendCommand()` provided in the original `echo_client.go` file.
If you've done something fancy, and want to show us that's fine, but it won't be considered by the quickfeed.

## Echo Server Specification

The `SendCommand()` takes the following arguments:

| Argument  | Description                                                                  |
| --------- | ---------------------------------------------------------------------------- |
| `udpAddr` | UDP address of the server (`localhost:12110`)                                |
| `cmd`     | Command (as a text string) that the server should interpret and execute      |
| `txt`     | Text string on which the server should perform the command provided in `cmd` |

The `SendCommand()` function produces a string composed of the following

```text
cmd|:|txt
```

For example:

```text
UPPER|:|i want to be upper case
```

From this, the server is expected to produce the following reply:

```text
I WANT TO BE UPPER CASE
```

See below for more details about the specific behaviors of the server.

1. For each of the following commands, implement the corresponding functions, so that the returned value corresponds to the expected test outcome.
   Here you are expected to implement demultiplexer that demultiplexes the input (the command) so that different actions can be taken.
   A hint is to use the `switch` statement. You will probably also need the `strings.Split()` function.

   | Command | Action                                                                                                                                                                                                                           |
   | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
   | UPPER   | Takes the provided input string `txt` and applies the translates it to upper case using `strings.ToUpper()`.                                                                                                                     |
   | LOWER   | Same as UPPER, but lower case instead.                                                                                                                                                                                           |
   | CAMEL   | Same as UPPER, but title or camel case instead.                                                                                                                                                                                  |
   | ROT13   | Takes the provided input string `txt` and applies the rot13 translation to it; see lab1 for an example.                                                                                                                          |
   | SWAP    | Takes the provided input string `txt` and inverts the case. For this command you will find the `strings.Map()` function useful, together with the `unicode.IsUpper()` and `unicode.ToLower()` and a few other similar functions. |

2. The server should reply `Unknown command` if it receives an unknown command or fails to interpret a request in any way.

3. Make sure that your server continues to function even if one client's connection or datagram packet caused an error.

### Echo Server Implementation

You should implement the specification by extending the skeleton code found in `echo_server.go`:
