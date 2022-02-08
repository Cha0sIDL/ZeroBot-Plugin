package jx3

import (
	"fmt"
)

type mental struct {
	ID          uint64 `db:"mentalID"`
	Name        string `db:"mentalName"`
	MentalIcon  string `db:"mentalIcon"`
	Accept      string `db:"acceptName"`
	MentalColor string `db:"mentalColor"`
	Works       int    `db:"works"`
	Relation    int    `db:"relation"`
}

func getMental(mentalName string) string {
	db.Open()
	var mental mental
	arg := fmt.Sprintf("WHERE acceptName LIKE '%%%s%%' OR mentalName='%s'", mentalName, mentalName)
	db.Find("ns_mental", &mental, arg)
	defer db.Close()
	return mental.Name
}
