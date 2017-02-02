package main

import (
	"time"
	"database/sql"
	"fmt"
)

type Sales struct {
	ID int `db:"id"`
	FirstName string `db:"firstname"`
	LastName string `db:"lastname"`
}

type Customer struct {
	ID int `db:"id"`
	FirstName string `db:"firstname"`
	LastName string `db:"lastname"`
	Email string `db:"email"`
	ContactNumber string `db:"contact_number"`
}

type Supplier struct {
	ID int `db:"id"`
	Name string `db:"name"`
	ContactNumber string `db:"main_contact_number"`
}

type CustomerAddress struct {
	ID int `db:"id"`
	CustomerID int64 `db:"customer_id"`
	Address1 string `db:"line1"`
	Address2 string `db:"line2"`
	City string `db:"city"`
	Postcode string `db:"postcode"`
}

type Stock struct {
	ID int `db:"id"`
	Title string `db:"title"`
	Description string `db:"description"`
	Price string `db:"price"`
	SupplierID int `db:"supplier_id"`
	SupplierName string
	TotalStock int `db:"total_stock"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Order struct {
	ID int `db:"id"`
	Quantity int `db:"quantity"`
	StatusID int `db:"status_id"`
	StockID int `db:"stock_id"`
	CustomerID int `db:"customer_id"`
	Address1 string `db:"address_line1"`
	Address2 string `db:"address_line2"`
	City string `db:"address_city"`
	Postcode string `db:"address_postcode"`
	StockName string
	CustomerName string
}

type SalesBonus struct {
	ID int `db:"id"`
	SalesName string `db:"sales_name"`
	Bonus float64 `db:"bonus"`
}

func CreateCustomer(customer *Customer) (sql.Result, bool) {
	db := DatabaseContext()

	fmt.Println(customer)
	stmt, err := db.Prepare("INSERT INTO customers SET firstname=?,lastname=?,email=?,contact_number=?")
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	res, err := stmt.Exec(customer.FirstName, customer.LastName, customer.Email, customer.ContactNumber)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	return res, true
}

func UpdateCustomer(customer *Customer, customerId string) bool {
	db := DatabaseContext()

	query := fmt.Sprintf("UPDATE customers SET firstname='%s',lastname='%s',email='%s',contact_number='%s' WHERE id = %s", customer.FirstName, customer.LastName, customer.Email, customer.ContactNumber, customerId)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func GetSalesBonusResults() ([]SalesBonus, bool) {
	db := DatabaseContext()

	rows, err := db.Query("CALL `baltic`.`GetSalesBonuses`()")
	if err != nil {
		return nil, false
	}

	ents := []SalesBonus{}
	for rows.Next() {
		var r SalesBonus
		err = rows.Scan(&r.ID, &r.SalesName, &r.Bonus)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func CreateCustomerAddress(address *CustomerAddress) (sql.Result, bool) {
	db := DatabaseContext()

	stmt, err := db.Prepare("INSERT INTO customer_addresses SET customer_id=?,country_code_id=?,line1=?,line2=?,city=?,postcode=?")
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	res, err := stmt.Exec(address.CustomerID, 1, address.Address1, address.Address2, address.City, address.Postcode)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	return res, true
}

func UpdateCustomerAddress(address *CustomerAddress, customerId string) bool {
	db := DatabaseContext()

	query := fmt.Sprintf("UPDATE customer_addresses SET line1='%s',line2='%s',city='%s',postcode='%s' WHERE customer_id = %s", address.Address1, address.Address2, address.City, address.Postcode, customerId)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func DeleteCustomer(id int) bool {
	db := DatabaseContext()

	custQuery := fmt.Sprintf("DELETE FROM customers WHERE id = %d", id)
	_, err := db.Query(custQuery)
	if err != nil {
		return false
	}

	return true
}

func DeleteSupplier(id int) bool {
	db := DatabaseContext()

	query := fmt.Sprintf("DELETE FROM supplier WHERE id = %d", id)
	_, err := db.Query(query)
	if err != nil {
		return false
	}

	return true
}

func CreateSales(sales *Sales) (sql.Result, bool) {
	db := DatabaseContext()

	stmt, err := db.Prepare("INSERT INTO sales SET firstname=?,lastname=?")
	if err != nil {
		return nil, false
	}

	res, err := stmt.Exec(sales.FirstName, sales.LastName)
	if err != nil {
		return nil, false
	}

	return res, true
}

func DeleteSales(salesId int) bool {
	db := DatabaseContext()

	custQuery := fmt.Sprintf("DELETE FROM sales WHERE id = %d", salesId)
	_, err := db.Query(custQuery)
	if err != nil {
		return false
	}

	return true
}

func GetAllSales() ([]Sales, bool) {
	db := DatabaseContext()

	rows, err := db.Query("SELECT id, firstname, lastname FROM sales")
	if err != nil {
		return nil, false
	}

	ents := []Sales{}
	for rows.Next() {
		var r Sales
		err = rows.Scan(&r.ID, &r.FirstName, &r.LastName)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func GetAllCustomers() ([]Customer, bool) {
	db := DatabaseContext()

	rows, err := db.Query("SELECT id, firstname, lastname, email FROM customers")
	if err != nil {
		return nil, false
	}

	ents := []Customer{}
	for rows.Next() {
		var r Customer
		err = rows.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func GetAllSuppliers() ([]Supplier, bool) {
	db := DatabaseContext()

	rows, err := db.Query("SELECT s.id, s.name, s.main_contact_number FROM baltic.supplier AS s")
	if err != nil {
		return nil, false
	}

	ents := []Supplier{}
	for rows.Next() {
		var r Supplier
		err = rows.Scan(&r.ID, &r.Name, &r.ContactNumber)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func GetCustomerByID (customerId int) (Customer, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id, firstname, lastname, email, contact_number FROM customers WHERE id = ?", customerId)

	var r Customer
	err := row.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email, &r.ContactNumber)
	if err != nil {
		fmt.Printf("Scan: %v", err)
		return Customer{}, false
	}

	return r, true
}

func GetSalesByID (salesId int) (Sales, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id, firstname, lastname FROM sales WHERE id = ?", salesId)

	var r Sales
	err := row.Scan(&r.ID, &r.FirstName, &r.LastName)
	if err != nil {
		fmt.Printf("Scan: %v", err)
		return Sales{}, false
	}

	return r, true
}

func UpdateSales(sales *Sales, salesId string) bool {
	db := DatabaseContext()

	query := fmt.Sprintf("UPDATE sales SET firstname='%s',lastname='%s' WHERE id = %s", sales.FirstName, sales.LastName, salesId)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func GetFirstCustomerAddressByID(customerId int) (CustomerAddress, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id,line1,line2,city,postcode FROM customer_addresses WHERE customer_id = ?", customerId)

	var r CustomerAddress
	err := row.Scan(&r.ID,&r.Address1, &r.Address2,&r.City,&r.Postcode)
	if err != nil {
		fmt.Printf("Scan: %v", err)
		return CustomerAddress{}, false
	}

	return r, true
}

func GetAllStock(supplierId int) ([]Stock, bool) {
	db := DatabaseContext()

	var query string
	if supplierId == 0 {
		query = "SELECT id, title, description, price, supplier_id, total_stock, (SELECT name FROM supplier as sup WHERE sup.id = s.supplier_id) as supplier_name FROM stock as s"
	} else {
		query = fmt.Sprintf("SELECT id, title, description, price, supplier_id, total_stock, (SELECT name FROM supplier as sup WHERE sup.id = s.supplier_id) as supplier_name FROM stock as s WHERE s.supplier_id = %d", supplierId)
	}

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	ents := []Stock{}
	for rows.Next() {
		var r Stock
		err = rows.Scan(&r.ID, &r.Title, &r.Description, &r.Price, &r.SupplierID, &r.TotalStock, &r.SupplierName)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func GetStockByID(stockId int) (Stock, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id, title, description, price, supplier_id, total_stock, (SELECT name FROM supplier as sup WHERE sup.id = s.supplier_id) as supplier_name FROM stock as s WHERE s.id = ?", stockId)

	var r Stock
	err := row.Scan(&r.ID, &r.Title, &r.Description, &r.Price, &r.SupplierID, &r.TotalStock, &r.SupplierName)
	if err != nil {
		fmt.Printf("Scan: %v", err)
		return Stock{}, false
	}

	return r, true
}

func GetOrderByID(orderId int) (Order, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id, quantity, stock_id, status_id, (SELECT title FROM stock WHERE id = o.stock_id) as stock_name, (SELECT CONCAT(c.firstname, ' ', c.lastname) FROM customers as c WHERE id = o.customer_id) as customer_name,address_line1, address_line2, address_city, address_postcode FROM orders as o WHERE o.id = ?", orderId)

	var r Order
	err := row.Scan(&r.ID, &r.Quantity, &r.StockID, &r.StatusID, &r.StockName, &r.CustomerName, &r.Address1, &r.Address2, &r.City, &r.Postcode)
	if err != nil {
		fmt.Printf("Scan: %v", err)
		return Order{}, false
	}

	return r, true
}

func AddStock(stock *Stock) bool {
	db := DatabaseContext()

	stmt, err := db.Prepare("INSERT INTO stock SET title=?,description=?,price=?,supplier_id=?,added_by=?,total_stock=?")
	if err != nil {
		fmt.Println(err)
		return false
	}

	res, err := stmt.Exec(stock.Title, stock.Description, stock.Price, stock.SupplierID, 1, stock.TotalStock)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if _, err := res.RowsAffected(); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func UpdateStockItemQuantity(stockId, quantity int) bool  {
	db := DatabaseContext()

	query := fmt.Sprintf("UPDATE stock SET total_stock=%d WHERE id=%d", quantity, stockId)
	_, err := db.Query(query)
	if err != nil {
		return false
	}

	return true
}

func DeleteStock(stockId int) bool {
	db := DatabaseContext()

	query := fmt.Sprintf("DELETE FROM stock WHERE id = %d", stockId)
	_, err := db.Query(query)
	if err != nil {
		return false
	}

	return true
}

func GetAllOrders() ([]Order, bool) {
	db := DatabaseContext()

	var query string
	query = "SELECT id, quantity, status_id, (SELECT title FROM stock WHERE id = o.stock_id) as stock_name, (SELECT CONCAT(c.firstname, ' ', c.lastname) FROM customers as c WHERE id = o.customer_id) as customer_name,address_line1, address_line2, address_city, address_postcode FROM orders as o"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	ents := []Order{}
	for rows.Next() {
		var r Order
		err = rows.Scan(&r.ID, &r.Quantity, &r.StatusID, &r.StockName, &r.CustomerName, &r.Address1, &r.Address2, &r.City, &r.Postcode)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func CreateOrder(stockId, quantity, customerId, salesId int) (bool, string) {
	stockItem, success := GetStockByID(stockId)
	if !success {
		return false, "Failed to get stock record"
	}

	if stockItem.TotalStock < quantity {
		return false, "Not enough stock"
	}
	
	customerAddr, success := GetFirstCustomerAddressByID(customerId)
	if !success {
		return false, "Unable to gather customer address, please ensure that there is an address added for this customer"
	}

	db := DatabaseContext()

	stmt, err := db.Prepare("INSERT INTO orders SET stock_id=?,customer_id=?,sales_id=?,quantity=?,address_line1=?,address_line2=?,address_city=?,address_postcode=?")
	if err != nil {
		return false, "There was an issue preparing the query to insert into the database"
	}

	res, err := stmt.Exec(stockId, customerId, salesId, quantity, customerAddr.Address1, customerAddr.Address2, customerAddr.City, customerAddr.Postcode)
	if err != nil {
		fmt.Println(err)
		return false, "Failed to create order"
	}

	if _, err := res.RowsAffected(); err != nil {
		fmt.Println(err)
		return false, "Failed to create order"
	}

	stockLeft := stockItem.TotalStock - quantity
	updated := UpdateStockItemQuantity(stockId, stockLeft)
	if !updated {
		return false, "Successfully created record, however, failed to update total stock left for stock item."
	}

	return true, ""
}

func CancelOrder(orderId int) (bool, string) {
	order, success := GetOrderByID(orderId)
	if !success {
		return false, "Unable to find order"
	}

	stockItem, success := GetStockByID(order.StockID)
	if !success {
		return false, "Unabe to get stock information"
	}

	db := DatabaseContext()

	query := fmt.Sprintf("DELETE FROM orders WHERE id = %d", orderId)
	_, err := db.Query(query)
	if err != nil {
		return false, "Failed to cancel order"
	}

	newQuantity := stockItem.TotalStock + order.Quantity
	updated := UpdateStockItemQuantity(stockItem.ID, newQuantity)
	if !updated {
		return false, "Failed to update stock"
	}

	return true, ""
}