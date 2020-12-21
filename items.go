package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		var database string
		database = os.Getenv("DATABASE")
		if key[0] == database {
			return true
		}
		return false
	}
	return false
}

//ItemInfo type struct
type ItemInfo struct {
	ProductCode        string `json:"code"`
	Name               string `json:"name"`
	ProductDescription string `json:"description"`
	Quantity           int    `json:"quantity"`
}

var items = make(map[string]ItemInfo)

func createItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {
		var newItem ItemInfo
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &newItem)
			if newItem.ProductCode == "" || newItem.Name == "" || newItem.ProductDescription == "" || newItem.Quantity == 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}
			if _, ok := items[params["itemid"]]; !ok {
				if newItem.ProductCode == params["itemid"] {
					items[params["itemid"]] = newItem
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Item added: " + params["itemid"]))
					query := fmt.Sprintf("INSERT INTO Items VALUES ('%s', '%s', '%s', '%d')", newItem.ProductCode, newItem.Name, newItem.ProductDescription, newItem.Quantity)
					_, err := db.Query(query)
					if err != nil {
						panic(err.Error())
					}
				} else {
					fmt.Println(errors.New("item id has to match with the product code"))
				}
			} else {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("409 - Duplicate item code"))
			}
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply course information in JSON format"))
	}
}

func retrieveItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	results, err := db.Query("SELECT ProductCode, Name, ProductDescription, Quantity FROM items_db.Items WHERE ProductCode = ?", params["itemid"])
	if err != nil {
		panic(err.Error())
	}
	var existingItem ItemInfo
	for results.Next() {
		err := results.Scan(&existingItem.ProductCode, &existingItem.Name, &existingItem.ProductDescription, &existingItem.Quantity)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(existingItem)
}

func allItems(w http.ResponseWriter, r *http.Request) {
	results, err := db.Query("SELECT ProductCode, Name, ProductDescription, Quantity from Items")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var existingItems ItemInfo
		err := results.Scan(&existingItems.ProductCode, &existingItems.Name, &existingItems.ProductDescription, &existingItems.Quantity)
		if err != nil {
			panic(err.Error())
		}
		items[existingItems.ProductCode] = existingItems
	}
	json.NewEncoder(w).Encode(items)
	fmt.Println(items)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-Type") == "application/json" {
		var updatedItemInfo ItemInfo
		var existingItemInfo ItemInfo
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &updatedItemInfo)
			if updatedItemInfo.ProductCode == "" || updatedItemInfo.Name == "" || updatedItemInfo.ProductDescription == "" || updatedItemInfo.Quantity == 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply item information in JSON format"))
				return
			}
			results, err := db.Query("SELECT ProductCode, Name, ProductDescription, Quantity from Items")
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var existingItems ItemInfo
				err := results.Scan(&existingItems.ProductCode, &existingItems.Name, &existingItems.ProductDescription, &existingItems.Quantity)
				if err != nil {
					panic(err.Error())
				}
				items[existingItems.ProductCode] = existingItems
			}
			for k, v := range items {
				if params["itemid"] == k {
					existingItemInfo = v
					for key := range items {
						if updatedItemInfo.ProductCode == key && existingItemInfo.ProductCode != updatedItemInfo.ProductCode {
							w.WriteHeader(http.StatusConflict)
							w.Write([]byte("409 - There will be a duplicate course code if the edit is in this manner"))
							return
						}
					}
					items[params["itemid"]] = updatedItemInfo
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Item updated: " + params["itemid"]))
					query := fmt.Sprintf("UPDATE Items SET ProductCode='%s', Name='%s', ProductDescription='%s', Quantity='%d' WHERE ProductCode='%s'", updatedItemInfo.ProductCode, updatedItemInfo.Name, updatedItemInfo.ProductDescription, updatedItemInfo.Quantity, existingItemInfo.ProductCode)
					_, err := db.Query(query)
					if err != nil {
						panic(err.Error())
					}
					delete(items, existingItemInfo.ProductCode)
					return
				}
			}
			items[params["itemid"]] = updatedItemInfo
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("item added: " + params["itemid"]))
			query := fmt.Sprintf("INSERT INTO Items VALUES ('%s', '%s', '%s', '%d')", updatedItemInfo.ProductCode, updatedItemInfo.Name, updatedItemInfo.ProductDescription, updatedItemInfo.Quantity)
			_, err = db.Query(query)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := fmt.Sprintf("DELETE FROM Items WHERE ProductCode='%s'", params["itemid"])
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}
