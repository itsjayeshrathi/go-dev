package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func modifyPersonName(person *Person, name string) {
	person.Name = name
}

func (p *Person) modifyPersonAge(age int) {
	p.Age = age
}

func main() {

	person := Person{"Jayesh", 23}
	fmt.Printf("This is our person %+v\n", person)

	modifyPersonName(&person, "Melkey")
	person.modifyPersonAge(24)
	fmt.Printf("This is our person %+v\n", person)
	//anonymous structs

	employee := struct {
		name string
		id   int
	}{
		name: "Alice",
		id:   23,
	}
	fmt.Printf("This is our employee %+v\n", employee)

	type Address struct {
		Street string
		City   string
	}

	type Contact struct {
		Name    string
		Address Address
		Phone   string
	}

	contact := Contact{
		Name:    "Mark",
		Address: Address{Street: "street 1", City: "Peru"},
		Phone:   "+22323232",
	}

	fmt.Printf("this is the first contract %+v\n", contact)

	// var name = "string"
	// fmt.Printf("Your name is %s\n", name)
	// fmt.Println("Hello World")

	// age := 29
	// fmt.Printf("Your age is %d\n", age)

	// var city string
	// city = "Seattle"
	// fmt.Printf("Your city is %s\n", city)

	// var country, continent = "USA", "North America"
	// fmt.Printf("this is my country %s, and this is my continent is %s\n", country, continent)

	// if age >= 18 {
	// 	fmt.Println("you are able to drive through man.")
	// }

	// day := "Tuesday"

	// switch day {
	// case "Monday":
	// 	fmt.Println("Start of the week")
	// case "Tuesday", "Wednesday", "Thursday":
	// 	fmt.Println("Midweek")
	// case "Friday":
	// 	fmt.Println("TGIF")
	// default:
	// 	fmt.Println("its the weekend")
	// }

	// for i := 0; i <= 10; i++ {
	// 	fmt.Println(i)
	// }

	// numbers := [5]int{1, 2, 3, 4, 5}

	// //numbersAtInit := [...]int{1, 2, 3, 4, 5}

	// fmt.Printf("this is our array %v\n", numbers)

	// //Slices
	// fruits := []string{"apple", "banana", "strawberry"}

	// fmt.Printf("these are my fruits %v\n", fruits)

	// fruits = append(fruits, "kiwi")
	// fmt.Printf("these are my fruits with kiwi appened %v\n", fruits)

	// for index, value := range fruits {
	// 	fmt.Printf("the index is %d and the value is %s\n", index, value)
	// }

	// capitalCity := map[string]string{
	// 	"USA":   "Washington DC",
	// 	"India": "New Delhi",
	// 	"UK":    "London",
	// }

	// fmt.Printf("the capital of usa is %v\n", capitalCity["USA"])

	// capital, exists := capitalCity["Germany"]

	// if exists {
	// 	fmt.Println("this is the capital ", capital)
	// } else {
	// 	fmt.Println("does not exist")
	// }

}

func Add(a, b int) int {
	return a + b
}
