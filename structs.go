package main

type enotfounderr struct{ name string }

func (e *enotfounderr) Error() string {
	return "strict mode on: environment variable not found: $" + e.name
}
