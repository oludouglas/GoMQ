package main

import (
	"encoding/json"
	"fmt"
)

type Employee struct {
	ID     float32 `json:"id"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
}

func (q Employee) String() string {
	return fmt.Sprintf("ID: %.0f, Name: %s, Salary %.2f", q.ID, q.Name, q.Salary)
}

func main() {

	activeMQ := NewClient("192.168.1.103:61616")
	activeMQ.Publish("/helloQueue", []byte("This is a simple message"))

	activeMQ.Subscribe("/employeeQueue", func(err error, message []byte) {
		var employee Employee
		if err := json.Unmarshal(message, &employee); err != nil {
			panic(err)
		}
		fmt.Println("received", employee)
		//publish the data back onto another queue
		activeMQ.Publish("helloQueue", message)
	})
}
