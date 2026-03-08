package main

type option struct {
	label  string
	desc   string
	result string
	danger bool
}

var options = []option{
	{
		label:  "exp1 - Hello World",
		desc:   "Classic greeting, the OG",
		result: "Hello, World! 👋",
	},
	{
		label:  "exp2 - Farewell",
		desc:   "Sometimes you gotta dip",
		result: "Bye! See ya on the other side 👋",
	},
	{
		label:  "exp3 - Ping",
		desc:   "Check if anything is alive",
		result: "Testing... ✓ all systems nominal",
	},
	{
		label:  "exp4 - Blab",
		desc:   "Absolutely meaningless output",
		result: "Blabla blab bla blaaaa 🗣️",
	},
	{
		label:  "exp5 - Matrix",
		desc:   "You took the red pill",
		result: "Wake up, Neo... 💊",
	},
	{
		label:  "uninstall - Remove this tool",
		desc:   "Kill switch — removes binary and PATH entry",
		result: "Running uninstaller...",
		danger: true,
	},
}
