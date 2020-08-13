package main

import (
	"container/list"
	"errors"
	"sync"

	"github.com/golang/glog"
)

type Queue struct {
	queue *list.List
	mutex sync.Mutex
}

func GetQueue() *Queue {
	return &Queue{
		queue: list.New(),
	}
}
func (this *Queue) Push(data interface{}) {
	if data == nil {
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.queue.PushBack(data)
}

func (this *Queue) Front() interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	front := this.queue.Front()
	glog.Info(front)
	return front
}

func (this *Queue) Pop() (interface{}, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	front := this.queue.Front()
	if front == nil {
		return nil, errors.New("try Pop from an empty queue")
	}
	this.queue.Remove(front)
	return front, nil
}
