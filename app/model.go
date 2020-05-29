package main

type Employee struct {
	AccountName string `json:"accountname,omitempty" bson:"accountname,omitempty"`
	FullName    string `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string `json:"phone,omitempty" bson:"phone,omitempty"`
}
