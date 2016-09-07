package main

import (
	"errors"
	"sync"
	"time"
)

//ErrQueueFull is returned when trying to push an item onto a full queue
var ErrQueueFull = errors.New("Queue full")

//WorkQueue is an interface for work item queues
type WorkQueue interface {
	Push(Identifiable) (int, error) //push item to the queue, setting its Id and returning it
	Peek() (Identifiable, error)
	Fail(id int) error //put peeked item back into main queue
	Delete(id int) error
}

type leasedItem struct {
	item       Identifiable
	leasedTime time.Time
}

//MemoryQueue is a FIFO in-memory queue
type MemoryQueue struct {
	sync.Mutex
	lastID      int
	queue       []Identifiable
	leasedItems []leasedItem
}

//NewMemoryQueue creates a new in-memory FIFO queue.
func NewMemoryQueue(capacity int) *MemoryQueue {
	return &MemoryQueue{
		queue:       make([]Identifiable, 0, capacity),
		leasedItems: make([]leasedItem, 0, capacity),
	}
}

//Push pushes the item onto the queue and returns the Id
func (mq *MemoryQueue) Push(item Identifiable) (int, error) {
	mq.Lock()
	defer mq.Unlock()
	if len(mq.queue)+len(mq.leasedItems) < cap(mq.queue) {
		item.SetId(mq.lastID)
		id := mq.lastID
		mq.lastID++
		mq.queue = append(mq.queue, item)
		return id, nil
	}
	return -1, ErrQueueFull
}

//Delete permanently removes the item from the queue by Id
func (mq *MemoryQueue) Delete(id int) error {
	mq.Lock()
	defer mq.Unlock()

	//1 first try deleting from leased items, since we assume that Peek() was used to get it in the first place
	for i := range mq.leasedItems {
		if mq.leasedItems[i].item.ID() == id {
			mq.leasedItems = append(mq.leasedItems[:i], mq.leasedItems[i+1:]...)
			return nil
		}
	}

	//2 now try deleting from queue items
	for i := range mq.queue {
		if mq.queue[i].ID() == id {
			mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
			return nil
		}
	}
	return errors.New("Item not found")
}

//Peek returns the item at the head of the queue
func (mq *MemoryQueue) Peek() (Identifiable, error) {
	mq.Lock()
	defer mq.Unlock()
	if len(mq.queue) == 0 {
		return nil, nil
	}

	item := mq.queue[0]
	mq.queue = mq.queue[1:]
	mq.leasedItems = append(mq.leasedItems, leasedItem{
		item:       item,
		leasedTime: time.Now(),
	})
	return item, nil
}

//Expire returns leased items to the main queue after the timeout so they are available from other agents
func (mq *MemoryQueue) Expire(timeout int) {
	mq.Lock()
	defer mq.Unlock()
	for i := range mq.leasedItems {
		if mq.leasedItems[i].leasedTime.Add(time.Second * time.Duration(timeout)).Before(time.Now()) {
			item := mq.leasedItems[i]
			mq.leasedItems = append(mq.leasedItems[:i], mq.leasedItems[i+1:]...)
			mq.queue = append(mq.queue, item.item)
		}

	}
}

func (mq *MemoryQueue) Fail(id int) error {
	mq.Lock()
	defer mq.Unlock()
	for i := range mq.leasedItems {
		if mq.leasedItems[i].item.ID() == id {
			item := mq.leasedItems[i]
			mq.leasedItems = append(mq.leasedItems[:i], mq.leasedItems[i+1:]...)
			mq.queue = append(mq.queue, item.item)
			return nil
		}
	}
	return errors.New("Item not found")
}
