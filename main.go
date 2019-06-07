package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	url1 = "https://api.packet.net/projects"
	url2 = "https://api.packet.net/projects/ca73364c-6023-4935-9137-2132e73c20b4/devices"
	url3 = "https://api.packet.net/devices"
)

type obj struct {
	Projects []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
}

type device struct {
	Id       string `json:"id"`
	State    string `json:"state"`
	Hostname string `json: "hostname"`
}

func getProjects(body []byte) (*obj, error) {
	var s = new(obj)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func getDeviceInfo(body []byte) (*device, error) {
	var s = new(device)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}
// this function takes a string "GET" and string url and returns a response status   
func GetRequestMethod(a, b string) string {
	req, err := http.NewRequest(a, b, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", "wbrYPDxpE1y8bT95WknGyJgrwPdsteVw")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.Status
}
//Inputs: "GET" string and url string, Outputs: Device info
func RetrieveDeviceInfo(a, b string) (string, string) {
	req, err := http.NewRequest(a, b, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", "wbrYPDxpE1y8bT95WknGyJgrwPdsteVw")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	Dev, err := getDeviceInfo(body)
	if err != nil {
		panic(err)
	}
	return Dev.Id, Dev.Hostname

}
// This function creates a device and returns device info
func PostRequestMethod(a, b string) (string, string) {

	device := []byte(`{"facility": "ewr1","plan":"t1.small.x86","operating_system":"ubuntu_16_04"}`)
	req, err := http.NewRequest(a, b, bytes.NewBuffer(device))
	req.Header.Set("X-Auth-Token", "wbrYPDxpE1y8bT95WknGyJgrwPdsteVw")
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	r, err := getDeviceInfo(body)
	if err != nil {
		panic(err)
	}
	return r.Id, r.State
}
// Deletes an existing device
func DeleteDevice(a, b, c string) string {
	fmt.Println("Deleting Device...")
	url := []string{b, c}
	req, err := http.NewRequest(a, strings.Join(url, "/"), nil)
	req.Header.Set("X-Auth-Token", "wbrYPDxpE1y8bT95WknGyJgrwPdsteVw")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.Status
}

func main() {
	fmt.Println("URL:>", url1)
	status := GetRequestMethod("GET", url1)
	fmt.Println("Response Status:", status)
	fmt.Println("Ready to create device")
	deviceId, deviceStat := PostRequestMethod("POST", url2)
	fmt.Println("device: [\n   { Id:", deviceId, "\n     Status:", deviceStat, "\n   }\n]")
	fmt.Println("Datacenter is provisioning device.")
	timer1 := time.NewTimer(45 * time.Second)
	<-timer1.C
	deleteDevice := DeleteDevice("DELETE", url3, deviceId)
	fmt.Println("Response Satuts: ", deleteDevice)
	fmt.Print("Device has been deleted\nTerminating Program.")

}
