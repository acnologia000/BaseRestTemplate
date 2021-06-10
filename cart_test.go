package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

const (
	SampleUserID = "hr1YYQeWLfcImYGitSoa95OXx4AkrOXl5aXdz1L"
	SampleCartID = "DjkCSbvTaZndVIUbSmLbrJ59zlKHoOWXkqfm89U"
	SampleItem1  = "PORK_SLAB"
	SampleItem2  = "BACON"
)

func printAsJSON(t *testing.T, data interface{}) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	json.NewEncoder(buffer).Encode(data)
	t.Logf("current cart in json\n %s", buffer.String())
}

func TestAddCard(t *testing.T) {
	addCart(SampleUserID)
	if !checkForCart(SampleUserID) {
		t.Error("cart not added")
	}
	t.Logf("%d", len(Carts))
}

func TestCheckout(t *testing.T) {
	addCart(SampleUserID)

	checkOutCart(SampleUserID)

	if checkForCart(SampleUserID) {
		t.Error("removal failed")
	}
	t.Logf("%d", len(Carts))
}

func TestAddItem(t *testing.T) {
	addCart(SampleUserID)
	addItem(SampleUserID, SampleItem1)

	if len(Carts[SampleUserID]) != 1 {
		t.Error("Add Item Failed")
	}

	t.Log("Calling function again to test incriment")

	addItem(SampleUserID, SampleItem1)

	if Carts[SampleUserID][SampleItem1] != 2 {
		t.Error("Incriment feature failed")
	}
	printAsJSON(t, Carts)
}

func TestDeleteItem(t *testing.T) {
	addCart(SampleUserID)
	addCart("some sample string")
	addItem(SampleUserID, SampleItem1)
	addItem(SampleUserID, SampleItem1)
	addItem(SampleUserID, SampleItem2)
	addItem(SampleUserID, SampleItem2)
	printAsJSON(t, Carts)

	t.Log("deleting...")

	deleteItem(SampleUserID, SampleItem1)

	_, exits := Carts[SampleUserID][SampleItem1]

	if exits {
		t.Error("Deletion Failed")
	}

	printAsJSON(t, Carts)
}

func TestReduceItem(t *testing.T) {
	addCart(SampleUserID)
	addItem(SampleUserID, SampleItem1)
	addItem(SampleUserID, SampleItem1)
	addItem(SampleUserID, SampleItem1)

	var count = Carts[SampleUserID][SampleItem1]

	reduceItem(SampleUserID, SampleItem1)

	if count == Carts[SampleUserID][SampleItem1] {
		t.Error("reduction failed")
	}

	reduceItem(SampleUserID, SampleItem1)
	reduceItem(SampleUserID, SampleItem1)

	if checkForCartItem(SampleUserID, SampleItem1) {
		t.Error("deletion failed")
	}

	printAsJSON(t, Carts)
}
