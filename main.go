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

type obj struct {
	Projects []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
}

type obj1 struct {
	Id    string `json:"id"`
	State string `json:"state"`
}

func getProjects(body []byte) (*obj, error) {
	var s = new(obj)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func getDeviceInfo(body []byte) (*obj1, error) {
	var s = new(obj1)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func main() {
	// URL's with specific endpoints
	url1 := "https://api.packet.net/projects"
	url2 := "https://api.packet.net/projects/ca73364c-6023-4935-9137-2132e73c20b4/devices"
	url3 := "https://api.packet.net/devices"
	fmt.Println("URL:>", url1)
	//Requesting a GET, for recieving project detials
	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		panic(err)
	}
	//Authentication, Packet API requires an API token in the header of new request
	req.Header.Set("X-Auth-Token", "wbrYPDxpE1y8bT95WknGyJgrwPdsteVw")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println("Creating Device")
	//Creating a device with specific detials and proper JSON format
	device := []byte(`{"facility": "ewr1","plan":"t1.small.x86","operating_system":"ubuntu_16_04"}`)
	req, err = http.NewRequest("POST", url2, bytes.NewBuffer(device))
	req.Header.Set("X-Auth-Token", "wbrYPDxpE1y8bT95WknGyJgrwPdsteVw")
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	r, err := getDeviceInfo(body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Device:\n Id:", r.Id, "\nSate:", r.State)
	fmt.Println("Deleting Device...")
	//Timer allows for Device Processing after initial creation of device
	timer1 := time.NewTimer(45 * time.Second)
	<-timer1.C
	
	//Extracted device Id to then concatenate with url3
	together := []string{url3, r.Id}
	req, err = http.NewRequest("DELETE", strings.Join(together, "/"), nil)
	req.Header.Set("X-Auth-Token", "wbrYPDxpE1y8bT95WknGyJgrwPdsteVw")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)

	fmt.Print("Device has been deleted\nTerminating Program.")
}
