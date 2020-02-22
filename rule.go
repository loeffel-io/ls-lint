package main

var definitions = map[string]Rule{
	"lowercase":  new(RuleLowercase).Init(),
	"camelcase":  new(RuleCamelCase).Init(),
	"pascalcase": new(RulePascalCase).Init(),
	"sneakcase":  new(RuleSneakCase).Init(),
}

var rules = map[string]Rule{
	"lowercase": definitions["lowercase"],

	"camelcase": definitions["camelcase"],
	"camelCase": definitions["camelcase"],

	"pascalcase": definitions["pascalcase"],
	"PascalCase": definitions["pascalcase"],

	"sneakcase":  definitions["sneakcase"],
	"sneak_case": definitions["sneakcase"],
}

type Rule interface {
	Init() Rule
	GetName() string
	Validate(value string) (bool, error)
}
