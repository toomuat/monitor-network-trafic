// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Message struct {
	Os      string `json:"os"`
	Counter uint64 `json:"counter"`
}

func main() {
	bytes, err := ioutil.ReadFile("sample_data.json")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%v\n", bytes)

	var messages [3]Message
	err = json.Unmarshal(bytes, &messages)
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range messages {
		fmt.Printf("%s: %d\n", message.Os, message.Counter)
	}

	fmt.Printf("%v\n", messages)

	messages = [3]Message{
		{"Windows", 76},
		{"Android", 77},
		{"iPhone", 323},
	}

	fmt.Printf("%v\n", messages)

	var messages_map map[string]Message
	messages_map = map[string]Message{
		"Windows": Message{"Windows", 76},
		"Android": Message{"Android", 77},
		"iPhone":  Message{"iPhone", 323},
	}

	fmt.Printf("%v\n", messages_map)
	fmt.Printf("%v\n", messages_map["Windows"])
	fmt.Printf("%v\n", messages_map["Android"])
	fmt.Printf("%v\n", messages_map["iPhone"])

	var messages_array []Message
	messages_array = []Message{messages_map["Windows"], messages_map["Android"], messages_map["iPhone"]}
	fmt.Printf("%v\n", messages_array)
	messages_array2 := [3]Message{messages_map["Windows"], messages_map["Android"], messages_map["iPhone"]}
	fmt.Printf("%v\n", messages_array2)
	var messages_array_append []Message
	for key, msg := range messages_map {
		fmt.Printf("%v: %v\n", key, msg)
		messages_array_append = append(messages_array_append, msg)
	}
	fmt.Printf("**%v\n", messages_array_append)

	var os [3]string
	os[0] = "Windows"
	os[1] = "Android"
	os[2] = "iPhone"
	fmt.Printf("%v\n", os)

	jsonBytes, err := json.Marshal(messages_array_append)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%v\n", jsonBytes)
	jsonStr := string(jsonBytes)
	fmt.Printf("%v\n", jsonStr)

	fmt.Printf("%v\n", messages_map["Windows"].Counter)
}
