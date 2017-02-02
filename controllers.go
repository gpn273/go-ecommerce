package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	"time"
	"os"
)

type Page struct {
	Title string
	Path string
	Error string
	Data interface{}
}

type StockViewModel struct {
	Stock []Stock
	Suppliers []Supplier
}

type OrderViewModel struct {
	Orders []Order
	Stocks []Stock
	Customers []Customer
	Sales []Sales
}

type CustomerEditViewModel struct {
	Customer Customer
	Address CustomerAddress
}

func LoadView(viewName string) *Page {
	filePath := fmt.Sprintf("views/%s.html", viewName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Errorf("File %s does not exist", filePath)
	}

	return &Page{Title: viewName, Path: filePath}
}

func (p *Page) Render(w http.ResponseWriter, data interface{}) {
	t, _ := template.ParseFiles(p.Path)
	p.Data = data
	t.Execute(w, p)
}

func GetHomeIndex(w http.ResponseWriter, r *http.Request) {
	p := LoadView("home")
	p.Render(w, nil)
}

func GetCustomerIndex(w http.ResponseWriter, r *http.Request) {
	p := LoadView("customer_index")
	customers, success := GetAllCustomers()
	if !success {
		p.Error = "Issue gathering user information"
		p.Render(w, nil)
	} else {
		p.Render(w, customers)
	}
}

func GetSupplierIndex(w http.ResponseWriter, r *http.Request) {
	p := LoadView("supplier_index")
	suppliers, success := GetAllSuppliers()
	if !success {
		p.Error = "Issue gathering suppliers information"
		p.Render(w, nil)
	} else {
		p.Render(w, suppliers)
	}
}

func PostCreateCustomer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	customer := &Customer{
		FirstName: r.PostFormValue("firstname"),
		LastName: r.PostFormValue("lastname"),
		Email: r.PostFormValue("email"),
		ContactNumber: r.PostFormValue("contactnumber"),
	}

	cusRes, success := CreateCustomer(customer)
	if !success {
		p := LoadView("customer_index")
		p.Error = "Unable to create customer, please check that all fields are filled out correctly."
		p.Render(w, nil)
		return
	}

	rowId, err := cusRes.LastInsertId()
	if err != nil {
		panic(err)
	}

	address := &CustomerAddress{
		Address1: r.PostFormValue("address1"),
		Address2: r.PostFormValue("address2"),
		City: r.PostFormValue("city"),
		Postcode: r.PostFormValue("postcode"),
		CustomerID: rowId,
	}

	cusAddrRes, success := CreateCustomerAddress(address)
	if !success {
		p := LoadView("customer_index")
		p.Error = "Unable to create customer address, please check that all fields are filled out correctly."
		p.Render(w, nil)
		return
	}

	if _, err := cusAddrRes.RowsAffected(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/customers", 301)
}

func GetCustomerDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	cusId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	DeleteCustomer(cusId)
	http.Redirect(w, r, "/customers", 301)
}

func GetSupplierDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	cusId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	DeleteSupplier(cusId)
	http.Redirect(w, r, "/suppliers", 301)
}

func GetStockIndex(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("supplier_id")
	var supplierId int = 0
	var err error

	if len(id) > 0 {
		supplierId, err = strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
	}

	viewModel := &StockViewModel{}
	stock, success := GetAllStock(supplierId)
	if !success {
		panic("There was an issue gathering stock information")
	}

	viewModel.Stock = stock

	suppliers, success := GetAllSuppliers()
	if !success {
		panic("There was na issue gathering supplier information")
	}

	viewModel.Suppliers = suppliers
	p := LoadView("stock_index")
	p.Render(w, viewModel)
}

func PostStockCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	stock := &Stock{}
	stock.Title = r.PostFormValue("title")
	stock.Description = r.PostFormValue("description")
	stock.Price = r.PostFormValue("price")

	totalStock, err := strconv.Atoi(r.PostFormValue("totalstock"))
	if err != nil {
		p := LoadView("stock_index")
		p.Error = "Please add total stock value"
		p.Render(w, nil)
		return
	}

	supplierId, err := strconv.Atoi(r.PostFormValue("supplier"))
	if err != nil {
		p := LoadView("stock_index")
		p.Error = "Please select a supplier"
		p.Render(w, nil)
		return
	}

	stock.TotalStock = totalStock
	stock.SupplierID = supplierId
	stock.CreatedAt = time.Now()
	stock.UpdatedAt = time.Now()

	created := AddStock(stock)
	if !created {
		p := LoadView("stock_index")
		p.Error = "There was an issue creating the stock, please ensure you have filled out all fields"
		p.Render(w, nil)
	}

	http.Redirect(w, r, "/stock", 301)
}

func GetStockDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	stockId, err := strconv.Atoi(id)
	if err != nil {
		p := LoadView("stock_index")
		p.Error = "Unable to gather stock ID"
		p.Render(w, nil)
		return
	}

	success := DeleteStock(stockId)
	if !success {
		p := LoadView("stock_index")
		p.Error = "There was error deleting the stock record"
		p.Render(w, nil)
		return
	}

	http.Redirect(w, r, "/stock", 301)
}

func GetOrderIndex(w http.ResponseWriter, r *http.Request) {
	orders, success := GetAllOrders()
	if !success {
		p := LoadView("home")
		p.Error = "Unable to gather order information"
		p.Render(w, nil)
		return
	}

	stocks, success := GetAllStock(0)
	if !success {
		p := LoadView("home")
		p.Error = "Unable to gather stock information"
		p.Render(w, nil)
		return
	}

	customers, success := GetAllCustomers()
	if !success {
		p := LoadView("home")
		p.Error = "Unable to gather order information"
		p.Render(w, nil)
		return
	}

	sales, success := GetAllSales()
	if !success {
		p := LoadView("home")
		p.Error = "Unable to gather salesmen information"
		p.Render(w, nil)
		return
	}

	ordersViewModel := &OrderViewModel{
		Orders: orders,
		Stocks: stocks,
		Customers: customers,
		Sales: sales,
	}

	p := LoadView("order_index")
	p.Render(w, ordersViewModel)
}

func PostOrderCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	stockId, err := strconv.Atoi(r.PostFormValue("stockid"))
	if err != nil {
		p := LoadView("order_index")
		p.Error = "Unable to gather stock ID"
		p.Render(w, nil)
		return
	}

	quantity, err := strconv.Atoi(r.PostFormValue("quantity"))
	if err != nil {
		p := LoadView("order_index")
		p.Error = "Unable to gather quanity of item"
		p.Render(w, nil)
		return
	}

	customerId, err := strconv.Atoi(r.PostFormValue("customerid"))
	if err != nil {
		p := LoadView("order_index")
		p.Error = "Unable to gather customer ID"
		p.Render(w, nil)
		return
	}

	salesId, err := strconv.Atoi(r.PostFormValue("salesid"))
	if err != nil {
		p := LoadView("home")
		p.Error = "Unable to gather customer ID"
		p.Render(w, nil)
		return
	}

	success, message := CreateOrder(stockId, quantity, customerId, salesId)
	if !success {
		p := LoadView("home")
		p.Error = message
		p.Render(w, nil)
		return
	}

	http.Redirect(w, r, "/orders", 301)
}

func GetOrderCancel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		p := LoadView("home")
		p.Error = "Unable to gather order ID"
		p.Render(w, nil)
		return
	}

	orderId, err := strconv.Atoi(id)
	if err != nil {
		p := LoadView("home")
		p.Error = "Unable to gather order ID"
		p.Render(w, nil)
		return
	}

	success, message := CancelOrder(orderId)
	if !success {
		p := LoadView("home")
		p.Error = message
		p.Render(w, nil)
		return
	}

	http.Redirect(w, r, "/orders", 301)
}

func GetSalesIndex(w http.ResponseWriter, r *http.Request) {
	sales, success := GetAllSales()
	if !success {
		p := LoadView("home")
		p.Error = "There was an internal error whilst trying to gather salesmen information"
		p.Render(w, nil)
		return
	}

	p := LoadView("sales_index")
	p.Render(w, sales)
}

func PostSalesCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sales := &Sales{
		FirstName: r.PostFormValue("firstname"),
		LastName: r.PostFormValue("lastname"),
	}

	res, success := CreateSales(sales)
	if !success {
		p := LoadView("home")
		p.Error = "There was an internal error whilst trying to create salesmen"
		p.Render(w, nil)
		return
	}

	if _, err := res.LastInsertId(); err != nil {
		p := LoadView("home")
		p.Error = "There was an internal error whilst trying to create salesmen"
		p.Render(w, nil)
		return
	}

	http.Redirect(w, r, "/sales", 301)
}

func GetSalesDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		p := LoadView("home")
		p.Error = "Unable to gather salesmen ID"
		p.Render(w, nil)
		return
	}

	salesId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	success := DeleteSales(salesId)
	if !success {
		p := LoadView("home")
		p.Error = "Failed to delete salesmen"
		p.Render(w, nil)
		return
	}

	http.Redirect(w, r, "/sales", 301)
}

func GetSalesBonus(w http.ResponseWriter, r *http.Request) {
	sb, success := GetSalesBonusResults()
	if !success {
		p := LoadView("home")
		p.Error = "Failed to get sales bonus information"
		p.Render(w, nil)
		return
	}

	m := make(map[string]float64)

	for _, v := range sb {
		if val, ok := m[v.SalesName]; ok {
			m[v.SalesName] = val + v.Bonus
		} else {
			m[v.SalesName] = v.Bonus
		}
	}

	p := LoadView("sales_bonus")
	p.Render(w, m)
}

func GetCustomerEdit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	customerId, err := strconv.Atoi(id)
	if err != nil {
		p := LoadView("home")
		p.Error = "Unable to get customerID"
		p.Render(w, nil)
		return
	}

	customer, success := GetCustomerByID(customerId)
	if !success {
		p := LoadView("home")
		p.Error = "Unable to find customer"
		p.Render(w, nil)
		return
	}

	address, success := GetFirstCustomerAddressByID(customerId)
	if !success {
		p := LoadView("home")
		p.Error = "Unable to find customer address record"
		p.Render(w, nil)
		return
	}

	viewModel := &CustomerEditViewModel{
		Customer: customer,
		Address: address,
	}

	fmt.Println(viewModel)

	p := LoadView("customer_edit")
	p.Render(w, viewModel)
}

func PostCustomerEdit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	customer := &Customer{
		FirstName: r.PostFormValue("firstname"),
		LastName: r.PostFormValue("lastname"),
		Email: r.PostFormValue("email"),
		ContactNumber: r.PostFormValue("contactnumber"),
	}

	success := UpdateCustomer(customer, r.PostFormValue("id"))
	if !success {
		p := LoadView("home")
		p.Error = "Unable to update customer, please check that all fields are filled out correctly."
		p.Render(w, nil)
		return
	}

	address := &CustomerAddress{
		Address1: r.PostFormValue("address1"),
		Address2: r.PostFormValue("address2"),
		City: r.PostFormValue("city"),
		Postcode: r.PostFormValue("postcode"),
	}

	success2 := UpdateCustomerAddress(address, r.PostFormValue("id"))
	if !success2 {
		p := LoadView("home")
		p.Error = "Unable to update customer address, please check that all fields are filled out correctly."
		p.Render(w, nil)
		return
	}

	http.Redirect(w, r, "/customers", 301)
}

func GetSalesmenEdit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	salesId, err := strconv.Atoi(id)
	if err != nil {
		p := LoadView("home")
		p.Error = "Unable to get salesID"
		p.Render(w, nil)
		return
	}

	sales, success := GetSalesByID(salesId)
	if !success {
		p := LoadView("home")
		p.Error = "Unable to find sales record"
		p.Render(w, nil)
		return
	}

	p := LoadView("sales_edit")
	p.Render(w, sales)
}

func PostSalesEdit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sales := &Sales{
		FirstName: r.PostFormValue("firstname"),
		LastName: r.PostFormValue("lastname"),
	}

	success := UpdateSales(sales, r.PostFormValue("id"))
	if !success {
		p := LoadView("home")
		p.Error = "Unable to update sales record"
		p.Render(w, nil)
		return
	}

	http.Redirect(w, r, "/sales", 301)
}