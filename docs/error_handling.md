# Error Handling

## How to handle errors

### Returning Errors

In most cases, if a function encounters an error, it does not have the necessary context to properly handle the error by itself, so it has to pass the error back to its caller.

As an example, see `func ReadFile()` from the sample code (`readfile.go`):

```go
func ReadFile(path string) ([]byte, error) {
    if path == "" {
       // Create an error with errors.New()
       return nil, errors.New("path is empty")
    }
    f, err := os.Open(path)
    if err != nil {
       // Wrap the error.
       // If the format string uses %w to format the error,
       // fmt.Errorf() returns an error that has the
       // method "func Unwrap() error" implemented.
       return nil, fmt.Errorf("open failed: %w", err)
    }
    defer f.Close()

    buf, err := io.ReadAll(f)
    if err != nil {
       return nil, fmt.Errorf("read failed: %w", err)
    }
    return buf, nil
}
```

`ReadFile()` tests the received path, and if the path is empty, it creates a new error and returns it. The data that `ReadFile()` was supposed to return does not exist; therefore, `ReadFile()` returns a `nil` value:

```go
if path == "" {
    return nil, errors.New("path is empty")
}
```

Conventionally, if a function returns an error value, it's always the last (rightmost) value in the list of return values:

```go
func ReadFile(path string) ([]byte, error) {
```

The caller of `ReadFile()` receives the error value along with the result value. By default, the returned error value is assigned to a variable named `err`:

```go
_, err := ReadFile("no/file")
if err != nil {
    fmt.Println("Error:", err)
}
```

### Panic and Recover

Go newcomers might miss the `try...catch` mechanism that other languages provide. However, Go has something that fulfills a similar purpose: `panic` and `recover`. But beware! Unlike `try...catch`, `panic` and `recover` is not, and should not be, the standard way of handling errors. Panicking is only acceptable if an error is indeed unexpected and there is no way of handling it. In such cases, it's better to have the application crash early and restart it.

In certain cases, crashing the app might not be an option. Consider an HTTP server that must be up and running without disruption. If a panic occurs when handling a request, all other requests should continue being handled, if possible. To do this, the `net/http` package uses Go's recovery technique.

### Logging Errors

If a function can handle an error it receives from a called function, it might want to write information about the error to a log file.

> [!NOTE]
> If you write code for a library, consider not logging anything. The library clients will have different opinions about which logger to use and what is printed to stdout or stderr. So, it is almost always better to only return errors and let the library clients do the logging they want.

### Using Error Wrapping

An error often "bubbles up" a call chain of multiple functions. In other words, a function receives an error and passes it back to its caller through a return value. The caller might do the same, and so on, until a function up the call chain handles or logs the error. Each function involved in this "bubbling up" can add valuable contextual information to the error before handing it back to its caller.

A function should only pass the error on unchanged if it cannot add any helpful information:

```go
if err != nil {
    // Only do that if no additional context can be added!
    return err
}
```

In all other cases, it should add appropriate contextual information. However, simply concatenating a new error message with the original one does not work:

```go
// WRONG!
if err != nil {
    return errors.New("open failed:" + err.Error())
}
```

This would only preserve the original error message, while all other information gets lost.

Instead, you should use error wrapping. An error can be "wrapped" around another error using `fmt.Errorf()` and the special formatting verb `%w`. See the `ReadFile()` function in the file `readfile.go`:

```go
f, err := os.Open(path)
if err != nil {
    return nil, fmt.Errorf("open failed: %w", err)
}
```

`os.Open()` returns an error type that contains additional information, as you will see later. Wrapping the error preserves all this additional information.

### Unwrapping Wrapped Errors

An error returned by a function might contain one or more wrapped errors. Printing or logging the received error will also include all error messages from the wrapped errors. However, sometimes you need to know if a particular type of error is nested somewhere inside the layers of errors.

For example, let's see how to handle `ReadFile()`'s errors in `func main()`:

```go
_, err := ReadFile("no/file")
log.Println("err = ", err)

// Unwrap the error returned by os.Open()
log.Println("errors.Unwrap(err) = ", errors.Unwrap(err))
```

This code snippet prints the following:

```go
Reading a single file: err =  open failed: open no/file: no such file or directory
Reading a single file: errors.Unwrap(err) =  open no/file: no such file or directory
```

While the wrapped error message is `open failed: open no/file: no such file or directory`, the unwrapped error contains only `open no/file: no such file or directory`, excluding the open failed: message that was added to the wrapped error.

This way, you can unwrap one error after another until you hit the end of the chain.

### Testing for Specific Error Types

Occasionally, you need to know if any of the errors inside a chain of wrapped errors are of a particular type.

