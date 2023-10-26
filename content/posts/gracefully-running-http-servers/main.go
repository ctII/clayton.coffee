package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main1() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ok")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Shutdown complete")
}

func main() {
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, os.Interrupt)

	httpErr := make(chan error, 1)
	server := http.Server{Addr: ":8080"}
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Ok")
		})

		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		httpErr <- err
	}()

	var err error
	select {
	case err = <-httpErr:
	case <-osSig:
		err = errors.Join(server.Shutdown(context.Background()), <-httpErr)
	}

	if err != nil {
		log.Fatalf("failed to shutdown http server %v", err)
	}

	fmt.Println("Shutdown complete")
}
