package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//Timeslot struct
type Timeslot struct {
	date    string
	slot    string
	doctor  string
	patient *Patient
	next    *Timeslot
}

//Patient Struct
type Patient struct {
	name       string
	hasBooking bool
	timeslot   *Timeslot
}

type appointment map[string]*linkedList

var doctorList = []string{"Dr.Andy", "Dr.Bob", "Dr.Chris", "Dr.Denny"}
var slotMap = map[string]string{"1": "9am-11am", "2": "11am-1pm", "3": "2pm-4pm", "4": "4pm-6pm"}
var appointmentMap = &appointment{}
var patientBst = &BST{}
var reader = bufio.NewReader(os.Stdin)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	for {
		user := login()
	out:
		for {
			fmt.Printf("\nWelcome %v. What do you want to do today:\n", user.name)
			if user.name == "admin" {
				fmt.Println("1. Browse appointment for a doctor")
				fmt.Println("2. Search appointment for a patient")
				fmt.Println("3. Sign out")
			} else {
				fmt.Println("1. List available time slots of a doctor")
				fmt.Println("2. List available doctors of a timeslot")
				fmt.Println("3. Show my appointment")
				fmt.Println("4. Make appointment")
				fmt.Println("5. Edit appointment")
				fmt.Println("6. Sign out")
			}
			fmt.Println("\nSelect your choice:")

			c, _ := reader.ReadString('\n')
			c = strings.Replace(c, "\n", "", -1)
			choice, err := strconv.Atoi(c)

			if err != nil || (user.name == "admin" && (choice > 3 || choice < 1)) || (user.name != "admin" && (choice > 6 || choice < 1)) {
				fmt.Println("\nPlease enter a valid number")
				continue
			}

			if user.name == "admin" {
				switch choice {
				case 1:
					doctorAppointment()
				case 2:
					patientAppointment()
				case 3:
					break out
				}
			}

			if user.name != "admin" {
				switch choice {
				case 1:
					listAvailTime()
				case 2:
					listAvailDoc()
				case 3:
					showAppointment(user)
				case 4:
					makeAppointment(user)
				case 5:
					editAppointment(user)
				case 6:
					break out
				}
			}
			time.Sleep(750 * time.Millisecond)
		}
		fmt.Println("\nSigning out...")
	}
}

func login() *Patient {
	fmt.Println("\n\n============================================")
	fmt.Println("Welcome to the Online Dental Appointment System.")
	fmt.Println("Input your username to login/register.")
	fmt.Println("Input \"admin\" to login as admin.")
	fmt.Println("============================================")
	fmt.Println("\nEnter username: ")

	username, _ := reader.ReadString('\n')
	username = strings.ToLower(strings.Replace(strings.Replace(username, "\n", "", -1), " ", "", -1))

	if len(username) <= 0 {
		panic("Invalid username. Username must have at least 1 character")
	} else if len(username) > 15 {
		panic("Invalid username. Username cannot have more than 15 characters")
	}

	user := patientBst.search(username)
	if user != nil {
		return user
	}

	fmt.Println("User not found. Do you want to sign up? Press enter to sign up, input 'N' to quit.")
	var choice string
	c, _ := reader.ReadString('\n')
	choice = strings.ToUpper(strings.Trim(strings.Replace(c, "\n", "", -1), " "))

	if choice != "N" {
		user = &Patient{name: username}
		patientBst.insert(user)
		return user
	}
	panic("Quitting without signing up")
}

func listAvailTime() {
	doctor, _ := selectDoctor()
	nextWeek := time.Now().AddDate(0, 0, 7)
	currentTimeslot := (*appointmentMap)[doctor].head

	fmt.Println("\nDate         Available Slots")

	for date := time.Now(); date.Before(nextWeek); date = date.AddDate(0, 0, 1) {
		fmt.Printf("%v    ", date.Format("2006-01-02"))
		for slot := 1; slot < 5; slot++ {
			if currentTimeslot == nil {
				fmt.Printf("%v ", slot)
			} else if date.Format("2006-01-02") != currentTimeslot.date || strconv.Itoa(slot) != currentTimeslot.slot {
				fmt.Printf("%v ", slot)
			} else {
				fmt.Printf("  ")
				currentTimeslot = currentTimeslot.next
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func listAvailDoc() {
	date, _ := selectDate()
	slot, _ := selectSlot()
	availability := make(chan string)
	fmt.Println("")
	for _, doctor := range doctorList {
		go checkDocAvail(availability, doctor, date, slot)
	}
	for i := 1; i <= len(doctorList); i++ {
		fmt.Println(<-availability)
	}
}

func checkDocAvail(availability chan string, doctor string, date string, slot string) {
	currentTimeslot := (*appointmentMap)[doctor].head
	avail := true
	for currentTimeslot != nil {
		if date == currentTimeslot.date && slot == currentTimeslot.slot {
			avail = false
			break
		}
		currentTimeslot = currentTimeslot.next
	}

	if avail == true {
		availability <- fmt.Sprintf("%s is available", doctor)
	} else {
		availability <- fmt.Sprintf("%s is NOT available", doctor)
	}

}

func showAppointment(patient *Patient) {
	if patient.hasBooking == true {
		fmt.Printf("\nYou have an appointment with %v on %v, %v\n", patient.timeslot.doctor, patient.timeslot.date, slotMap[patient.timeslot.slot])
	} else {
		fmt.Println("\nYou do not have an appointment yet.")
	}
}

func makeAppointment(patient *Patient) {
	if patient.hasBooking == true {
		fmt.Println("\nYou already have an appointment.")
		return
	}
	doctor, _ := selectDoctor()
	date, _ := selectDate()
	slot, _ := selectSlot()

	fmt.Printf("\nAre you sure you want to make an appointment with %v on %v -- slot %v\nPress enter to proceed with booking, press 'N' to cancel and return to menu.\n", doctor, date, slot)
	var choice string
	c, _ := reader.ReadString('\n')
	choice = strings.ToUpper(strings.Trim(strings.Replace(c, "\n", "", -1), " "))
	if choice == "N" {
		fmt.Println("Returning to menu...")
		return
	}
	_, err := addAppointment(patient, doctor, date, slot, false)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\nBooking successful.\n\n")
	}
}

func editAppointment(patient *Patient) {
	if patient.hasBooking == false {
		fmt.Println("\nYou do not have an appointment yet.")
		return
	}
	fmt.Printf("\nYou have an appointment with %v on %v, %v\n", patient.timeslot.doctor, patient.timeslot.date, slotMap[patient.timeslot.slot])
	fmt.Println("\nAre you sure you want to edit your appointment? Press enter to proceed with booking, press 'N' to cancel and return to menu.")
	var choice string
	c, _ := reader.ReadString('\n')
	choice = strings.ToUpper(strings.Trim(strings.Replace(c, "\n", "", -1), " "))
	if choice == "N" {
		fmt.Println("Returning to menu...")
		return
	}
	doctor, _ := selectDoctor()
	date, _ := selectDate()
	slot, _ := selectSlot()

	_, err := addAppointment(patient, doctor, date, slot, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\nBooking successful.\n\n")
	}
}

func doctorAppointment() {
	doctor, _ := selectDoctor()
	currentTimeslot := (*appointmentMap)[doctor].head
	fmt.Printf("\nUpcoming appointment of %v:\n", doctor)
	for currentTimeslot != nil {
		fmt.Printf("%v %-8v %v\n", currentTimeslot.date, slotMap[currentTimeslot.slot], currentTimeslot.patient.name)
		currentTimeslot = currentTimeslot.next
	}
}

func patientAppointment() {
	patientBst.inOrder()
}
