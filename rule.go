package main

var definitions = map[string]Rule{
	"lowercase":  new(RuleLowercase).Init(),
	"camelcase":  new(RuleCamelCase).Init(),
	"pascalcase": new(RulePascalCase).Init(),
	"snakecase":  new(RuleSnakeCase).Init(),
	"kebabcase":  new(RuleKebabCase).Init(),
	"pointcase":  new(RulePointCase).Init(),
}

var rules = map[string]Rule{
	"lowercase": definitions["lowercase"],

	"camelcase": definitions["camelcase"],
	"camelCase": definitions["camelcase"],

	"pascalcase": definitions["pascalcase"],
	"PascalCase": definitions["pascalcase"],

	"snakecase":  definitions["snakecase"],
	"snake_case": definitions["snakecase"],

	"kebabcase":  definitions["kebabcase"],
	"kebab-case": definitions["kebabcase"],

	"pointcase":  definitions["pointcase"],
	"point.case": definitions["pointcase"],
}

type Rule interface {
	Init() Rule
	GetName() string
	Validate(value string) (bool, error)
}
