package domain_test

import (
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/testkit"
)

var annaOpensHerFirstAccount = commands.OpenAccountForNewCustomer{
	CustomerID:   "C010",
	CustomerName: "Anna Smith",
	AccountID:    "A010",
	AccountName:  "Anna Smith",
}

var annaOpensASecondAccount = commands.OpenAccount{
	CustomerID:  "C010",
	AccountID:   "A011",
	AccountName: "Anna Smith",
}

var annaHasAnOpenAccount = testkit.
	Scenario("anna has an open account").
	ExecuteCommand(annaOpensHerFirstAccount)

var annaHasTwoOpenAccounts = annaHasAnOpenAccount.
	Scenario("anna has two open accounts").
	ExecuteCommand(annaOpensASecondAccount)
