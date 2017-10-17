// This quickly became the messiest part of this whole bot.
// Expect confusing weirdness in regards to error handling.

package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)

// Gets a standard key from the database
func GetKey(key string, user *discordgo.User) (*redis.StringCmd, error) {
	result := client.Get(key)
	_, err := result.Result()
	return result, err
}

// Gets a hash from the database
func GetHashAll(key string, user *discordgo.User) (*redis.StringStringMapCmd, error) {
	result := client.HGetAll(key)
	_, err := result.Result()
	return result, err
}

// Adds new user to the database. Will not check for existing users, WILL overwrite (but not all fields!).
func CreateUser(user *discordgo.User) (*redis.StatusCmd, error) {

	DebugPrint("Adding user to database.")
	newUser := client.HMSet("user:"+user.ID, map[string]interface{}{
		"username": user.Username,
		"isBot":    user.Bot,
		"posts":    1,
	})

	err := newUser.Err()

	if err == nil {
		DebugPrint("New user added to database.\n" + newUser.String())
	}

	return newUser, err
}

// Gets a user from the database; can be told to automatically create a user if it does not exist
// result.Val()["fieldname"] for field as string
// result.String() for entire hash human-readable
func GetUser(user *discordgo.User, createUser bool) (*redis.StringStringMapCmd, error) {
	result := client.HGetAll("user:" + user.ID)
	err := result.Err()

	// guess we have to do this? redis.Nil apparently isn't a thing for hashes
	if result.Val()["username"] == "" {
		if createUser {
			// add to database
			_, err := CreateUser(user)
			if err != nil {
				return nil, err
			}

			// update this
			result = client.HGetAll("user:" + user.ID)
		} else {
			// don't add to database
			DebugPrint("User does not exist; function call said not to create new user.")
			return nil, redis.Nil // it's technically not from redis, but it makes sense.
		}
	}

	// Did we get an error?
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sets values to a user object.
// This will check for an existing user, but it will not create one.
func SetUser(user *discordgo.User, fields map[string]interface{}) (*redis.StatusCmd, error) {

	// make sure we don't create a new one
	_, err := GetUser(user, false)
	if err != nil {
		return nil, err
	}

	// set the user
	userDb := client.HMSet("user:"+user.ID, fields)
	if userDb.Err() == nil {
		DebugPrint("User object was just modified. Changed values:\n" + userDb.String()) // this will only contain the changed values
	}
	return userDb, userDb.Err()
}
