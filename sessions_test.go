package main

import "testing"

func TestAddSession(t *testing.T) {
	addSession(SampleUserID)
	_, exists := sessions[SampleUserID]
	if !exists {
		t.Error("failed at adding session")
	}
}

func TestGetSessionAndAddData(t *testing.T) {
	addSession(SampleUserID)

	session, isValid := getSession(SampleUserID)

	if !isValid {
		t.Error("adding session failed")
	}

	addSessionData(SampleUserID, "name", "name")

	name, Present := session["name"]
	if !Present {
		t.Error("name is not present")
	} else {
		t.Log(name)
	}

	addSessionData(SampleUserID, "userID", SampleUserID)

	id, Present := session["userID"]
	if !Present {
		t.Error("userID is not present")
	} else {
		t.Log(id)
	}

	time, Present := session["time"]
	if !Present {
		t.Error("name is not present")
	} else {
		t.Log(time)
	}
}
