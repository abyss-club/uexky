// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package entity

type User struct {
	ID      int     `json:"id"`
	Name    *string `json:"name"`
	Level   int     `json:"level"`
	Friends []*User `json:"friends"`
}
