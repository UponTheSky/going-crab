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
handlerOptions := slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
handler := slog.NewTextHandler(writer, &handlerOptions)

logger := slog.New(handler)
// [...]
```

## Logging with `Attr`s

## Grouping Messages for Organized Logs

## Optimizing Logging with `With`

## Conclusion

## Exercise
