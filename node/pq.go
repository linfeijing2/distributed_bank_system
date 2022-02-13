package main

import (
	"container/heap"
	"errors"
	"strconv"
)

// PriorityQueue represents the queue
type PriorityQueue struct {
	itemHeap *itemHeap
	lookup   map[string]*item
}

// New initializes an empty priority queue.
func NewPriorityQueue() PriorityQueue {
	return PriorityQueue{
		itemHeap: &itemHeap{},
		lookup:   make(map[string]*item),
	}
}

// Len returns the number of elements in the queue.
func (p *PriorityQueue) Len() int {
	return p.itemHeap.Len()
}

// Insert inserts a new element into the queue. No action is performed on duplicate elements.
func (p *PriorityQueue) Insert(v string, priority int) {
	_, ok := p.lookup[v]
	if ok {
		return
	}

	newItem := &item{
		value:    v,
		priority: priority,
	}
	heap.Push(p.itemHeap, newItem)
	p.lookup[v] = newItem
}

// Pop removes the element with the highest priority from the queue and returns it.
// In case of an empty queue, an error is returned.
func (p *PriorityQueue) Pop() (string, error) {
	if len(*p.itemHeap) == 0 {
		return "", errors.New("empty queue")
	}

	item := heap.Pop(p.itemHeap).(*item)
	delete(p.lookup, item.value)
	return item.value, nil
}

// Peek return the value of the element with the highest priority from the queue.
// In case of an empty queue, an error is returned.
func (p *PriorityQueue) Peek() (string, error) {
	if len(*p.itemHeap) == 0 {
		return "", errors.New("empty queue")
	}

	item := (*p.itemHeap)[0]
	return item.value, nil
}

func (p *PriorityQueue) PeekPriority() (int, error) {
	if len(*p.itemHeap) == 0 {
		return -1, errors.New("empty queue")
	}

	item := (*p.itemHeap)[0]
	return item.priority, nil
}

func (p *PriorityQueue) PrintQueue() {
	if len(*p.itemHeap) == 0 {
		return
	}
	print("start printing queue: ")
	for i, v := range p.lookup {
		print(i + " " + strconv.Itoa(v.priority) + " " + msgMap[i].ProposedNodeID + ", ")
	}
	println()
}

// Delete element from the queue
// In case of an empty queue, an error is returned.
func (p *PriorityQueue) Delete(x string) bool {
	item, ok := p.lookup[x]
	if !ok {
		return false
	}

	heap.Remove(p.itemHeap, item.index)
	delete(p.lookup, item.value)
	return true
}

// UpdatePriority changes the priority of a given item.
// If the specified item is not present in the queue, no action is performed.
func (p *PriorityQueue) UpdatePriority(x string, newPriority int) bool {
	item, ok := p.lookup[x]
	if !ok {
		return false
	}

	item.priority = newPriority
	heap.Fix(p.itemHeap, item.index)
	return true
}

type itemHeap []*item

type item struct {
	value    string
	priority int
	index    int
}

func (ih *itemHeap) Len() int {
	return len(*ih)
}

func (ih *itemHeap) Less(i, j int) bool {
	if (*ih)[i].priority < (*ih)[j].priority {
		return true
	}
	if (*ih)[i].priority == (*ih)[j].priority {
		return msgMap[(*ih)[i].value].ProposedNodeID < msgMap[(*ih)[j].value].ProposedNodeID
	}
	return false
}

func (ih *itemHeap) Swap(i, j int) {
	(*ih)[i], (*ih)[j] = (*ih)[j], (*ih)[i]
	(*ih)[i].index = i
	(*ih)[j].index = j
}

func (ih *itemHeap) Push(x interface{}) {
	it := x.(*item)
	it.index = len(*ih)
	*ih = append(*ih, it)
}

func (ih *itemHeap) Pop() interface{} {
	old := *ih
	item := old[len(old)-1]
	item.index = -1
	old[len(old)-1] = nil
	*ih = old[0 : len(old)-1]
	return item
}
