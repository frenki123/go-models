package db

import "fmt"

func CharField(fieldPtr *string, maxLen int) *field {
	v := fmt.Sprintf("VARCHAR(%d)", maxLen)
	return newField(fieldPtr, v)
}
func EmailField(fieldPtr *string) *field {
	return newField(fieldPtr, "VARCHAR(255)")
}
func TextField(fieldPtr *string) *field {
	return newField(fieldPtr, "TEXT")
}
func IntField(fieldPtr *int) *field {
	return newField(fieldPtr, "INT")
}
