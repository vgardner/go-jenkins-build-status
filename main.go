package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tarm/serial"
)

// API represents an API message.
type API struct {
	Message string "json:message"
}

// Connect to Arduino board by USB serial and send the color code.
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

// Handles request to chnage light color.
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

// Sets the color of the light on the Arduino board.
func setArduinoLightColor(lightColorCode string) {
	fmt.Println("Done", lightColorCode)
	initializeSerializer(lightColorCode)
}

// Initializes server and routes http requests.
func handleHTTPRequests() {
	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/light/{color}", handleLightColorRequest)
	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":3010", nil)
}

// Send request to Jenkins API and retrieves the status Json.
func sendRequestToJenkinsAPI() string {
	resp, err := http.Get("http://localhost:8080/api/json?pretty=true")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	jenkinsJSON, err := ioutil.ReadAll(resp.Body)

	return string(jenkinsJSON)
}

// Parses the JSON from Jenkins and retrives the specified job's status code.
func getJobsStatusFromJenkinsJSON(jenkinsJSON string) string {

	type JobItem struct {
		Name  string
		URL   string
		Color string
	}

	type JenkinsAPIData struct {
		Mode            string
		NodeDescription string
		Jobs            []JobItem
	}

	msg := new(JenkinsAPIData)
	_ = json.Unmarshal([]byte(jenkinsJSON), &msg)

	monitoredJob := JobItem{}

	for _, jobItem := range msg.Jobs {
		if jobItem.Name == os.Getenv("job_name") {
			monitoredJob = jobItem
			break
		}
	}

	return monitoredJob.Color
}

// Handles the polling of requests to Jenkins.
func getFrequentStatusFromJenkins() {
	var pollingFrequency int64

	pollingFrequency, _ = strconv.ParseInt(os.Getenv("polling_frequency"), 10, 64)

	for {
		time.Sleep(time.Second * time.Duration(pollingFrequency))

		jenkinsJSON := sendRequestToJenkinsAPI()

		jobStatus := getJobsStatusFromJenkinsJson(jenkinsJSON)

		lightColorCode := getLightColorCode(jobStatus)

		setArduinoLightColor(lightColorCode)

		fmt.Println("Done", jobStatus)
	}
}

// Get associated light code from Jenkins job for Arduino.
func getLightColorCode(statusColor string) string {

	lightColorCodes := map[string]string{
		"blue":           "1",
		"blue_anime":     "0",
		"yellow":         "2",
		"yellow_anime":   "20",
		"red":            "3",
		"red_anime":      "30",
		"grey":           "4",
		"grey_anime":     "40",
		"aborted":        "4",
		"aborted_anime":  "40",
		"disabled":       "4",
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
