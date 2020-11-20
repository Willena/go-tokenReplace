package main

import (
	"fmt"
)

func main() {

	fromVal := "4856d&"

	sqlQuery := `SELECT $fields, $func FROM $table WHERE $cond1 LIMIT $value OFFSET $int $extra`
	c := CreateCompound(sqlQuery).WithFailures().
		Put("fields",
			CreateCompound("$f1,$f2,$f3").
				PutString("f1", "Name").
				PutString("f2", "Age").
				PutString("f3", "Address")).
		Put("func", CreateSanitizedCompound("$f($c)").
			PutRaw("f", "mean").
			PutRaw("c", "Age")).
		PutString("table",fromVal).
		Put("cond1", CreateCompound("$n=$_n").
			PutString("n", "Name").
			PutString("_n", "Guillaume d'alambert \"Le barons\"")).
		PutString("value", "156").
		PutString("int", "1455").
		PutString("extra", "This is an extra with unicode ☼")
	sqlCmd, err := c.GetReplacementValue()
	if err != nil{
		fmt.Println("ERROR:", err)
	}

	fmt.Println("SQL:", sqlCmd)

}
