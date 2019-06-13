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
	url4 = "https://api.packet.net/devices/08cbe037-d56c-4311-af47-5ee2d3de15b7"
)

type obj struct {
	Projects []struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Device device `json:"device"`
	}
}

type device struct {
	Id               string `json:"id"`
	State            string `json:"state"`
	Hostname         string `json:"hostname"`
	Operating_System osData `json:"operating_system"`
}

type osData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//Returns project data, but is not used
func getProjects(body []byte) (*obj, error) {
	var s = new(obj)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("projects could not be retrieved:", err)
	}
	return s, err
}

//Returns device info given a response body used in func PostRequestMethod
func getDeviceInfo(body []byte) (*device, error) {
	var s = new(device)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Device info could not be retrieved:", err)
	}
	return s, err
}

//Takes a url with specified endpoints, returns response status.
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

//Returns an already exsisting device id and OS id, given "GET" /device/{id}
func GetDeviceInfo(a, b, c string) (string, osData) {
	url := []string{b, c}
	req, err := http.NewRequest(a, strings.Join(url, "/"), nil)
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
	return Dev.Id, Dev.Operating_System

}

//creates a device and returns device ID givrn "POST" /projects/{id}/deivce
func PostRequestMethod(a, b string) string {

	device := []byte(`{"facility": "ewr1","plan":"t1.small.x86","operating_system":"ubuntu_16_04","tags":["Hello "]}`)
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
	return r.Id
}

// Deletes a deivce given "DELETE" /device/{id}
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
	fmt.Println("Ready to create device.")
	deviceId := PostRequestMethod("POST", url2)
	devId, devOsId := GetDeviceInfo("GET", url3, deviceId)
	fmt.Println("Device Id:", devId, "\nOperating System Id:", devOsId.Id, "\nDatacenter is provisioning device...")
	timer1 := time.NewTimer(60 * time.Second)
	<-timer1.C
	deleteDevice := DeleteDevice("DELETE", url3, deviceId)
	fmt.Println("Response Satuts: ", deleteDevice)
	fmt.Print("Device has been deleted\nTerminating Program.")
}
