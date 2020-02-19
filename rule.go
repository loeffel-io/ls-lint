package main

var rules = map[string]Rule{
	"lowercase": new(RuleLowercase).Init(),
}

type Rule interface {
	Init() Rule
	GetName() string
	Validate(value string) (bool, error)
}
