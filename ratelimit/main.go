package main

import (
	"time"
)

const (
	brokenRulesLimit = 10
	timeLimit        = 2500 * time.Millisecond
)

var (
	//Name that keybase will use for background plugins
	Name         = "ratelimit"
	userAccounts = make(map[string]*limiter)
)

type limiter struct {
	brokenRuleCount uint64
	penaltyBonus    uint64
	lastContact     time.Time
}

func updateMap() {
	for {
		for user, t := range userAccounts {
			if time.Since(t.lastContact) > timeLimit {
				delete(userAccounts, user)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

//Validate the user hasn't hit their request limit
func Validate(user string) bool {
	if _, ok := userAccounts[user]; !ok {
		l := &limiter{}
		l.lastContact = time.Now()
		l.brokenRuleCount = 0
		userAccounts[user] = l
		return true
	}
	if time.Since(userAccounts[user].lastContact) < timeLimit {
		if userAccounts[user].brokenRuleCount > brokenRulesLimit {
			userAccounts[user].penaltyBonus++
			userAccounts[user].lastContact = time.Now().Add(time.Second * time.Duration(userAccounts[user].penaltyBonus))
		} else {
			userAccounts[user].brokenRuleCount++
			userAccounts[user].lastContact = time.Now()
		}
		return false
	}
	userAccounts[user].lastContact = time.Now()
	return true
}

func init() {
	go updateMap()
}

func main() {}
