package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//Node type struct
type Node struct {
	name     string
	quantity int
	next     *Node
}

type linkedList struct {
	head *Node
	size int
}

func (p *linkedList) addNode(name string, quantity int) {
	newNode := &Node{
		name:     name,
		quantity: quantity,
		next:     nil,
	}
	if p.head == nil {
		p.head = newNode
	} else {
		currentNode := p.head
		if currentNode.name == name {
			p.head.quantity = p.head.quantity + quantity
			return
		}
		for currentNode.next != nil {
			currentNode = currentNode.next
			if currentNode.name == name {
				currentNode.quantity = currentNode.quantity + quantity
				return
			}
		}
		currentNode.next = newNode
	}
	p.size++
	return
}

func viewItems(res http.ResponseWriter, req *http.Request) {
	var sliceOfProductCodes []string
	var productCodes string
	results, err := db.Query("SELECT ProductCode from Items")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err := results.Scan(&productCodes)
		if err != nil {
			panic(err.Error())
		}
		sliceOfProductCodes = append(sliceOfProductCodes, productCodes)
	}
	checkItems := map[string]int{}
	var existingItemsQuantity int
	for _, v := range sliceOfProductCodes {
		query := fmt.Sprintf("SELECT Quantity FROM items_db.Items WHERE ProductCode='%s'", v)
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
		checkItems[v] = existingItemsQuantity
		if existingItemsQuantity == 0 {
			checkItems[v+"0"] = 1
		}
	}
	fmt.Println(checkItems)

	if req.Method == http.MethodPost {
		if !alreadyLoggedIn(req) {
			http.Redirect(res, req, "/", http.StatusSeeOther)
		}
		myCookie, err := req.Cookie("myCookie")
		if err != nil {
			panic(err.Error())
		}
		myList := &linkedList{nil, 0}
		if username, ok := mapSessions[myCookie.Value]; ok {
			*myList = mapShoppingCart[username]
			handSanitiser := req.FormValue("handsanitiser")
			faceMask := req.FormValue("facemask")
			gloves := req.FormValue("gloves")
			handSanitiserQty, _ := strconv.Atoi(handSanitiser)
			faceMaskQty, _ := strconv.Atoi(faceMask)
			glovesQty, _ := strconv.Atoi(gloves)
			errorMessage := map[string]string{}
			log.Println("Items added to shopping cart")
			fmt.Println(handSanitiserQty)
			fmt.Println(faceMaskQty)
			fmt.Println(glovesQty)
			for k, v := range checkItems {
				if k == "QR101" {
					if faceMaskQty > v {
						errorMessage["exceedQR101quantity"] = "Unable to add to cart as number of items requested exceeds quantity available"
						tpl.ExecuteTemplate(res, "itemlist.html", errorMessage)
						delete(errorMessage, "exceedQR101quantity")
						return
					}
				} else if k == "QR102" {
					if glovesQty > v {
						errorMessage["exceedQR102quantity"] = "Unable to add to cart as number of items requested exceeds quantity available"
						tpl.ExecuteTemplate(res, "itemlist.html", errorMessage)
						delete(errorMessage, "exceedQR102quantity")
						return
					}
				} else if k == "QR103" {
					if handSanitiserQty > v {
						errorMessage["exceedQR103quantity"] = "Unable to add to cart as number of items requested exceeds quantity available"
						tpl.ExecuteTemplate(res, "itemlist.html", errorMessage)
						delete(errorMessage, "exceedQR103quantity")
						return
					}
				}
			}
			myList.addNode("handsanitiser", handSanitiserQty)
			myList.addNode("facemask", faceMaskQty)
			myList.addNode("gloves", glovesQty)
			log.Println("WIP")
			mapShoppingCart[username] = *myList
			http.Redirect(res, req, "/checkout", http.StatusSeeOther)
		}
	}
	tpl.ExecuteTemplate(res, "itemlist.html", checkItems)
	for k := range checkItems {
		if k == "QR1010" || k == "QR1020" || k == "QR1030" {
			delete(checkItems, k)
		}
	}
}
