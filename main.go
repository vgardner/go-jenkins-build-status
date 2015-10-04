package main

import (
        "log"
        "github.com/tarm/serial"
        "fmt"
        "time"
        "net/http"
        "github.com/gorilla/mux"
        "encoding/json"
        "io/ioutil"
)

type API struct {
    Message string "json:message"
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

func handleLightColor(w http.ResponseWriter, r *http.Request) {
    urlParams := mux.Vars(r)
    color := urlParams["color"]
    HelloMessage := "Hello, " + color

    message := API{HelloMessage}
    output, err := json.Marshal(message)

    if err != nil {
        fmt.Println("Something went wrong!")
    }

    if color == "red" {
       initializeSerializer("1");
    }
    if color == "blue" {
        initializeSerializer("0");
    }

    fmt.Fprintf(w, string(output))
}

func handleHttpRequests() {
    gorillaRoute := mux.NewRouter()
    gorillaRoute.HandleFunc("/light/{color}", handleLightColor)
    http.Handle("/", gorillaRoute)
    http.ListenAndServe(":3010", nil)
}

func sendRequestToJenkins() {
    resp, err := http.Get("http://localhost:8080/api/json?pretty=true")

    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    type B struct {
        Mode string
        NodeDescription string
        Jobs string
    }

    msg := new(B)
    _ = json.Unmarshal([]byte(string(body)), &msg)
    fmt.Println("Done", msg.NodeDescription)

    //log.Printf(string(body))
}

func getFrequentStatusFromJenkins() {
    for {
        time.Sleep(time.Second * 3)
        sendRequestToJenkins()
    }
}

func main() {
    //handleHttpRequests()
    getFrequentStatusFromJenkins()
}