For example, `os.Open` returns an error of type `fs.PathError` that not only records the error but also the operation and the path that caused it. If you can find out that the error chain contains this error, you can make use of the additional information for troubleshooting.

To achieve this, the errors package provides two functions: `Is()` and `As()`.

#### errors.Is()

Function `func Is(err, target error) bool` returns `true` if error `err` is of the same type as `target`.

In the case of the `ReadFile()` function, you can verify that the returned error is, or wraps, an `fs.ErrNotExist` error:

```go
_, err := ReadFile("no/file")

log.Println("err is fs.ErrNotExist:", errors.Is(err, fs.ErrNotExist))
```

#### errors.As()

You'll also want to access the path information. For this, you not only need to ensure the error wraps an `fs.PathError` but also access this `PathError` and all its methods.

To do this, use the function `func As(err error, target any) bool`. Like `Is()`, function `As()` returns `true` if `err` is or wraps an error of the same type as `target`, and it also unwraps that error and assigns it to `target`.

This requires defining a variable of type `fs.PathError` and passing a pointer to that variable to `As()`:

```go
target := &fs.PathError{}
if errors.As(err, &target) {
    log.Printf("err as PathError: path is '%s'\n", target.Path)
	log.Printf("err as PathError: op is '%s'\n", target.Op)
}
```

### Joining Errors

Typically, errors get wrapped one by one while being returned to the respective caller. Sometimes, a function needs to collect multiple errors and wrap them into one.

Take the function `ReadFiles()` (note the plural) as an example. This function reads multiple files and returns all file contents that were successfully read. If one or more files fail to be read, `ReadFiles()` will collect the errors and join them into one.

For this, the errors package provides the `Join()` function (since Go 1.20). Let's see how `ReadFiles()` makes use of the `Join()` function:

```go
func ReadFiles(paths []string) ([][]byte, error) {
    var errs error
    var contents [][]byte

    if len(paths) == 0 {
       // Create a new error with fmt.Errorf() (but without using %w):
       return nil, fmt.Errorf("no paths provided: paths slice is %v", paths)
    }

    for _, path := range paths {
       content, err := ReadFile(path)
       if err != nil {
        errs = errors.Join(errs, fmt.Errorf("reading %s failed: %w", path, err))
          continue
       }
       contents = append(contents, content)
    }

    return contents, errs
}
```

If an error occurs inside the `for` loop, it does not break the loop. Instead, it is joined to variable `errs`, and the loop continues, joining more records as they occur.

Finally, `ReadFiles()` returns both the contents read successfully and the joined error messages.

### Handling Joined Errors

Now, you might expect that joined errors can be unwrapped like single errors. Unfortunately, this is not the case. A joined error is actually a slice of errors, `[]error`. The `Unwrap()` function, however, returns a single `error`. If called on a joined error, `Unwrap()` returns `nil`:

```go
_, err = ReadFiles([]string{"no/file/a", "no/file/b", "no/file/c"})
log.Println("joined errors = ", err)

log.Println("errors.Unwrap(err) = ", errors.Unwrap(err))
```

The second log line prints:

```go
errors.Unwrap(err) =  <nil>
```

Fortunately, there is a way to unwrap the slice of joined errors. The joined error type itself helps you do this by providing an `Unwrap() []error` method that returns the error slice.

To access this `Unwrap()` method, you only need to type-assert that the error variable implements this method. You can then call it safely:

```go
e, ok := err.(interface{ Unwrap() []error })
if ok {
    log.Println("e.Unwrap() = ", e.Unwrap())
}
```

This prints the full set of joined errors:

```
Reading multiple files: e.Unwrap() =  [reading no/file/a failed: open failed: open no/file/a: no such file or directory
reading no/file/b failed: open failed: open no/file/b: no such file or directory reading no/file/c failed: open failed: open no/file/c: no such file or directory]
```

### Context-Based Error Handling

The `context` package is popular for controlling timeouts of requests or canceling multiple goroutines upon request. If you use a cancelable context, you can inspect and handle the error that caused the cancellation.

Since Go 1.20, you can even send a custom error message when canceling a context by using a `WithCancelCause` context. The following is a basic example:

```go
parent := context.Background()
ctx, cancel := context.WithCancelCause(parent)
defer cancel(nil)             // Set the cause to Canceled
cancel(fmt.Errorf("myError")) // Set the cause to myError

fmt.Println(ctx.Err())          // Output: context.Canceled
fmt.Println(context.Cause(ctx)) // Output: myError
```

The context function `WithCancelCause()` returns a context and a cancel function that expects an error type. When calling `cancel`, a custom error message can be passed as input. All interested parties that have access to the context can retrieve the custom error through `context.Cause(ctx)`.

## Best Practices

### Use the `defer` Function

