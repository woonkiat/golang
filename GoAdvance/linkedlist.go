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
	} else if p.size == 1 && newTimeslot.date <= p.head.date && newTimeslot.slot < p.head.slot {
		newTimeslot.next = p.head
		p.head = newTimeslot
	} else if p.size == 1 {
		p.head.next = newTimeslot
	} else {
		currentTimeslot := p.head
		var prevTimeslot *Timeslot
		for currentTimeslot.date < newTimeslot.date || (currentTimeslot.date == newTimeslot.date && currentTimeslot.slot < newTimeslot.slot) {
			prevTimeslot = currentTimeslot
			currentTimeslot = currentTimeslot.next
			if currentTimeslot == nil {
				break
			}
		}
		newTimeslot.next = currentTimeslot
		if prevTimeslot == nil {
			p.head = newTimeslot
		} else {
			prevTimeslot.next = newTimeslot
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
		deleted = p.head
		p.head = nil
	} else {
		var currentTimeslot *Timeslot = p.head
		var prevTimeslot *Timeslot
		for currentTimeslot.date != newTimeslot.date && currentTimeslot.slot != newTimeslot.slot {
			prevTimeslot = currentTimeslot
			currentTimeslot = currentTimeslot.next
			if currentTimeslot == nil {
				return nil, errors.New("timeslot not found")
			}
		}
		deleted = currentTimeslot
		if prevTimeslot == nil {
			p.head = currentTimeslot.next
		} else {
			prevTimeslot.next = currentTimeslot.next
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
	for date != currentNode.date || slot != currentNode.slot {
		currentNode = currentNode.next
		if currentNode == nil {
			return nil, errors.New("Timeslot not found")
		}
	}
	return currentNode, nil
}
