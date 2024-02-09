package pubsub

import (
	"fmt"
	"sync"
	"time"
)

const (
	PubSubInitMsg   = "PubSub system created !"
	SubCreatedMsg   = "Subscriber created id:%s Topics: %v"
	PubSubStartMsg  = "PubSub system started !"
	PublishingMsg   = "Publishing message on topic: %s"
	SubscribedMsg   = "Subscriber %s subscribed to topics: %v"
	UnsubscribedMsg = "Subscriber %s unsubscribed"
	DeliveringMsg   = "Delivering message to Subscriber %s: %s %s"
	PubSubClosedMsg = "PubSub system closed"
)

type Message struct {
	Topic   string
	Content string
}

type MsgHandler func(Message)

type Subscriber struct {
	ID       string
	Topics   []string
	Callback MsgHandler
}

type PubSub struct {
	subscribers map[string]*Subscriber
	messages    chan Message
	lock        sync.RWMutex
	done        chan struct{}
}

func NewPubSub() *PubSub {
	ps := &PubSub{
		subscribers: make(map[string]*Subscriber),
		messages:    make(chan Message),
		done:        make(chan struct{}),
	}

	ps.log(fmt.Sprintf(PubSubInitMsg))

	return ps
}

func (ps *PubSub) NewSubsciber(
	id string,
	topics []string,
	callback MsgHandler,
) *Subscriber {

	ps.log(fmt.Sprintf(SubCreatedMsg, id, topics))

	return &Subscriber{
		ID:       id,
		Topics:   topics,
		Callback: callback,
	}
}

func (ps *PubSub) Publish(topic, content string) {
	ps.log(fmt.Sprintf(PublishingMsg, topic))
	ps.log(fmt.Sprintln("Message Channel", ps.messages))

	ps.messages <- Message{Topic: topic, Content: content}
}

func (ps *PubSub) Subscribe(subscriber *Subscriber) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.subscribers[subscriber.ID] = subscriber
	ps.log(fmt.Sprintf(SubscribedMsg, subscriber.ID, subscriber.Topics))
}

func (ps *PubSub) Unsubscribe(subscriberID string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.subscribers, subscriberID)
	ps.log(fmt.Sprintf(UnsubscribedMsg, subscriberID))

}

func (ps *PubSub) Start() {

	ps.log(fmt.Sprintf(PubSubStartMsg))

	for {
		select {
		case message := <-ps.messages:
			ps.lock.RLock()

			fmt.Printf("Meassage: %v\n", message)

			for _, subscriber := range ps.subscribers {
				for _, topic := range subscriber.Topics {
					if topic == message.Topic {
						ps.log(
							fmt.Sprintf(
								DeliveringMsg,
								subscriber.ID,
								message.Topic,
								message.Content,
							),
						)
						go subscriber.Callback(message)
					}
				}
			}
			ps.lock.RUnlock()
		case <-ps.done:
			ps.log(fmt.Sprintf(PubSubClosedMsg))
			return
		}
	}
}

func (ps *PubSub) Close() {
	close(ps.done)
}

func (ps *PubSub) log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), message)
}
