+++ 
draft = true
date = 2022-08-29T16:12:07-05:00
title = "Horizontal Error Handling in Go"
description = "Error handling in Go is fairly simple, well, *most of the time*, and then you encounter an error while *already* handling an error."
tags = ["go"]
+++

# Error handling in Go is fairly simple, well, most of the time.
Generally you just call a function that returns an error, check it, and then 
``` go
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

Generally this works, except that it ignores an error from f.Close().

I'm going to ignore what potential errors that f.Close() can return and argue that it better to handle the 
error even if you generally don't expect an error, even based on the documentation stating it is fine to do so.

By handling errors consistently, readers of the code don't have to wonder if that missed error could be trouble later.

## Handling errors that generally happen in defers
Handling these errors is sadly generally done as was in replacefileNewLines. Most just

```go
f.Close()
```

and call it good enough. But what other solutions are there?

### Logging the error

One *could* simply just log the error if it happened:

```go
defer func() {
	if err := f.Close(); err != nil {
		log.Printf("could not close file (%v)", err)
	}
}()
```

But this kind of error handling is lackluster at best. It is now left up to whoever reads the logs (if that ever happens) to decide how the error is going to be handled, when almost everywhere else we let the caller handle 
the error. Not to mention this leads to logs being written that don't really have any context 
