package main

import "net/http"

func main() {
	http.HandleFunc("/", GetHomeIndex)
	http.HandleFunc("/customers", GetCustomerIndex)
	http.HandleFunc("/customer/create", PostCreateCustomer)
	http.HandleFunc("/customer/delete", GetCustomerDelete)
	http.HandleFunc("/suppliers", GetSupplierIndex)
	http.HandleFunc("/supplier/delete", GetSupplierDelete)
	http.HandleFunc("/stock", GetStockIndex)
	http.HandleFunc("/stock/create", PostStockCreate)
	http.HandleFunc("/stock/delete", GetStockDelete)
	http.HandleFunc("/orders", GetOrderIndex)
	http.HandleFunc("/order/create", PostOrderCreate)
	http.HandleFunc("/order/cancel", GetOrderCancel)
	http.HandleFunc("/sales", GetSalesIndex)
	http.HandleFunc("/sales/create", PostSalesCreate)
	http.HandleFunc("/sales/delete", GetSalesDelete)
	http.HandleFunc("/sales/edit", GetSalesmenEdit)
	http.HandleFunc("/sales/edit/update", PostSalesEdit)
	http.HandleFunc("/sales/bonus", GetSalesBonus)
	http.HandleFunc("/customer/edit", GetCustomerEdit)
	http.HandleFunc("/customer/edit/update", PostCustomerEdit)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
