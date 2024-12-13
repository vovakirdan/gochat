package main

// todo hash all passwords
type Database struct {
	users map[string]string  // username -> password
	rooms map[string]string  // rooms -> password; if no password room is open
}

func NewDatabase() *Database {
	db := &Database{
		users: make(map[string]string),
		rooms: make(map[string]string),
	}
	return db
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

func (db *Database) AddRoom(room, password string) bool {
	if _, ok := db.rooms[room]; ok {
		return false
	}
	db.rooms[room] = password
	return true
}

func (db *Database) IsPrivateRoom(room string) bool {
	// if password is empty, room is open
	room_password, exists := db.rooms[room]
	return exists && room_password != "" 
}

func (db *Database) IsUserExists(username string) bool {
	_, exists := db.users[username]
	return exists
}
