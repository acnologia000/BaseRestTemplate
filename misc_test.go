package main

import (
	"testing"
)

func TestSnapshot(t *testing.T) {
	err := snapShotCarts()
	if err != nil {
		t.Errorf("CART SNAPSHOT FAILED \n %s ", err.Error())
	}
}
