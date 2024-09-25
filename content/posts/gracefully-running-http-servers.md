+++
title = "Graceful Go HTTP Server Shutdown"
date = "2023-09-18T14:36:30-05:00"
author = "Clayton Townsend II"
cover = ""
tags = ["go", "golang", "http"]
keywords = ["go", "golang", "http"]
description = ""
showFullContent = false
readingTime = true
draft = true
+++

Generally when I see http servers written by new Go programmers, I see [`net/http`](https://pkg.go.dev/net/http) used in a manner that does not handle shutdowns gracefully.

It seems most copy verbatim the [Go documentation](https://pkg.go.dev/net/http@go1.21.1#hdr-Servers) example into something like this:

```go
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Ok")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

While this works, it glosses over an important aspect on how this applications exits: the program running this one.

What happens if a outside program (such as docker, systemd, sysvinit, etc) wants to shutdown the application as a part of some routine like container node migration, or the host is restarting?

In this case, a problem arises. For example in docker a SIGTERM is sent to the process which as of [go1.21.1 defaults](https://pkg.go.dev/os/signal@go1.21.1#hdr-Default_behavior_of_signals_in_Go_programs) to stopping the program immediately.

This means the following program would never print "Shutdown complete".

```go
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Ok")
    })

    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Shutdown complete")
}
```

This *also* means that any request that is currently being processed gets immediately stopped without cleanup or finishing the response to the user. Likely that user just sees an [`io.EOF`](https://pkg.go.dev/io#pkg-variables) with no explaination other than the connection has closed.


