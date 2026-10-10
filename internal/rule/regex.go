package rule

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

const negate = '!'

type Regex struct {
	name         string
	exclusive    bool
	regexPattern string
	negate       bool
	*sync.RWMutex
}

func (rule *Regex) Init() Rule {
	rule.name = "regex"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Regex) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

// 0 = regex pattern
func (rule *Regex) SetParameters(params []string) error {
	rule.Lock()
	defer rule.Unlock()

	if len(params) == 0 {
		return fmt.Errorf("regex pattern not exists")
	}

	if params[0] == "" {
		return fmt.Errorf("regex pattern is empty")
	}

	if params[0][0] == negate {
		rule.negate = true
		rule.regexPattern = params[0][1:]
		return nil
	}

	rule.negate = false
	rule.regexPattern = params[0]
	return nil
}

func (rule *Regex) GetParameters() []string {
	if rule.negate {
		return []string{string(negate) + rule.regexPattern}
	}

	return []string{rule.regexPattern}
}

func (rule *Regex) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if full string matches regex
func (rule *Regex) Validate(value string, path string, _ bool) (bool, error) {
	regexPattern := rule.getRegexPattern()
	if path != "" && strings.ContainsAny(regexPattern, "$") {
		pathSplit := strings.Split(path, "/")
		replaces := make([]string, len(pathSplit)*2)
		for i := 0; i < len(pathSplit); i++ {
			replaces[i*2] = fmt.Sprintf("${%d}", len(pathSplit)-1-i)
			replaces[i*2+1] = pathSplit[i]
		}

		regexPattern = strings.NewReplacer(replaces...).Replace(regexPattern)
	}

	match, err := regexp.MatchString("^"+regexPattern+"$", value)
	return match != rule.negate, err
}

func (rule *Regex) getRegexPattern() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.regexPattern
}

func (rule *Regex) GetErrorMessage() string {
	if rule.negate {
		return fmt.Sprintf("%s:%s", rule.GetName(), string(negate)+rule.getRegexPattern())
	}

	return fmt.Sprintf("%s:%s", rule.GetName(), rule.getRegexPattern())
}

func (rule *Regex) Copy() Rule {
	rule.RLock()
	defer rule.RUnlock()

	c := new(Regex)
	c.Init()
	c.regexPattern = rule.regexPattern
	c.negate = rule.negate
	return c
}
