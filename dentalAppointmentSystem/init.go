package main

import (
	"time"
)

func init() {
	for _, doc := range doctorList {
		(*appointmentMap)[doc] = &linkedList{}
	}

	admin := &Patient{name: "admin"}
	patient1 := &Patient{name: "joe"}
	patient2 := &Patient{name: "kelvin"}
	patient3 := &Patient{name: "stella"}
	patient4 := &Patient{name: "andrew"}
	patient5 := &Patient{name: "john"}
	patient6 := &Patient{name: "zoey"}
	patient7 := &Patient{name: "emma"}

	patientBst.insert(admin)
	patientBst.insert(patient1)
	patientBst.insert(patient2)
	patientBst.insert(patient3)
	patientBst.insert(patient4)
	patientBst.insert(patient5)
	patientBst.insert(patient6)
	patientBst.insert(patient7)

	addAppointment(patient1, "Dr.Andy", time.Now().AddDate(0, 0, 2).Format("2006-01-02"), "1", false)
	addAppointment(patient2, "Dr.Bob", time.Now().AddDate(0, 0, 2).Format("2006-01-02"), "1", false)
	addAppointment(patient3, "Dr.Chris", time.Now().AddDate(0, 0, 2).Format("2006-01-02"), "1", false)
	addAppointment(patient4, "Dr.Andy", time.Now().AddDate(0, 0, 2).Format("2006-01-02"), "3", false)
	addAppointment(patient5, "Dr.Andy", time.Now().AddDate(0, 0, 1).Format("2006-01-02"), "2", false)
	addAppointment(patient6, "Dr.Chris", time.Now().AddDate(0, 0, 3).Format("2006-01-02"), "2", false)
}
