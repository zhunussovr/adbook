package model

// Person is main entity here
type Person struct {
	ID        string `json:"accountname,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	FullName  string `json:"fullname,omitempty"`
	Title     string `json:"title,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}
