+++ 
draft = true
date = 2022-08-29T16:12:07-05:00
title = "Horizontal Error Handling in Go"
tags = ["go"]
+++

# Errors in Go are fairly simple most of the time.
You just call a function that returns an error, check it, maybe wrap it with context, then return ```if err != nil```. And that covers the vast majority of error handling people do when writing Go. 

For example, this simple function that replaces new lines in a file:

```go
func openReadWriteReplaceNewLine(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	bs = bytes.ReplaceAll(bs, []byte("\n"), []byte(""))

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(bs)
	if err != nil {
		return err
	}

	return nil
}
```
