package main

import (
        "log"
        "github.com/tarm/serial"
        "fmt"
        "net/http"
)

func initializeHttpHandlers() {
        http.HandleFunc("/light/", handler)
        http.ListenAndServe(":3010", nil)
}

func initializeSerializer(color string) {
        config := &serial.Config{Name: "/dev/tty.usbserial-A60080Ig", Baud: 9600}
        s, err := serial.OpenPort(config)
        if err != nil {
                log.Fatal(err)
        }

        n, err := s.Write([]byte(color))
        if err != nil {
                log.Fatal(err)
        }

        buf := make([]byte, 128)
        n, err = s.Read(buf)
        if err != nil {
                log.Fatal(err)
        }
        log.Printf("%q", buf[:n])
}

func handler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    color := r.URL.Path[len("/light/"):]
    fmt.Fprintf(w, color);

    if color == "red" {
       initializeSerializer("1");
    }
    if color == "blue" {
        initializeSerializer("0");
    }
}

func main() {
        initializeHttpHandlers();
}

