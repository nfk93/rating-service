// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package database

import ()

type Game struct {
	ID   string
	Name string
}

type Rating struct {
	UserID string
	GameID string
	Rating int32
}

type User struct {
	ID   string
	Name string
}
