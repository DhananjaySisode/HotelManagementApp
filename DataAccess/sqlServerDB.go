package DataAccess

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type CustomerMaster struct {
	CustId       int
	CustName     string
	CustLoginId  string
	CustPassword string
	CustEmail    string
	CustPhone    string
}

type Login struct {
	CustLoginId  string
	CustPassword string
}

type HotelMenu struct {
	MenuId          int
	CustId          int
	MenuName        string
	MenuDescription string
	Price           int
	IsVeg           int
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/hotelmanagementdev")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetAllMenu() ([]HotelMenu, error) {
	var MenuList []HotelMenu
	db := connect()
	defer db.Close()
	result, err := db.Query("select menuId, custId, menuName, menuDescription, price, isVeg from hotelmanagementdev.tbl_customer_menu")
	if err != nil {
		return MenuList, err
	}

	for result.Next() {
		var menu HotelMenu
		errScan := result.Scan(&menu.MenuId, &menu.CustId, &menu.MenuName, &menu.MenuDescription, &menu.Price, &menu.IsVeg)
		if errScan != nil {
			return MenuList, errScan
		}
		MenuList = append(MenuList, menu)
	}
	return MenuList, nil
}

func GetMenuByCustomer(custId int) ([]HotelMenu, error) {
	var MenuList []HotelMenu
	db := connect()
	defer db.Close()
	result, err := db.Query("select menuId, custId, menuName, menuDescription, price, isVeg from hotelmanagementdev.tbl_customer_menu where custId = ?", custId)
	if err != nil {
		return MenuList, err
	}

	for result.Next() {
		var menu HotelMenu
		errScan := result.Scan(&menu.MenuId, &menu.CustId, &menu.MenuName, &menu.MenuDescription, &menu.Price, &menu.IsVeg)
		if errScan != nil {
			return MenuList, errScan
		}
		MenuList = append(MenuList, menu)
	}
	return MenuList, nil
}

func LoginCustomer(cust Login) (CustomerMaster, error) {
	var custDetails CustomerMaster
	db := connect() //db connection is opened
	defer db.Close()
	row := db.QueryRow("SELECT custId, custName, custLoginId, custPassword, custEmail, custPhone FROM hotelmanagementdev.tbl_customer_master WHERE custLoginId = ? and custPassword= ? and isActive=1", cust.CustLoginId, cust.CustPassword)

	err1 := row.Scan(&custDetails.CustId, &custDetails.CustName, &custDetails.CustLoginId, &custDetails.CustPassword, &custDetails.CustEmail, &custDetails.CustPhone)

	if err1 != nil {
		return custDetails, err1
	}
	return custDetails, nil

}

func RegisterCustomer(cust CustomerMaster) (CustomerMaster, error) {
	var custDetails CustomerMaster
	db := connect() //db connection is opened
	defer db.Close()

	fmt.Println(cust)
	stmt, err := db.Prepare("Insert into hotelmanagementdev.tbl_customer_master(custName, custLoginId, custPassword, custEmail, custPhone,isActive)values(?,?,?,?,?,?)")
	if err != nil {
		return custDetails, err
	}
	res, _ := stmt.Exec(cust.CustName, cust.CustLoginId, cust.CustPassword, cust.CustEmail, cust.CustPhone, 1)

	affected, _ := res.RowsAffected()
	if int(affected) == 1 {
		row := db.QueryRow("SELECT custId, custName, custLoginId, custPassword, custEmail, custPhone FROM hotelmanagementdev.tbl_customer_master WHERE custLoginId = ? and custPassword= ? and isActive=1", cust.CustLoginId, cust.CustPassword)
		err1 := row.Scan(&custDetails.CustId, &custDetails.CustName, &custDetails.CustLoginId, &custDetails.CustPassword, &custDetails.CustEmail, &custDetails.CustPhone)
		if err1 != nil {
			return custDetails, err1
		}
	}
	return custDetails, nil
}
