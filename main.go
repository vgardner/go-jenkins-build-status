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
        "strconv"
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

func handleLightColorRequest(w http.ResponseWriter, r *http.Request) {
    urlParams := mux.Vars(r)
    color := urlParams["color"]
    HelloMessage := "Hello, " + color

    message := API{HelloMessage}
    output, err := json.Marshal(message)

    if err != nil {
        fmt.Println("Something went wrong!")
    }

    fmt.Fprintf(w, string(output))
}

func setArduinoLightColor(lightColorCode string) {
    fmt.Println("Done", lightColorCode)
    initializeSerializer(lightColorCode);
}

func handleHttpRequests() {
    gorillaRoute := mux.NewRouter()
    gorillaRoute.HandleFunc("/light/{color}", handleLightColorRequest)
    http.Handle("/", gorillaRoute)
    http.ListenAndServe(":3010", nil)
}

func sendRequestToJenkinsAPI() string {
    resp, err := http.Get("http://localhost:8080/api/json?pretty=true")

    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()
    jenkinsJson, err := ioutil.ReadAll(resp.Body)

    return string(jenkinsJson)
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

    monitoredJob := JobItem{}

    for _, jobItem := range msg.Jobs {
        if jobItem.Name == os.Getenv("job_name") {
            monitoredJob = jobItem;
            break;
        }
    }

    return monitoredJob.Color
}

func getFrequentStatusFromJenkins() {
    var pollingFrequency int64

    pollingFrequency, _ = strconv.ParseInt(os.Getenv("polling_frequency"), 10, 64)

    for {
        time.Sleep(time.Second * time.Duration(pollingFrequency))

        jenkinsJson := sendRequestToJenkinsAPI()

        jobStatus := getJobsStatusFromJenkinsJson(jenkinsJson)

        lightColorCode := getLightColorCode(jobStatus)

        setArduinoLightColor(lightColorCode)

        fmt.Println("Done", jobStatus)
    }
}

func getLightColorCode(statusColor string) string{

    lightColorCodes := map[string]string{
      "blue": "1",
      "blue_anime": "0",
      "yellow": "2",
      "yellow_anime": "20",
      "red": "3",
      "red_anime": "30",
      "grey": "4",
      "grey_anime": "40",
      "aborted": "4",
      "aborted_anime": "40",
      "disabled": "4",
      "disabled_anime": "40",
    }

    return lightColorCodes[statusColor]
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    //handleHttpRequests()
    getFrequentStatusFromJenkins()
}

