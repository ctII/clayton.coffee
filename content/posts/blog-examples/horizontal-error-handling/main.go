package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	err := openReadWriteReplaceNewLine(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
}

func openReadWriteReplaceNewLine(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			fmt.Println(err)
		}
	}()

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

func connectDatabase() error {
	return nil
}
func connectRedis() error {
	return nil
}

func startRequirements() error {
	databaseErr := make(chan error)
	go func() { databaseErr <- connectDatabase() }()

	redisErr := make(chan error)
	go func() { redisErr <- connectRedis() }()

	err := <-databaseErr
	if err != nil {
		return err
	}

	err = <-redisErr
	if err != nil {
		return err
	}

	return nil
}

func startRequirementsSequentially() error {
	err := connectDatabase()
	if err != nil {
		return err
	}

	err = connectRedis()
	if err != nil {
		return err
	}

	return nil
}

func startRequirementsTogether() error {
	databaseErr := make(chan error)
	go func() { databaseErr <- connectDatabase() }()

	redisErr := make(chan error)
	go func() { redisErr <- connectRedis() }()

	err := <-databaseErr
	if err != nil {
		return err
	}

	err = <-redisErr
	if err != nil {
		return err
	}

	return nil
}
