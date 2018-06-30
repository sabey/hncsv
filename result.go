package main

import (
	"fmt"
)

type Result struct {
	ID       uint64
	URL      string
	Title    string
	Points   uint64
	User     string
	Age      string
	Comments uint64
	IsJob    bool
}

func (self *Result) IsValid() bool {
	if self == nil {
		return false
	}
	return true
}
func (self *Result) GetCSVFields() []string {
	return []string{
		fmt.Sprintf("%d", self.ID),
		self.URL,
		self.Title,
		fmt.Sprintf("%d", self.Points),
		self.User,
		self.Age,
		fmt.Sprintf("%d", self.Comments),
		fmt.Sprintf("%t", self.IsJob),
	}
}
