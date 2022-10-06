package main

type ruleTest struct {
	value    string
	expected bool
	err      error
}

type readConfigFileTestCase struct {
	name					string
	cwd						string
	config_file 	string
	expected_error string
}
