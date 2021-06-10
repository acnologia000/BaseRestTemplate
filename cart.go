package main

// map[userID]map[ProductID]<Number of items for given product ID>
var Carts = make(map[string]map[string]int)

func checkForCart(userID string) bool {
	_, isCartPresent := Carts[userID]
	return isCartPresent
}
func checkForCartItem(userID, item string) bool {
	_, isCartItemPresent := Carts[userID][item]
	return isCartItemPresent
}

func addCart(userID string) bool {

	if checkForCart(userID) {
		return false
	}

	nmap := make(map[string]int)
	Carts[userID] = nmap

	return true
}

func checkOutCart(userID string) {
	delete(Carts, userID)
}

func addItem(userID, itemCode string) { // both adds to cart and incriments

	_, CartExists := Carts[userID][itemCode]

	if !CartExists {
		Carts[userID][itemCode] = 1
	} else {
		Carts[userID][itemCode]++
		print("path called\n")
	}
}

func reduceItem(userID, item string) bool {
	if !checkForCartItem(userID, item) {
		return false
	} else if Carts[userID][item] == 1 {
		deleteItem(userID, item)
		return true
	} else {
		Carts[userID][item]--
		return true
	}
}

func deleteItem(userID, ItemCode string) {
	delete(Carts[userID], ItemCode)
}
