package main

import (
	//"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQueuePush(t *testing.T) {

	Convey("When an item is pushed to an empty queue", t, func() {
		queue := NewMemoryQueue(1)
		item := WorkItem{}
		_, err := queue.Push(&item)

		Convey("It should succeed", func() {
			So(err, ShouldBeNil)
			So(len(queue.queue), ShouldEqual, 1)
		})
		Convey("When another item is pushed, and the queue is full", func() {
			_, err = queue.Push(&WorkItem{})
			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("When two items are pushed to a queue", t, func() {
		queue := NewMemoryQueue(3)
		id1, err := queue.Push(&WorkItem{})
		So(err, ShouldBeNil)
		//id2, err := queue.Push(&WorkItem{})
		So(err, ShouldBeNil)
		Convey("When peeking an item", func() {
			Convey("The first item should be returned", func() {
				item, err := queue.Peek()
				So(err, ShouldBeNil)
				So(item.ID(), ShouldEqual, id1)
			})
		})
		Convey("When peeking a second item", func() {
			queue := NewMemoryQueue(3)
			_, err := queue.Push(&WorkItem{})
			So(err, ShouldBeNil)
			id2, err := queue.Push(&WorkItem{})
			So(err, ShouldBeNil)
			_, err = queue.Peek()
			So(err, ShouldBeNil)
			Convey("The second item should be returned", func() {
				item, err := queue.Peek()
				So(err, ShouldBeNil)
				So(item.ID(), ShouldEqual, id2)
			})
		})
	})

	Convey("When an item is pushed to the queue", t, func() {
		queue := NewMemoryQueue(3)
		id1, err := queue.Push(&WorkItem{})
		So(err, ShouldBeNil)
		_, err = queue.Peek()
		So(err, ShouldBeNil)
		Convey("When this item is later deleted", func() {
			err := queue.Delete(id1)
			So(err, ShouldBeNil)
			Convey("It should be permanently removed", func() {
				So(len(queue.queue), ShouldEqual, 0)
				So(len(queue.leasedItems), ShouldEqual, 0)
			})
		})
	})
}
