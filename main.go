package pubsub

import (
	"fmt"
	"time"
)

func main() {
	pubsub := NewPubSub()

	subscriber1 := &Subscriber{ID: "Sub1", Topics: []string{"topic1", "topic2"}, Callback: trigger}
	subscriber2 := &Subscriber{ID: "Sub2", Topics: []string{"topic2", "topic3"}, Callback: trigger}

	pubsub.Subscribe(subscriber1)
	pubsub.Subscribe(subscriber2)

	go pubsub.Start()

	pubsub.Publish("topic1", "Hello, topic1!")
	pubsub.Publish("topic2", "Hello, topic2!")
	pubsub.Publish("topic3", "Hello, topic3!")

	// Signal to close and wait a moment for the goroutine to finish
	pubsub.Close()
	time.Sleep(time.Second) // Adjust the sleep time as needed
}

func trigger(msg Message) {
	fmt.Printf("Triggered %s: %s\n", msg.Topic, msg.Content)
}
