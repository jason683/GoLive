package main

import (
	"fmt"
	"net/http"
)

const (
	handSanitiserPrice = 5
	facemaskPrice      = 3
	glovesPrice        = 2
)

func viewShoppingCart(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		panic(err.Error())
	}
	if username, ok := mapSessions[myCookie.Value]; ok {
		myList := mapShoppingCart[username]
		mapItemsAddedToCart := myList.returnAllNodes()
		fmt.Println(mapItemsAddedToCart)
		if mapItemsAddedToCart["emptycart"] == 1 {
			tpl.ExecuteTemplate(res, "emptycart.html", mapItemsAddedToCart)
		} else {
			tpl.ExecuteTemplate(res, "shoppingcart.html", mapItemsAddedToCart)
		}
	}
}

func (p *linkedList) returnAllNodes() map[string]int {
	mapItemsAddedToCart := map[string]int{}
	currentNode := p.head
	if currentNode == nil {
		mapItemsAddedToCart["emptycart"] = 1
		return mapItemsAddedToCart
	}
	mapItemsAddedToCart[currentNode.name] = currentNode.quantity
	fmt.Printf("Name: %s, Quantity: %v\n", currentNode.name, currentNode.quantity)
	for currentNode.next != nil {
		currentNode = currentNode.next
		fmt.Printf("Name: %s, Quantity: %v\n", currentNode.name, currentNode.quantity)
		mapItemsAddedToCart[currentNode.name] = currentNode.quantity
	}
	return mapItemsAddedToCart
}

func (p *linkedList) computeTotalCost() map[string]int {
	currentNode := p.head
	var handSanitiserCost int
	var faceMaskCost int
	var glovesCost int
	mapCheckout := map[string]int{}
	if currentNode == nil {
		return nil
	}
	if currentNode.name == "handsanitiser" {
		handSanitiserCost = currentNode.quantity * handSanitiserPrice
		mapCheckout["handsanitiser"] = handSanitiserCost
	} else if currentNode.name == "facemask" {
		faceMaskCost = currentNode.quantity * facemaskPrice
		mapCheckout["facemask"] = faceMaskCost
	} else if currentNode.name == "gloves" {
		glovesCost = currentNode.quantity * glovesPrice
		mapCheckout["gloves"] = glovesCost
	}

	for currentNode.next != nil {
		currentNode = currentNode.next
		if currentNode.name == "handsanitiser" {
			handSanitiserCost = currentNode.quantity * handSanitiserPrice
			mapCheckout["handsanitiser"] = handSanitiserCost
		} else if currentNode.name == "facemask" {
			faceMaskCost = currentNode.quantity * facemaskPrice
			mapCheckout["facemask"] = faceMaskCost
		} else if currentNode.name == "gloves" {
			glovesCost = currentNode.quantity * glovesPrice
			mapCheckout["gloves"] = glovesCost
		}
	}
	mapCheckout["totalcost"] = handSanitiserCost + faceMaskCost + glovesCost
	fmt.Println(mapCheckout)
	return mapCheckout
}
