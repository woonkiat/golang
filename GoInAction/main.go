package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//Timeslot struct
type Timeslot struct {
	Date   string
	Slot   string
	Doctor string
	User   *User
	Next   *Timeslot
}

//User Struct
type User struct {
	Name       string
	Password   []byte
	HasBooking bool
	Timeslot   *Timeslot
}

var tpl *template.Template

type appointment map[string]*linkedList

var doctorList = []string{"Dr.Andy", "Dr.Bob", "Dr.Chris", "Dr.Denny"}
var slotMap = map[string]string{"1": "9am-11am", "2": "11am-1pm", "3": "2pm-4pm", "4": "4pm-6pm"}
var appointmentMap = &appointment{}
var sessionMap = map[string]string{}
var patientBst = &BST{}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/makeappointment", makeAppointment)
	http.HandleFunc("/adminedit", adminEdit)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":5221", nil)
}

//struct serving as input to html templates
type pageData struct {
	Title           string
	User            User
	Timeslot        Timeslot
	ErrorMsg        string
	DateList        []string
	DoctorList      []string
	SlotMap         map[string]string
	AvailabilityMap map[string]map[string][]bool
	UserMap         map[string]Timeslot
}

func index(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	pd := pageData{Title: "Home Page"}

	//populate inputs for template
	if myUser != nil {
		pd.User = *myUser
		if myUser.HasBooking == true {
			pd.Timeslot = *(myUser.Timeslot)
		}
	}
	tpl.ExecuteTemplate(res, "index.gohtml", pd)
}

func signup(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	pd := pageData{Title: "Sign up"}
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		password := req.FormValue("password")
		if username != "" && password != "" {
			// check if username exist/ taken
			if patientBst.search(username) != nil {
				pd.ErrorMsg = "Username already taken"
				tpl.ExecuteTemplate(res, "signup.gohtml", pd)
				return
			}
			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}
			http.SetCookie(res, myCookie)
			sessionMap[myCookie.Value] = username
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			myUser := &User{Name: username, Password: bPassword}
			patientBst.insert(myUser)
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return

		}
		pd.ErrorMsg = "Invalid username or password"
		tpl.ExecuteTemplate(res, "signup.gohtml", pd)
		return

	}
	tpl.ExecuteTemplate(res, "signup.gohtml", pd)
}

func login(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	pd := pageData{Title: "Login"}
	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		// check if user exist with username
		myUser := patientBst.search(username)
		if myUser == nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)
		sessionMap[myCookie.Value] = username
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "login.gohtml", pd)
}

func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(sessionMap, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func makeAppointment(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if myUser == nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		booking := req.FormValue("booking")
		if booking == "" {
			http.Error(res, "Please select a timeslot", http.StatusForbidden)
			return
		}
		bookingData := strings.Split(booking, "|")
		doctor := bookingData[0]
		date := bookingData[1]
		sslot, _ := strconv.Atoi(bookingData[2])
		slot := strconv.Itoa(sslot + 1)

		_, err := addAppointment(myUser, doctor, date, slot, true)

		if err != nil {
			http.Error(res, err.Error(), http.StatusForbidden)
			return
		}
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	pd := pageData{Title: "Make Appointment"}
	if myUser != nil {
		pd.User = *myUser
		if myUser.HasBooking == true {
			pd.Timeslot = *(myUser.Timeslot)
		}
	}
	pd.DateList = dateList()
	pd.DoctorList = doctorList
	pd.AvailabilityMap = availabilityMap()
	pd.SlotMap = slotMap

	tpl.ExecuteTemplate(res, "makeappointment.gohtml", pd)

}

func adminEdit(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if myUser == nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	} else if myUser.Name != "admin" {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	pd := pageData{Title: "Admin : Edit appointment"}

	if req.Method == http.MethodPost {
		u := req.FormValue("user")
		if u == "" {
			http.Error(res, "Please select a user to edit", http.StatusForbidden)
			return
		}
		user := patientBst.search(u)
		if user == nil {
			http.Error(res, "User not found", http.StatusForbidden)
			return
		}

		booking := req.FormValue("booking")
		if booking == "" {
			http.Error(res, "Please select a timeslot", http.StatusForbidden)
			return
		}
		if booking == "remove" {
			_, err := removeAppointment(user)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return
		}

		bookingData := strings.Split(booking, "|")
		doctor := bookingData[0]
		date := bookingData[1]
		sslot, _ := strconv.Atoi(bookingData[2])
		slot := strconv.Itoa(sslot + 1)

		_, err := addAppointment(user, doctor, date, slot, true)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}

	pd.UserMap = userMap()
	pd.DateList = dateList()
	pd.DoctorList = doctorList
	pd.AvailabilityMap = availabilityMap()
	pd.SlotMap = slotMap

	tpl.ExecuteTemplate(res, "adminedit.gohtml", pd)

}

func getUser(res http.ResponseWriter, req *http.Request) *User {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
	}
	http.SetCookie(res, myCookie)
	// if the user exists already, get user
	var myUser *User
	if username, ok := sessionMap[myCookie.Value]; ok {
		myUser = patientBst.search(username)
	}
	return myUser
}

func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := sessionMap[myCookie.Value]
	if patientBst.search(username) != nil {
		return true
	}
	return false
}
