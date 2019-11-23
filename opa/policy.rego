# POLICY
package play

default message = false

message {
	users[f].name == input.from
	users[t].name == input.to
    
    users[f].friends[_] == input.to
    users[t].friends[_] == input.from
    
    not is_blocked(users[t].blockeds, input.from)
}

is_blocked(friends, user){
	friends[_] == user
}


# DATA

users := [
	{"name": "rashad",    "friends": ["shahriyar", "farid", "bob"], "blockeds":["bob"]},
    {"name": "shahriyar",    "friends": ["rashad"], "blockeds":[]},
    {"name": "farid",    "friends": ["rashad"], "blockeds":[]},
    {"name": "bob",    "friends": ["rashad"], "blockeds":[]},
]

# TEST

test_friends {
	myInput := {"from": "rashad", "to":"bob"}
    message with input as myInput
}