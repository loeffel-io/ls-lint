package main

var definitions = map[string]Rule{
	"lowercase":  new(RuleLowercase).Init(),
	"camelcase":  new(RuleCamelCase).Init(),
	"pascalcase": new(RulePascalCase).Init(),
	"snakecase":  new(RuleSnakeCase).Init(),
}

var rules = map[string]Rule{
	"lowercase": definitions["lowercase"],

	"camelcase": definitions["camelcase"],
	"camelCase": definitions["camelcase"],

	"pascalcase": definitions["pascalcase"],
	"PascalCase": definitions["pascalcase"],

	"snakecase":  definitions["snakecase"],
	"snake_case": definitions["snakecase"],
}

type Rule interface {
	Init() Rule
	GetName() string
	Validate(value string) (bool, error)
}
