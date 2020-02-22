package main

var rules = map[string]Rule{
	"lowercase":  new(RuleLowercase).Init(),

	"camelcase":  new(RuleCamelCase).Init(),
	"camelCase":  new(RuleCamelCase).Init(),

	"pascalcase": new(RulePascalCase).Init(),
	"PascalCase": new(RulePascalCase).Init(),
}

type Rule interface {
	Init() Rule
	GetName() string
	Validate(value string) (bool, error)
}
