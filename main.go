//c:\program files\postgresql\9.3\bin> pg_dump.exe -U postgres gauravpersoninfo >"d:\backup.sql" <-- to export db
//c:\program files\postgresql\9.3\bin> psql -h gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com -p 5432  -U postgres newDBname < D:\backup.sql <-- to import db
//C:\Program Files\PostgreSQL\9.0\bin\pg_dump.exe --host gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com --port 5432 --username postgres --format plain --ignore-version --verbose --file "C:\temp\filename.backup" --table public.tablename dbname <-- to export special table
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	echo "github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type Person struct {
	// gorm.Model
	ID        int     `gorm:"primary_key:true;"`
	Firstname string  `json:"firstname;"`
	Lastname  string  `json:"lastname;"`
	Age       uint    `json:"age"`
	Gender    string  `json:"gendertype;"`
	Address   Address `gorm:"foreignkey:ID;association_foreignkey:ID"`
	Contact   Contact `gorm:"foreignkey:ID;association_foreignkey:ID"`
}

type Address struct {
	// gorm.Model
	ID    uint
	City  string `json:"city;"`
	State string `json:"state;"`
	Pin   string `json:"pin";`
}
type Contact struct {
	// gorm.Model
	ID     uint
	Mobile string `json:"mobile;"`
	Email  string `json:"email;"`
}

func handlerequest() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", helloWorld)
	e.POST("/person", PostPerson)
	e.GET("/person", GetAll)
	e.GET("/person/:ID", GetPerson)
	e.PUT("/person/:ID", UpdatePerson)
	e.DELETE("/person/:ID", DeletePerson)
	e.GET("/person/a/:ID", GetAddress)
	e.GET("/person/f/:ID", GetFullname)
	e.GET("/person/c/:ID", GetContact)
	e.Logger.Fatal(e.Start(":12345"))

}

func checkError(err error) {
	if err != nil {
		log.Panic("Error detected-->", err)
	}
}

func initiatMigrate() {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	db.AutoMigrate(
		&Person{},
		&Address{},
		&Contact{},
	)
}

var Persons []Person
var people Person

// var addresses Address

var adresses []Address
var contacts []Contact

// Handler
func helloWorld(c echo.Context) (err error) {
	return c.String(http.StatusOK, "Hello, World!")
}

// to PostPerson
func PostPerson(c echo.Context) (err error) {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	u := new(Person)
	if err := c.Bind(u); err != nil {
		return err
	}
	db.Create(&u)
	fmt.Println("u->", u)
	return c.JSON(http.StatusOK, &u)
}

// to Getall records
func GetAll(c echo.Context) error {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	var all []Person
	db.Preload("Address").Preload("Contact").Find(&all)
	return c.JSON(http.StatusOK, &all)
}

// to GetPerson by id
func GetPerson(c echo.Context) error {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	e := new(Person)
	id := c.Param("ID")
	db.Where("id = ?", id).Preload("Address").Preload("Contact").Find(&e)
	return c.JSON(http.StatusOK, &e)
}

// to DeletePerson by id
func DeletePerson(c echo.Context) error {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	e := new(Person)
	a := new(Address)
	var g []Contact
	id := c.Param("ID")
	db.Where("id=?", id).Find(&e).Delete(&e)
	db.Where("id=?", id).Find(&a).Delete(&a)
	db.Where("id=?", id).Find(&g).Delete(&g)
	return c.JSON(http.StatusOK, 204)
}

// to UpdatePerson id
func UpdatePerson(c echo.Context) error {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	type UpdatedAddress struct {
		City  string `json:"city;"`
		State string `json:"state;"`
		Pin   string `json:"pin;"`
	}
	type UpdatedContact struct {
		ID     uint   `json:"id;"`
		Mobile string `json:"mobile;"`
		Email  string `json:"email;"`
	}
	type NewPerson struct {
		ID        int            `gorm:"primary_key:true;"`
		Firstname string         `json:"firstname;"`
		Lastname  string         `json:"lastname;"`
		Gender    string         `json:"gender;"`
		Age       uint           `json:"age"`
		Address   UpdatedAddress `json:"address"`
		Contact   UpdatedContact `json:"contact"`
	}
	n := new(NewPerson)
	if err := c.Bind(n); err != nil {
		fmt.Println(err)
		return err
	}

	p := new(Person)
	id, err := strconv.Atoi(c.Param("ID"))
	checkError(err)
	p.ID = id

	p.Firstname = n.Firstname
	p.Lastname = n.Lastname
	p.Age = n.Age
	p.Gender = n.Gender
	p.Address.State = n.Address.State
	p.Address.City = n.Address.City
	p.Address.Pin = n.Address.Pin
	p.Contact.Mobile = n.Contact.Mobile
	p.Contact.Email = n.Contact.Email

	db.Where("id = ?", id).Save(p)

	db.Find(&p)
	return c.JSON(http.StatusOK, &p)
}

// to GetAddress by id
func GetAddress(c echo.Context) error {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	e := new(Person)
	id := c.Param("ID")
	db.Where("id = ?", id).Preload("Address").Preload("Contact").Find(&e)
	// fmt.Println("e->", e)
	// fmt.Println("Address", e.Address.City)
	adrs := e.Address.City + " " + e.Address.State + e.Address.Pin
	return c.JSON(http.StatusOK, adrs)
}

// to GetFullname by id
func GetFullname(c echo.Context) error {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	e := new(Person)
	id := c.Param("ID")
	db.Where("id = ?", id).Preload("Address").Preload("Contact").Find(&e)
	Fullname := e.Firstname + " " + e.Lastname
	return c.JSON(http.StatusOK, Fullname)
}

// to GetContact by id
func GetContact(c echo.Context) error {
	db, err := gorm.Open("postgres", "host=gauravpersoninfo.co5lyabzhwet.us-east-2.rds-preview.amazonaws.com port=5432 user=gauravpersoninfo dbname=gauravpersoninfo password=gauravpersoninfo sslmode=disable")
	checkError(err)
	defer db.Close()

	e := new(Person)
	id := c.Param("ID")
	db.Where("id = ?", id).Preload("Address").Preload("Contact").Find(&e)
	return c.JSON(http.StatusOK, e.Contact)
}

//main func
func main() {
	fmt.Println("hello world")
	initiatMigrate()
	handlerequest()
}
