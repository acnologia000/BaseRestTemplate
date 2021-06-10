package main

import (
	"fmt"
	"sync"
	"time"
)

var sessions = make(map[string]map[string]string)

var tex = &sync.RWMutex{}

func addSession(userID string) {
	sessions[userID] = map[string]string{
		"time": time.Now().Format(TimeFormat),
	}
}

func addSessionData(userID, key, value string) bool {
	if _, f := getSession(userID); !f {
		return false
	}
	sessions[userID][key] = value
	return true
}

func deleteSession(sessionID string) {
	delete(sessions, sessionID)
}

func getSession(key string) (map[string]string, bool) {
	data, isValid := sessions[key]
	return data, isValid
}

func sessionExpireService(ticker *time.Ticker) {
	print("service started")
	go func() {
		for {
			select {
			case <-SessionServiceChannel:
				return
			case t := <-ticker.C:
				fmt.Printf("running Cleanup at %s\n", t.Format(TimeFormat))
				ReadSessionsAndExpire()
			}
		}
	}()

}

func ReadSessionsAndExpire() {
	tex.Lock()
	for k, v := range sessions {
		loginTime, _ := time.Parse(TimeFormat, v["time"])
		if time.Since(loginTime) > (time.Hour * 24) {
			delete(sessions, k)
		}
	}
	tex.Unlock()
}
