package main

import (
	"fmt"
	"time"
)

func sendEmail(emailChan chan string, done chan bool) {
	// Simulate sending email
	defer func() { done <- true }()

	for email := range emailChan {
		fmt.Println("Sending email to", email)
		time.Sleep(time.Second)
	}
}

func main() {
	emailChan := make(chan string, 100)
	done := make(chan bool)

	go sendEmail(emailChan, done)

	for i := 0; i < 10; i++ {
		emailChan <- fmt.Sprintf("%d@gmail.com", i)
	}

	fmt.Println("Done Senring")
	close(emailChan) // This is Important
	<-done
	fmt.Println("Nikunjjjjjj")
}
