package main

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

func addAppointment(user *User, doctor string, date string, slot string, override bool) (*Timeslot, error) {
	if user.HasBooking == true && override == false {
		return nil, errors.New("You already have an appointment")
	}

	timeslot, _ := (*appointmentMap)[doctor].get(date, slot)
	if timeslot != nil {
		return nil, errors.New("The selected timeslot is not available")
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if user.HasBooking == true && override == true {
		removeAppointment(user)
	}

	timeslot = &Timeslot{Date: date, Slot: slot, Doctor: doctor, User: user, Next: nil}
	(*appointmentMap)[doctor].add(timeslot)
	user.HasBooking = true
	user.Timeslot = timeslot

	return timeslot, nil

}

func removeAppointment(user *User) (*Timeslot, error) {
	if user.HasBooking == false {
		return nil, errors.New("You do not have an appointment yet")
	}

	timeslot := user.Timeslot
	doctor := user.Timeslot.Doctor

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	timeslot, _ = (*appointmentMap)[doctor].remove(timeslot)

	user.HasBooking = false
	user.Timeslot = nil
	return timeslot, nil
}

func dateList() []string {
	nextWeek := time.Now().AddDate(0, 0, 7)
	dList := []string{}
	for date := time.Now(); date.Before(nextWeek); date = date.AddDate(0, 0, 1) {
		dList = append(dList, date.Format("2006-01-02"))
	}
	return dList
}

func userMap() map[string]Timeslot {
	return patientBst.inOrder()
}

func availabilityMap() map[string]map[string][]bool {
	availMap := make(map[string]map[string][]bool)
	nextWeek := time.Now().AddDate(0, 0, 7)

	for _, doctor := range doctorList {
		currentTimeslot := (*appointmentMap)[doctor].head
		availMap[doctor] = map[string][]bool{}

		for date := time.Now(); date.Before(nextWeek); date = date.AddDate(0, 0, 1) {
			availMap[doctor][date.Format("2006-01-02")] = []bool{}

			for slot := 1; slot < 5; slot++ {
				if currentTimeslot == nil {
					availMap[doctor][date.Format("2006-01-02")] = append(availMap[doctor][date.Format("2006-01-02")], true)
				} else if date.Format("2006-01-02") != currentTimeslot.Date || strconv.Itoa(slot) != currentTimeslot.Slot {
					availMap[doctor][date.Format("2006-01-02")] = append(availMap[doctor][date.Format("2006-01-02")], true)
				} else {
					availMap[doctor][date.Format("2006-01-02")] = append(availMap[doctor][date.Format("2006-01-02")], false)
					currentTimeslot = currentTimeslot.Next
				}
			}
		}
	}
	return availMap
}
