package main

import (
	"fmt"
	"log"
	"net/http"
)

func (p *linkedList) checkoutItemQuantities() (map[string]int, error) {
	checkoutQuantities := map[string]int{}
	currentNode := p.head
	if currentNode == nil {
		return nil, nil
	}
	if currentNode.name == "handsanitiser" {
		checkoutQuantities["QR103"] = currentNode.quantity
	} else if currentNode.name == "facemask" {
		checkoutQuantities["QR101"] = currentNode.quantity
	} else if currentNode.name == "gloves" {
		checkoutQuantities["QR102"] = currentNode.quantity
	}
	for currentNode.next != nil {
		currentNode = currentNode.next
		if currentNode.name == "handsanitiser" {
			checkoutQuantities["QR103"] = currentNode.quantity
		} else if currentNode.name == "facemask" {
			checkoutQuantities["QR101"] = currentNode.quantity
		} else if currentNode.name == "gloves" {
			checkoutQuantities["QR102"] = currentNode.quantity
		}
	}
	return checkoutQuantities, nil
}

func (p *linkedList) removeAllNodes() {
	p.head = nil
	p.size = 0
}

func checkout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		panic(err.Error())
	}
	if username, ok := mapSessions[myCookie.Value]; ok {
		myList := mapShoppingCart[username]
		fmt.Println(myList)
		fmt.Println(mapShoppingCart)
		fmt.Println(username)
		mapCheckout := myList.computeTotalCost()
		fmt.Println(mapCheckout)
		if mapCheckout["totalcost"] == 0 {
			tpl.ExecuteTemplate(res, "emptycart.html", mapCheckout)
			return
		}
		if req.Method == http.MethodPost {
			myList := mapShoppingCart[username]
			checkoutQuantities, _ := myList.checkoutItemQuantities()
			lock := make(chan bool, 1)
			channel := make(chan bool)
			go reduceItemQty(checkoutQuantities, lock, channel)
			itemsDeducted := <-channel
			if itemsDeducted == true {
				http.Redirect(res, req, "/successfulcheckout", http.StatusSeeOther)
			} else if itemsDeducted == false {
				errorMessage := map[string]string{}
				errorMessage["notenoughavailableitems"] = "Sorry. There are not enough available items."
				tpl.ExecuteTemplate(res, "checkoutresult.html", errorMessage)
				return
			}
		}
		tpl.ExecuteTemplate(res, "cartcheckout.html", mapCheckout)
	}
}

func successfulCheckout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		panic(err.Error())
	}
	if username, ok := mapSessions[myCookie.Value]; ok {
		myList := mapShoppingCart[username]
		myList.removeAllNodes()
		delete(mapShoppingCart, username)
		mapCheckoutResult := map[string]string{}
		mapCheckoutResult["successfulcheckout"] = "You have successfully bought your items! You may now click on any link below to be redirected."
		tpl.ExecuteTemplate(res, "checkoutresult.html", mapCheckoutResult)
		delete(mapCheckoutResult, "successfulcheckout")
	}
	return
}

func reduceItemQty(checkoutQuantities map[string]int, lock chan bool, channel chan bool) {
	lock <- true
	var existingItemsQuantity int
	for k, v := range checkoutQuantities {
		query := fmt.Sprintf("SELECT Quantity FROM items_db.Items WHERE ProductCode='%s'", k)
		results, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			err := results.Scan(&existingItemsQuantity)
			if err != nil {
				panic(err.Error())
			}
		}
		if v > existingItemsQuantity {
			channel <- false
			close(channel)
			return
		}
	}
	for k, v := range checkoutQuantities {
		query := fmt.Sprintf("UPDATE Items SET Quantity=Quantity - '%d' WHERE ProductCode='%s'", v, k)
		_, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		delete(checkoutQuantities, k)
	}
	channel <- true
	close(channel)
	<-lock
}

func cancelCheckout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		panic(err.Error())
	}
	if username, ok := mapSessions[myCookie.Value]; ok {
		myList := mapShoppingCart[username]
		checkoutQuantities, _ := myList.checkoutItemQuantities()
		for k := range checkoutQuantities {
			delete(checkoutQuantities, k)
		}
		//addItemQty(checkoutQuantities)
		delete(mapShoppingCart, username)
		myList.removeAllNodes()
		mapCheckoutResult := map[string]string{}
		mapCheckoutResult["failedcheckout"] = "You have cancelled your order and removed all items from your shopping cart."
		tpl.ExecuteTemplate(res, "checkoutresult.html", mapCheckoutResult)
		delete(mapCheckoutResult, "failedcheckout")
	}
	log.Println("WIP")
}
