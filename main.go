package main

import "net/http"

func main() {
	http.HandleFunc("/", GetHomeIndex)
	http.HandleFunc("/users", GetUsersIndex)
	http.HandleFunc("/customers", GetCustomerIndex)
	http.HandleFunc("/customer/create", PostCreateCustomer)
	http.HandleFunc("/customer/delete", GetCustomerDelete)
	http.HandleFunc("/suppliers", GetSupplierIndex)
	http.HandleFunc("/supplier/delete", GetSupplierDelete)
	http.HandleFunc("/supplier/warehouses", GetSupplierWarehousesIndex)
	http.HandleFunc("/supplier/warehouse/delete", GetSupplierWarehouseDelete)
	http.HandleFunc("/stock", GetStockIndex)
	http.HandleFunc("/stock/create", PostStockCreate)
	http.HandleFunc("/stock/delete", GetStockDelete)
	http.HandleFunc("/orders", GetOrderIndex)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
