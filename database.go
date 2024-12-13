package main

type Database struct {
	users map[string]string  // username -> password
}

func NewDatabase() *Database {
	return &Database{
		users: make(map[string]string),
	}
}

func (db *Database) Register(username, password string) bool {
	if _, ok := db.users[username]; ok {
		return false
	}
	db.users[username] = password
	return true
}

func (db *Database) ValidateUser(username, password string) bool {
	pass, exists := db.users[username]
	return exists && pass == password
}
