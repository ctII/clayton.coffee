+++ 
draft = false
date = 2023-05-06 19:08:09-05:00
title = "Handling errors in a defer"
description = "Error handling in Go is fairly simple, well, *most of the time*, and then you encounter an error while *already* handling an error."
tags = ["go"]
+++

# Error handling in Go is fairly simple, well, most of the time.

Generally you just call a function that returns an error, check it, and then 

```go
if err != nil { 
    return err 
}
```

And that covers the vast majority of error handling people do when writing in Go. 

For example, this function that replaces new lines in a file:

```go
func replaceFileNewLines(path string, replace string) error {
    f, err := os.OpenFile(path, os.O_RDWR, 0755)
    if err != nil {
        return err
    }
    defer f.Close()

    fileContents, err := io.ReadAll(f)
    if err != nil {
        return err
    }

    fileContents = bytes.ReplaceAll(fileContents, []byte("\n"), []byte(replace))

    if err = f.Truncate(0); err != nil {
        return err
    }

    if _, err = f.Seek(0, 0); err != nil {
        return err
    }

    if _, err = f.Write(fileContents); err != nil {
        return err
    }

    return nil
}
```

Generally this works, except that it ignores an error from `f.Close()`.

I'm going to ignore what potential errors that `f.Close()` can return and argue that it's better to handle the 
error even if you generally don't expect an error, or have a clear cut way to handle that error other than returning it.

By handling errors consistently, readers of the code don't have to wonder if that missed error could be troublesome later.

## Handling errors that occur in defers

Sadly, most people just write

```go
defer f.Close()
```

and call it good enough. However, this ignores a potential error from `f.Close()`

### Logging the error

One *could* just log the error if it happened

```go
defer func() {
    if err2 := f.Close(); err2 != nil {
        log.Printf("could not close file (%v)", err2)
    }
}()
```

But this kind of error handling is lackluster at best. It is now left up to whoever reads the logs (if that ever happens) to decide how the error is going to be handled, when almost everywhere else the caller receives the error. This leads to logs being written that don't really have any context.

So what is the better solution?

How about adding err2's string to err?

```go
func replaceFileNewLines(path string, replace string) (err error) {
    [...]
    defer func() {
        if err2 := f.Close(); err2 != nil {
            err = fmt.Errorf("could not close file (%v) after another error occurred (%w)", err2, err)
        }
    }()
    [...]
}
```

Note: A named return is used here `(err error)` as they are modifiable by a `defer`. See [the Go spec](https://go.dev/ref/spec#Defer_statements). See also [my tool for detecting this](https://github.com/simplylib/defermodafterreturn).

But this loses some type information since `%v` only retrieves the error as a string. This means when attempting to use the built-in Go facilities, such as [`errors.Is`](https://pkg.go.dev/errors#Is) or [`errors.As`](https://pkg.go.dev/errors#As), the error from `f.Close()` is not receivable or detectable without attempting to do horribly error prone string comparisons.

### Pre Go1.20

The better solution is to use a multierror library, such as [hashicorp's go-multierror](https://github.com/hashicorp/go-multierror) or my own [multierror](https://github.com/simplylib/multierror). Both handle errors the same way by combining them into a single error that works with [`errors.Is`](https://pkg.go.dev/errors#Is) and [`errors.As`](https://pkg.go.dev/errors#As).

```go
func replaceFileNewLines(path string, replace string) (err error) {
    [...]
    defer func() {
        if err2 := f.Close(); err2 != nil {
            err = multierror.Append(err, fmt.Errorf("could not close file (%w)", err2))
        }
    }()
    [...]
}
```

### Post Go1.20

Starting in Go 1.20 the [errors package](https://pkg.go.dev/errors) now includes [`errors.Join`](https://pkg.go.dev/errors#Join). This lets errors be combined into a single error that still works with [`errors.Is`](https://pkg.go.dev/errors#Is) and [`errors.As`](https://pkg.go.dev/errors#As).

```go
func replaceFileNewLines(path string, replace string) (err error) {
    [...]
    defer func() {
        if err2 := f.Close(); err2 != nil {
            err = errors.Join(err, fmt.Errorf("could not close file (%w)", err2))
        }
    }()
    [...]
}
```

Alternatively, [`fmt.Errorf`](https://pkg.go.dev/fmt#Errorf) now supports multiple `%w` verbs.

```go
func replaceFileNewLines(path string, replace string) (err error) {
    [...]
    defer func() {
        if err2 := f.Close(); err2 != nil {
            if err == nil {
                err = fmt.Errorf("could not close file (%w)", err2)
                return
            }
            err = fmt.Errorf("could not close file (%w) after another error occurred (%w)", err2, err)
        }
    }()
    [...]
}
```

However, I argue the intention to combine errors is clearer, not to mention faster to read, than the first option.

## Conclusion

Unlike every other part of Go development the errors from `Close`, and every other error that occurs in a deferred call, are generally ignored by Go developers even when the Go dogma is to handle an error no matter what. This happens even when CI pipelines are giving warnings about unhandled errors, or even just VSCode/Goland highlighting these unhandled errors. 

Don't let an unhandled error cause the next high priority weekend meeting.
