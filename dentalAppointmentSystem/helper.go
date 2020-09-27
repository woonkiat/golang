package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func selectDoctor() (string, error) {
	fmt.Println("\nSelect a doctor:")

	for i, doc := range doctorList {
		fmt.Printf("%v. %v\n", i+1, doc)
	}
	fmt.Println("\nSelect your choice:")

	var choice int
	var err error
	for {
		c, _ := reader.ReadString('\n')
		c = strings.Replace(c, "\n", "", -1)
		choice, err = strconv.Atoi(c)
		if err != nil || choice > len(doctorList) || choice < 1 {
			fmt.Println("\nPlease enter a valid number:")
		} else {
			return doctorList[choice-1], nil
		}
	}
}

func selectDate() (string, error) {
	fmt.Println("\nSelect a date:")

	nextWeek := time.Now().AddDate(0, 0, 7)
	i := 1

	for date := time.Now(); date.Before(nextWeek); date = date.AddDate(0, 0, 1) {
		fmt.Printf("%v. %v\n", i, date.Format("2006-01-02"))
		i++
	}
	fmt.Println("\nSelect your option:")

	var choice int
	var err error
	for {
		c, _ := reader.ReadString('\n')
		c = strings.Replace(c, "\n", "", -1)
		choice, err = strconv.Atoi(c)
		if err != nil || choice > 7 || choice < 1 {
			fmt.Println("\nPlease enter a valid number:")
		} else {
			return time.Now().AddDate(0, 0, choice-1).Format("2006-01-02"), nil
		}
	}
}

func selectSlot() (string, error) {
	fmt.Println("\nSelect a slot.")
	for i := 0; i < len(slotMap); i++ {
		fmt.Printf("%v. %v\n", i+1, slotMap[strconv.Itoa(i+1)])
	}
	fmt.Println("\nSelect your option.")

	var choice int
	var err error
	for {
		c, _ := reader.ReadString('\n')
		c = strings.Replace(c, "\n", "", -1)
		choice, err = strconv.Atoi(c)
		if err != nil || choice > 4 || choice < 1 {
			fmt.Println("\nPlease enter a valid number:")
		} else {
			return strconv.Itoa(choice), nil
		}
	}
}

func addAppointment(patient *Patient, doctor string, date string, slot string, override bool) (*Timeslot, error) {
	if patient.hasBooking == true && override == false {
		return nil, errors.New("You already have an appointment")
	}

	timeslot, _ := (*appointmentMap)[doctor].get(date, slot)
	if timeslot != nil {
		return nil, errors.New("The selected timeslot is not available")
	}

	if patient.hasBooking == true && override == true {
		removeAppointment(patient)
	}

	timeslot = &Timeslot{date, slot, doctor, patient, nil}
	(*appointmentMap)[doctor].add(timeslot)
	patient.hasBooking = true
	patient.timeslot = timeslot
	return timeslot, nil
}

func removeAppointment(patient *Patient) (*Timeslot, error) {
	if patient.hasBooking == false {
		return nil, errors.New("You do not have an appointment yet")
	}

	timeslot := patient.timeslot
	doctor := patient.timeslot.doctor

	timeslot, _ = (*appointmentMap)[doctor].remove(timeslot)
	patient.hasBooking = false
	patient.timeslot = nil
	return timeslot, nil
}
