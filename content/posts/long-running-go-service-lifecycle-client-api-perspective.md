+++
title = "Long Running Go Service Lifecycle, An API Perspective"
date = "2022-11-14T17:08:31-06:00"
author = "Clayton Townsend II"
authorTwitter = "shctii" #do not include @
tags = ["go", "golang", "lifecycle management", "api design"]
description = "A overview of different approaches to exposing service lifecycles to a user of an api"
showFullContent = false
draft = true
readingTime = true
hideComments = false
+++

# The Problem

The biggest question people likely have when first using a package with some long running goroutine or background service is how that goroutine begins and ends.

A likely use case most Go developers have seen is http.ListenAndServe

```go
package main

import "net/http"

func main() {
    err := http.ListenAndServe(":8080", nil)

}
```
