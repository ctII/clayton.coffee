package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

func main() {
	err := replaceFileNewLines(os.Args[1], "---")
	if err != nil {
		panic(err)
	}
}

func replaceFileNewLines(path string, replace string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("could not close file (%v)", err)
		}
	}()

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