A function can exit at multiple points, through return statements as well as panics. Whenever a function allocates resources, such as files, network connections, or goroutines, use a defer() function to clean up any open resources at function exit.

### Provide Explicit Error Information:

Nothing is more frustrating than seeing some cryptic error message like `ERROR: EPIC FAIL` in the log files without any clue about the context in which the error occurred.

Therefore, if a function encounters an error, it should not pass the error verbatim up the call chain. Rather, if any contextual information is available to help troubleshoot the error, this information should be added to the error by wrapping it in a new error. (See the earlier section on using error wrapping.)

### Use `Panic` and `Recover` Only When Necessary

Go newcomers often frown upon Go's verbose error handling and want to save typing by letting a function panic instead of handling an error. At the top level, the panic is then recovered and handled. This approach, however, is unidiomatic Go and has downsides. First and foremost, adding useful contextual information (see the previous section) is not possible with this method. Moreover, because a panic unwinds the call stack outside the regular call/return flow, any function in the call chain between the top-level function and the panicking function contains no error handling code. How can a reader see that any of these functions might observe an error? For comparison, Java has the `throws` keyword to list all exceptions a function may emit. Go does not have such a feature. It's not possible to see if any of the callees of a function panics. Standard Go error handling makes the error flow clearly visible.

Go treats errors as a normal part of the program flow because they are exactly that. If an error occurs, it should be handled or passed to the caller until some function up the call chain handles the error or writes it to a log file for troubleshooting.

If you inspect a function, you'll want to immediately see which errors it may encounter and how it passes them up the call chain.

There are also some categories of errors that cannot be handled at all, such as an out-of-memory situation. If the required memory cannot be allocated, the application has no meaningful way to continue and should panic.

On the other hand, user input at runtime is expected to be unreliable. Any error resulting from user input, invalid or missing files, a network timeout, or other predictable sources of failure can and should be handled as an error.

### Use Libraries and Packages That Follow Error Handling Best Practices

If you have a choice between multiple third-party packages that deliver identical or similar functionality, choose the one that follows best practices for error handling.

You will not do yourself any favors if you decide to use the package with the fanciest API but with brittle error handling. Any package that suppresses errors rather than properly passing them back—or that provides no context for errors—will turn troubleshooting into a hit-or-miss debugging nightmare.

So, take a peek at the code inside a package to see if it contains robust code with proper error handling. This precautionary measure will pay off in the long run.

### Create Custom Error Types Wherever Suitable

Because `error` is an interface, you can build custom error types with extra functionality as long as they implement `Error() string`. You saw an example in the "Testing for Specific Error Types" section, where `os.Open` returned an `fs.PathError`.

This error is a `struct` that implements the methods `Error()`, `Unwrap()`, and `Timeout()` and provides the fields `Path`, `Op`, and `Error` to capture detailed error information:

```go
type PathError struct {
  Op   string
  Path string
  Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }

func (e *PathError) Unwrap() error { return e.Err }

// Timeout reports whether this error represents a timeout.func (e *PathError) Timeout() bool {
  t, ok := e.Err.(interface{ Timeout() bool })
  return ok && t.Timeout()
}
```

In the same manner, you can create your own error types. The only mandatory method to implement is `Error()`, but if you also implement the method `Unwrap()`, then the package function `errors.Unwrap()` will be able to unwrap your error.

## Handling Specific Types of Errors

Some types of errors require special treatment due to their specific nature. These types include network errors, I/O errors, and system errors.

### Network Errors

Failing network connections need special treatment. A network error can be caused by a permanent failure or by a temporary issue. Code that handles a network error needs to distinguish between these two situations.

Consider the task of opening a new TCP connection. This task can fail because the network is temporarily down or because the system at the other end of the connection is restarting or overloaded and cannot accept new connections at the moment.

In such cases, you'll want to try connecting again at a later time. The `net.Dial()` function, for example, supports this by returning a specific error type, `net.OpError`, that provides a method named `Temporary()` for testing if the error is expected to eventually go away.

With the `Temporary()` method, you can implement a simple retry algorithm like the one below or a more sophisticated strategy like [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff):

```go
func connectToTCPServer() error {
    var err error
    var conn net.Conn
    for retry := 3; retry > 0; retry-- {
       conn, err = net.Dial("tcp", "127.0.0.1:12345")
       if err != nil {
          // Check if err is a net.OpError
          opErr := &net.OpError{}
          if errors.As(err, &opErr) {
             log.Println("err is net.OpError:", opErr.Error())
             // test if the error is temporary
             if opErr.Temporary() {
                log.Printf("Retrying...\n")
                continue
             }
             retry = 0
          }
       }
    }
    if err != nil {
       return fmt.Errorf("connect failed: %w", err)
    }
    defer conn.Close()
    // send or receive data
    return nil
}
```

