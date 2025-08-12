# Advanced Log
If you recall the [basic log chapter](../../basic/log/README.md), we used [log](https://pkg.go.dev/log) package to write simple log messages. `log` package only provides functions that write messages as they are; there is no extra functionality, qed.

However, once you have a bigger and complex application, you may want to avoid duplicated logging processes and make them more systematic. For example, you could write log messages like

```go
log.Printf("[INFO] user %v has logged in", user.Id)
```

But this login information could also be written like:

```go
log.Printf("INFO: user with id=%v has logged in", user.Id)
```

This inconsistency is not good, because it
- confuses operators(or devops) cognitively, everytime they need to checkout the log messages
- is almost impossible to parse messages for obtaining useful data, usually for analytics
- makes logging process manual - you have to write log messages from scratch for all occasions

Fortunately, Go's [`log/slog`](https://pkg.go.dev/log/slog) package resolves this issue. We can write structured log messages that are good for parsing. Also, the package tries to optimize memory usage in constructing messages. So there is no reason why not using this package on behalf of our old friend `log`.

However, if you check out the package, a lot of features could make you feel daunted. So in this chapter, we would like to explain the usage of the package step by step, so that we can make the most of it. 

## Creating a Logger with Handlers
First of all, let's briefly overview how `slog` works. There are two basic parts in `slog` - [`Logger`](https://pkg.go.dev/log/slog#Logger) and [`Handler`](https://pkg.go.dev/log/slog#Logger.Handler)(please don't get confused with the one in the `net/http`). When we try to write a log message, the Logger constructs a [`Record`](https://pkg.go.dev/log/slog#Record)(which is just a collection of the information of the log) and passes it to its registered `Handler`. Then this `Handler` decides how to process the given `Record` data.

Creating a logger is simple. You need to prepare your `Handler` beforehand.

```go
// [...]
writer := os.Stdout
handlerOptions := slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}
handler := slog.NewTextHandler(writer, &handlerOptions)

logger := slog.New(handler)
// [...]
```

In `log/slog` package, there are two built-in handlers are provided - [`TextHandler`](https://pkg.go.dev/log/slog#TextHandler) and [`JsonHandler`](https://pkg.go.dev/log/slog#JSONHandler). In most cases, we won't need more than these two. Maybe XML? I have barely any idea about that. And both `TextHandler` and `JsonHandler` needs to have `io.Writer` and [`*HandlerOptions`](https://pkg.go.dev/log/slog#HandlerOptions) arguments to be generated and passed to a new `Logger`. Please checkout the `HandlerOptions` type for details, but it is pretty straightforward so I won't go into details. 

## Logging with `Attr`s
Once you created a new logger, it's time to make a log message! When making log messages with built-in `TextHandler` or `JsonHandler`, there are two kinds of information that you can provide to the logging functions - "message" and "arguments". 

A message is a plain simple text string. You can put any string you want. On the other hand, arguments are a collection of variables you put to the logging functions, which are parsed as a set of key-value pairs. For instance,

```go
logger.Info("hey", "MrChow", "Ken Jeong")
```

gives a log message as 

```sh
time=2025-08-13T00:22:24.758+09:00 level=INFO msg=hey MrChow="Ken Jeong"
```

However, internally arguments(and messages in usual cases) are converted to a type called [`Attr`](https://pkg.go.dev/log/slog#Attr). And using a generalized function [LogAttrs](https://pkg.go.dev/log/slog#Logger.LogAttrs) is recommended for unnecessary data allocation in memory. 

Hence, we will make log messages as follows, instead of using convenient methods like `.Info`:

```go
logger.LogAttrs(loggerCtx, slog.LevelInfo, "hey", slog.String("MrChow", "Ken Jeong"))
```

There are built-in basic `Attr` types and factory functions in `log/slog`, such as [`String`](https://pkg.go.dev/log/slog#String). Please look up and check out the documentation page for details. 

Note that we have now a context variable(`loggerCtx`). It is [recommended to pass contexts to the logging methods](https://pkg.go.dev/log/slog#hdr-Contexts).

## Grouping Messages for Organized Logs
When an application becomes complex, there are several business logics handling their own business logics, not necessarily separated. To recognize where our log messages come from, we would like to attach "group tags" for those messages.

In `log/slog` package, you can do it by using [`Logger.WithGroup`](https://pkg.go.dev/log/slog#Logger.WithGroup) to create a new logger, printing the tag-attached messages.

```go
// [...]
hangoverLogger := logger.WithGroup("hangover")
hangoverLogger.LogAttrs(loggerCtx, slog.LevelInfo, "hey", slog.String("MrChow", "Ken Jeong"))
```

which will give a log message like follows:

```sh
time=2025-08-13T00:40:00.020+09:00 level=INFO msg=hey hangover.MrChow="Ken Jeong"
```

You'll see that every key in key-value pairs has the group name(`"hangover"`) as its prefix prepended, while separated by `.`. This is because we are using `TextHandler`, but you'll see a big difference when you use `JsonHandler`. 

## Reducing Duplicative Logging with `With`
There could be cases in which we may print the same set of variables in multiple log messages. In such a case, we can use [`Logger.With`](https://pkg.go.dev/log/slog#Logger.With) to create a new logger that automatically put the set of variables as key-value pairs ahead of other arguments.

```go
// [...]
withLogger := hangoverLogger.With("Alan", "Zach Galifianakis")
withLogger.LogAttrs(loggerCtx, slog.LevelInfo, "hey", slog.String("MrChow", "Ken Jeong"))

// [...]
```

The log message will be like:

```sh
time=2025-08-13T00:52:24.752+09:00 level=INFO msg=hey hangover.Alan="Zach Galifianakis" hangover.MrChow="Ken Jeong"
```

## Conclusion
It wasn't that much complicated as it looked at first, right? Nevertheless, the `log/slog` package provides very comprehensive methods of making structured logs. 

To recap:
1. Define your `Handler` - either `TextHandler`, `JsonHandler`, or your own customized.
2. Create a new `Logger`.
3. (Optional) attach a group information, by creating a new logger with `Logger.WithGroup`.
4. (Optional) attach a set of key-value pairs for ergonomic logging, by creating a new loger with `Logger.With`.
5. Make a log message with `Attr` types explicitly, using `Logger.LogAttrs`.

## Exercise
Convert your `Handler`, from `TextHandler` to `JsonHandler`, and print a few logs. You'll see why I emphasized the power of `Logger.WithGroup`.
