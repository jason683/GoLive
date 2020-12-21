package main

import (
	"fmt"
	"testing"
)

func TestLinkedList(t *testing.T) {
	myList := &linkedList{nil, 0}
	myList.addNode("handsanitiser", 30)
	myList.addNode("facemask", 20)
	myList.addNode("gloves", 10)
	if myList.size == 3 {
		fmt.Println("Result is correct")
	} else {
		t.Error("Incorrect linked list size")
	}

	testMapResult := myList.computeTotalCost()
	if testMapResult["totalcost"] == 150+60+20 {
		fmt.Println("Result is correct")
	} else {
		t.Errorf("Incorrect total cost. Want %v", 150+60+20)
	}

	mapItemsAddedToCart := myList.returnAllNodes()
	if mapItemsAddedToCart["handsanitiser"] != 30 {
		t.Error("Incorrect quantity")
	} else {
		fmt.Println("Correct result")
	}

	myList.removeAllNodes()
	if myList.head == nil {
		fmt.Println("Correct result")
	} else {
		t.Error("Incorrect result")
	}

	if myList.size == 0 {
		fmt.Println("Correct result")
	} else {
		t.Errorf("Incorrect size. Want %v", 0)
	}

}