### I/O Errors

Recovering from an I/O error that occurs after having read or written large amounts of data can be costly. all the data that's already been processed up to the point where the error occurs might need to be read or written again.

To allow for a more efficient recovery, most I/O-related functions and methods in the standard library return not only an error but also the number of bytes that were successfully processed. A typical example is `io.Reader`'s `Read()` function:

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

An error recovery procedure could use this information to continue the I/O operation where it was interrupted.

> [!NOTE]
> The `io` package provides the sentinel error value `io.EOF` (that is defined as `errors.New("EOF")`) to signal the successful(!) end of reading an input stream. Every type that implements the `io.Reader` interface should stick to the documented semantics of returning an error:
>
> ```md
> …a Reader returning a non-zero number of bytes at the end of the input stream may return either `err == EOF` or `err == nil`. The next Read should return `0`, `EOF`.
> ```

## Common Mistakes to Avoid When Handling Errors in Go

While Go's error handling may seem unusual at first sight, it's logical and straightforward to use. However, this doesn't mean that you can't make errors with error handling. Here are some mistakes to avoid.

### Ignoring Errors

The biggest mistake a developer can make in any programming language is to ignore errors. Not catching errors early easily leads to follow-up errors that can be much more difficult to track down compared to the original error if it had been properly handled.

So, the number one rule for avoiding error handling mistakes is to never assign a returned error value to the blank identifier.

```
Fun fact: did you know that fmt.Println() returns an error value?
```

Bottom line—don't do this:

```go
WriteString(w, s)
```

Do this instead:

```go
n, err := WriteString(w, s)
// error handling here
```

### Not Wrapping Errors with Additional Context When Propagating

Often, if not always, a function that receives an error from calling another function can add valuable contextual information to the error.

So, whenever you find yourself writing this:

```go
n, err := WriteString(w, s)
if err != nil {
    err
}
```

Take a step back and see if you can include contextual information. In most cases, you can. Even the function name can be valuable information because it allows you to track the chain of function calls that lead to the error:

```go
n, err := WriteString(w, s)
if err != nil {
    return fmt.Errorf("after writing %d characters: %w", n, err)
}
```

It's a few more strokes on the keyboard for you, but it can be an enormous time-saver later on.

### Overgeneralizing Errors

When composing error messages, be as specific as you can. Include all the contextual information you have.

An error message like "database error" can have a truckload of different possible causes. The message "database error" is genuinely pointless and unhelpful.

Add as much information to the error message as you can. Consider creating custom error types that can carry additional information.

### Using Incorrect Error Types

The particular type of error value might seem like a negligible detail. After all, every error implements `type error interface{ Error() string }`, so in the end, errors are nothing but glorified string types, right?

Wrong. Custom error types can contain extra information and enable advanced error inspection through `errors.Is()` and `errors.As()`.

So, whenever you send an error back to a caller, make sure to use the error type that is appropriate for the given error context.

### Not Logging Errors

Error messages are indispensable for troubleshooting. Whether an app can handle an error or whether an error forces the app to terminate, the app should log that error for postmortem analysis.

In general, if a function observes an error, it should either handle the error or return it to its caller.

If it can handle the error or if it cannot return the error for some reason (maybe because it is function `main()`), the function should always log the error and all its contextual information.

Every error that occurs indicates an opportunity for fixing a bug or improving the code. Don't let this opportunity pass by unnoticed.

### Logging Errors with log.Fatal()

If your application encounters an unrecoverable error, it might feel natural to log this error by calling `log.Fatal()`, which conveniently logs a message and exits the process immediately.

However, there is a catch. `log.Fatal()` calls `os.Exit()`. Unlike a call to `panic()`, `os.Exit()` is not recoverable and skips all deferred functions.

A good practice is to write `func main()` so that it does not defer any functions and call `log.Fatal()` or `os.Exit()` exclusively in `main()`.

### Not Considering Error Recovery

"Crash early" is good advice in many circumstances. Crashing an app allows it to restart from a clean state. However, crashing is not always the best option.

- If an error is easy to recover from, crashing the whole application is an overreaction.
- If a process guarantees maximum uptime, it's better to do your best to recover from the error rather than disrupting the system with a restart.
- If a process spawns goroutines, it's often sufficient to exit a single goroutine that observes an error condition. `http.ListenAndServe()` is an example of this strategy. All incoming requests are handled in separate goroutines, and if one goroutine panics, `ListenAndServe()` recovers from that panic so that all other concurrent handlers can continue unaffected.

Bottom line: applications may benefit from well-designed error recovery, especially if crashing early entails the considerable cost of respawning the app.
