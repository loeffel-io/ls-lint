package rule

var RulesIndex = map[string]Rule{
	"lowercase": new(Lowercase).Init(),
	"regex":     new(Regex).Init(),
	"not_regex": new(NotRegex).Init(),
	"exists":    new(Exists).Init(),

	"camelcase":          new(CamelCase).Init(),
	"pascalcase":         new(PascalCase).Init(),
	"snakecase":          new(SnakeCase).Init(),
	"screamingsnakecase": new(ScreamingSnakeCase).Init(),
	"kebabcase":          new(KebabCase).Init(),
}

var Rules = map[string]Rule{
	"lowercase": RulesIndex["lowercase"],
	"regex":     RulesIndex["regex"],
	"not_regex": RulesIndex["not_regex"],
	"exists":    RulesIndex["exists"],

	"camelcase": RulesIndex["camelcase"],
	"camelCase": RulesIndex["camelcase"],

	"pascalcase": RulesIndex["pascalcase"],
	"PascalCase": RulesIndex["pascalcase"],

	"snakecase":  RulesIndex["snakecase"],
	"snake_case": RulesIndex["snakecase"],

	"screamingsnakecase":   RulesIndex["screamingsnakecase"],
	"SCREAMING_SNAKE_CASE": RulesIndex["screamingsnakecase"],

	"kebabcase":  RulesIndex["kebabcase"],
	"kebab-case": RulesIndex["kebabcase"],
}

type Rule interface {
	Init() Rule
	GetName() string
	SetParameters(params []string) error
	GetParameters() []string
	GetExclusive() bool
	Validate(value string, fail bool) (bool, error)
	GetErrorMessage() string
	Copy() Rule
}
