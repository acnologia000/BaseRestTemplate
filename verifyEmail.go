package main

import (
	"encoding/base64"
	"sync"
	"time"
)

var (
	UnverifiedEmails  = make(map[string]map[string]string)
	emailServiceMutex = &sync.RWMutex{}
)

func VerifyEmail(key string) (bool, error) {
	data, IsPresent := UnverifiedEmails[key]

	if !IsPresent {
		return false, nil
	}

	if !IsValidPendingVerification(key) {
		return false, nil
	}

	if err := SetEmailVerifiedInDatabase(data["mail"]); err != nil {
		return false, err
	}

	delete(UnverifiedEmails, key)

	return true, nil
}

func GenerateEmailVerificationRequest(email string) (string, error) {
	print("getting unique key")
	gString, err := GetUniqueKey()
	key := base64.RawStdEncoding.EncodeToString(gString)
	print("\nassigning\n")
	nMap := map[string]string{
		"mail": email,
		"time": time.Now().Format(TimeFormat),
	}
	UnverifiedEmails[key] = nMap
	return key, err
}

func SetEmailVerifiedInDatabase(email string) error {
	return nil
}

//-----------------------------------------------------------------------
func GetUniqueKey() ([]byte, error) {
	print("GetUniqueKey() called ")
	gString, err := GenerateRandomBytes(StdKeyLength)
	if err == nil && !IsKeyUnique(string(gString)) {
		return GetUniqueKey()
	}
	return gString, err
}

func IsKeyUnique(key string) bool {
	_, unique := UnverifiedEmails[key]
	return !unique
}

func IsValidPendingVerification(key string) bool {
	data, isPresent := UnverifiedEmails[key]

	if !isPresent {
		return false
	}

	ReqTime, err := time.Parse(TimeFormat, data["time"])

	if err != nil {
		return false
	}

	if time.Since(ReqTime) > (time.Hour * 24) {
		return false
	}

	return true
}

func DeleteInvalidPendingVerifications() {
	emailServiceMutex.Lock()
	for k := range UnverifiedEmails {
		if !IsValidPendingVerification(k) {
			delete(UnverifiedEmails, k)
		}
	}
	emailServiceMutex.Unlock()
}

func DeleteInvalidPendingVerificationsService(t *time.Ticker) {
	go func() {
		for {
			select {
			case <-VerificationServiceChanel:
				return
			case <-t.C:
				DeleteInvalidPendingVerifications()
			}
		}
	}()
}
