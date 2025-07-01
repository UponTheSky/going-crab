# Basic Log
Over the next three chapters, we will cover topics that are necessary for the robustness of our server: Logging, Error Handling, and Testing. This chapter is about Logging.

You can assume that a log is a text message telling you what is going on inside the server. It could be a simple message like a user has made a GET request, or some serious events that may directly affect the user experience, such as one of your DB has been going down due to resource exhaustion. 

You already have seen a few lines of logs in the code examples, such as `log.Fatal(err)`. In this chapter, we'll mainly deal with this [`log`](https://pkg.go.dev/log) package briefly. In a later chapter, we will have time to talk about more sophisticated techniques such as [`log/slog`](https://pkg.go.dev/log/slog) package.

Since this chapter will be very short, I wrote a simple Go code in `main.go` that write a short message to standard error. Check by running `go run .` as usual.

## Configure Log Messages
The first thing we need to do is to configure a set of prefixes showing information about the log message we will going to make, by generating a [`Logger`](https://pkg.go.dev/log#New) object:

```go
// [...]
logger := log.New(
    os.Stderr,
    "[Going Crab] ",
    log.LstdFlags|log.Lshortfile|log.LUTC,
)
// [...]
```
Here we create a logger object, such that the log messages created by this logger are flushed to the standard error, with prefix `[Going Crab]: `, and finally a set of configs additional to our own prefix.

An example of the log is like this: `[Going Crab] 2025/07/01 14:59:41 main.go:15: 42, test!`. Here we have

- Our own prefix "Going Crab"
- the time when this log was generated
- the line of our code where the log occurred
- Our message, `42, test!`

## Make Log Messages
In package `log`, there are in total nine methods for making log messages. But they can be classified by these two categories:

1. Format
- print as is: `Print`, `Panic`, `Fatal`
- print as is, but append a newline(`\n`): `Println`, `Panicln`, `Fatalln`
- print a formatted string: `Printf`, `Panicf`, `Fatalf`

There are more details for each of these categories but I don't think you have to know that further. Sometimes it is important to determine what is important to know or not, since your time is valuable and you should use it efficiently.

2. Behavior
- just print: `Print`, `Println`, `Printf`
- print and panic: `Panic`, `Panicln`, `Panicf`
- print and exit the current running application: `Fatal`, `Fatalln`, `Fatalf`

By far we have often used `Fatal` or its friends, but using `Fatal` bloated everywhere is definitely not good for your application. Except some serious logical errors or security issues, you won't want to shut down the entire application just for a small user input error like the 400 error. So in a few next chapters we will mostly use `log.Printf` or `log.Println`, by the time we finally cover the package `log/slog`.

About panic, I recommend this [Go Dev blog post](https://go.dev/blog/defer-panic-and-recover). For simplicity of this book, I won't use `log.Panic` and its family as well as `log.Fatal`. 

## Conclusion
This chapter has briefly skimmed through how to write a simple log using `log` package. The next chatper is about error handling in a running HTTP server.

## Exercise
This chapter is too short to do a fun exercise... I am so sorry but let us bring this joy to the next chapter!
