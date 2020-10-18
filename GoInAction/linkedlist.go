package main

import (
	"errors"
)

type linkedList struct {
	head *Timeslot
	size int
}

func (p *linkedList) add(newTimeslot *Timeslot) error {
	if p.size == 0 {
		p.head = newTimeslot
	} else if p.size == 1 && (newTimeslot.Date < p.head.Date || (newTimeslot.Date == p.head.Date && newTimeslot.Slot < p.head.Slot)) {
		newTimeslot.Next = p.head
		p.head = newTimeslot
	} else if p.size == 1 {
		p.head.Next = newTimeslot
	} else {
		currentTimeslot := p.head
		var prevTimeslot *Timeslot
		for currentTimeslot.Date < newTimeslot.Date || (currentTimeslot.Date == newTimeslot.Date && currentTimeslot.Slot < newTimeslot.Slot) {
			prevTimeslot = currentTimeslot
			currentTimeslot = currentTimeslot.Next
			if currentTimeslot == nil {
				break
			}
		}
		//currentTimeslot.Next = newTimeslot
		newTimeslot.Next = currentTimeslot
		if prevTimeslot == nil {
			p.head = newTimeslot
		} else {
			prevTimeslot.Next = newTimeslot
		}
	}
	p.size++
	return nil
}

func (p *linkedList) remove(newTimeslot *Timeslot) (*Timeslot, error) {
	var deleted *Timeslot

	if p.head == nil {
		return nil, errors.New("empty linked list")
	}
	if p.size == 1 {
		if p.head.Date == newTimeslot.Date && p.head.Slot == newTimeslot.Slot {
			deleted = p.head
			p.head = nil
		} else {
			return nil, errors.New("timeslot not found")
		}
	} else {
		var currentTimeslot *Timeslot = p.head
		var prevTimeslot *Timeslot
		for currentTimeslot.Date != newTimeslot.Date || currentTimeslot.Slot != newTimeslot.Slot {
			prevTimeslot = currentTimeslot
			currentTimeslot = currentTimeslot.Next
			if currentTimeslot == nil {
				return nil, errors.New("timeslot not found")
			}
		}
		deleted = currentTimeslot
		if prevTimeslot == nil {
			p.head = currentTimeslot.Next
		} else {
			prevTimeslot.Next = currentTimeslot.Next
		}
	}
	p.size--
	return deleted, nil
}

func (p *linkedList) get(date string, slot string) (*Timeslot, error) {
	if p.head == nil {
		return nil, errors.New("Empty Linked list")
	}
	currentNode := p.head
	for date != currentNode.Date || slot != currentNode.Slot {
		currentNode = currentNode.Next
		if currentNode == nil {
			return nil, errors.New("Timeslot not found")
		}
	}
	return currentNode, nil
}
