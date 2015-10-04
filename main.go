package main

import (
        "log"
        "os"
        "github.com/tarm/serial"
        "fmt"
        "time"
        "net/http"
        "github.com/gorilla/mux"
        "github.com/joho/godotenv"
        "encoding/json"
        "io/ioutil"
        "gopkg.in/yaml.v2"
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
    jenkinsJson, err := ioutil.ReadAll(resp.Body)

    jobStatus := getJobsStatusFromJenkinsJson(string(jenkinsJson))

    fmt.Println("Done", jobStatus)
}

func getJobsStatusFromJenkinsJson(jenkinsJson string) string {

    type JobItem struct {
        Name string
        Url string
        Color string
    }

    type JenkinsApiData struct {
        Mode string
        NodeDescription string
        Jobs []JobItem
    }

    msg := new(JenkinsApiData)
    _ = json.Unmarshal([]byte(jenkinsJson), &msg)

    fmt.Println(os.Getenv("HELLO"))

    monitoredJob := JobItem{}

    for _, jobItem := range msg.Jobs {
        if jobItem.Name == "golights" {
            monitoredJob = jobItem;
            break;
        }
    }

    return monitoredJob.Color
}

func getFrequentStatusFromJenkins() {
    for {
        time.Sleep(time.Second * 3)
        sendRequestToJenkins()
    }
}

func getConfiguration() {
    type Config struct {
        Job_name string
    }

    t := Config{}
    data, _ := ioutil.ReadFile("config.yaml")
    //fmt.Println(string(data))
    err := yaml.Unmarshal(data, &t)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    fmt.Println(t)
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    //handleHttpRequests()
    getFrequentStatusFromJenkins()
}

