package main

import (
	"html/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func init() {

}

func init() {
	for _, doc := range doctorList {
		(*appointmentMap)[doc] = &linkedList{}
	}

	tpl = template.Must(template.ParseGlob("templates/*"))

	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)

	admin := &User{Name: "admin", Password: bPassword}
	patient1 := &User{Name: "joe", Password: bPassword}
	patient2 := &User{Name: "kelvin", Password: bPassword}
	patient3 := &User{Name: "stella", Password: bPassword}
	patient4 := &User{Name: "andrew", Password: bPassword}
	patient5 := &User{Name: "john", Password: bPassword}
	patient6 := &User{Name: "zoey", Password: bPassword}
	patient7 := &User{Name: "emma", Password: bPassword}

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
