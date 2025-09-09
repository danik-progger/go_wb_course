package main

type Human struct {
	age  int
	name string
}

type Action struct {
	Human
	Active bool
}

func (a *Action) introduce() string {
	return "My name is " + a.name
}

// func main() {
// 	a := Action{
// 		Human: Human{
// 			age:  30,
// 			name: "Alice",
// 		},
// 		Active: true,
// 	}

// 	println(a.introduce())
// }
