package main

import (
	"time"
	"database/sql"
	"fmt"
	"crypto/sha1"
	"io"
)

type User struct {
	ID int `db:"id"`
	FirstName string `db:"firstname"`
	LastName string `db:"lastname"`
	Email string `db:"email"`
	Password string `db:"password"`
	Active bool `db:"active"`
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
	TotalWarehouses int
}

type SupplierWarehouses struct {
	ID int `db:"id"`
	SupplierName string
	CountryCodeID int `db:"country_code_id"`
	SupplierID int `db:"supplier_id"`
	Address1 string `db:"address_line1"`
	Address2 string `db:"address_line2"`
	City string `db:"city"`
	PostCode string `db:"postcode"`
	ContactNumber string `db:"contact_number"`
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

	success := DeleteSupplierWarehouseBySupplierID(id)
	if !success {
		return false
	}

	query := fmt.Sprintf("DELETE FROM supplier WHERE id = %d", id)
	_, err := db.Query(query)
	if err != nil {
		return false
	}

	return true
}

func DeleteSupplierWarehouse(id int) bool {
	db := DatabaseContext()

	query := fmt.Sprintf("DELETE FROM supplier_warehouses WHERE id = %d", id)
	_, err := db.Query(query)
	if err != nil {
		return false
	}

	return true
}

func DeleteSupplierWarehouseBySupplierID(id int) bool {
	db := DatabaseContext()

	query := fmt.Sprintf("DELETE FROM supplier_warehouses WHERE supplier_id = %d", id)
	_, err := db.Query(query)
	if err != nil {
		return false
	}

	return true
}

func CreateUser(user *User) (sql.Result, bool) {
	db := DatabaseContext()

	stmt, err := db.Prepare("INSERT INTO users SET firstname=?,lastname=?,email=?,password=?,created_at=?,active=?")
	if err != nil {
		return nil, false
	}

	h := sha1.New()
	io.WriteString(h, user.Password)

	res, err := stmt.Exec(user.FirstName, user.LastName, user.Email, h.Sum(nil), time.Now(), user.Active)
	if err != nil {
		return nil, false
	}

	return res, true
}

func GetAllUsers() ([]User, bool) {
	db := DatabaseContext()

	rows, err := db.Query("SELECT id, firstname, lastname, email, active FROM users")
	if err != nil {
		return nil, false
	}

	ents := []User{}
	for rows.Next() {
		var r User
		err = rows.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email, &r.Active)
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

	rows, err := db.Query("SELECT s.id, s.name, s.main_contact_number, (SELECT count(*) FROM baltic.supplier_warehouses as sw WHERE sw.supplier_id = s.id) as total_warehouses FROM baltic.supplier AS s")
	if err != nil {
		return nil, false
	}

	ents := []Supplier{}
	for rows.Next() {
		var r Supplier
		err = rows.Scan(&r.ID, &r.Name, &r.ContactNumber, &r.TotalWarehouses)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func GetAllSuppliersWarehouses(supplierId int) ([]SupplierWarehouses, bool) {
	db := DatabaseContext()
	query := fmt.Sprintf("SELECT sw.id, sw.address_line1, sw.address_line2, sw.city, sw.postcode, sw.contact_number FROM baltic.supplier_warehouses as sw WHERE supplier_id = %d", supplierId)
	rows, err := db.Query(query)
	if err != nil {
		return nil, false
	}

	ents := []SupplierWarehouses{}
	for rows.Next() {
		var r SupplierWarehouses
		err = rows.Scan(&r.ID, &r.Address1, &r.Address2, &r.City, &r.PostCode, &r.ContactNumber)
		if err != nil {
			fmt.Printf("Scan: %v", err)
			return nil, false
		}

		ents = append(ents, r)
	}

	return ents, true
}

func GetUserByEmail(email string) (User, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id, firstname, lastname, email, password, active FROM users WHERE email = ?", email)

	var r User
	err := row.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email, &r.Password, &r.Active)
	if err != nil {
		fmt.Printf("Scan: %v", err)
		return User{}, false
	}

	return r, true
}

func GetCustomerByID (customerId int) (Customer, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id, firstname, lastname, email FROM customers WHERE id = ?", customerId)

	var r Customer
	err := row.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email)
	if err != nil {
		fmt.Printf("Scan: %v", err)
		return Customer{}, false
	}

	return r, true
}

func GetFirstCustomerAddressByID(customerId int) (CustomerAddress, bool) {
	db := DatabaseContext()

	row := db.QueryRow("SELECT id,line1,line2,city,postcode FROM customers WHERE customer_id = ?", customerId)

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

func CreateOrder(stockId, quantity, customerId int) (bool, string) {
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

	stmt, err := db.Prepare("INSERT INTO orders SET stock_id=?,customer_id=?,quantity=?,address_line1=?,address_line2=?,city=?,postcode=?")
	if err != nil {
		return false, "There was an issue preparing the query to insert into the database"
	}

	res, err := stmt.Exec(stockId, customerId, quantity, customerAddr.Address1, customerAddr.Address2, customerAddr.City, customerAddr.Postcode)
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