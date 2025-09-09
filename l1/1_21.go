package main

// import (
// 	"fmt"
// )

// // Пример часто использующегося в Go адаптера
// // В Go чтобы поле было публичным нужно, чтобы оно имело название с большой буквы
// // В JSON принято использовать с маленькой буквы
// // Используя инструкции мы можем подружить два представления и использовать в
// // каждом месте нужные
// type Delivery struct {
// 	ID       uint   `gorm:"primaryKey" json:"-"`
// 	OrderUID string `gorm:"size:50;index" json:"-"`
// 	Name     string `gorm:"size:100" json:"name"`
// 	Phone    string `gorm:"size:20" json:"phone"`
// 	Zip      string `gorm:"size:20" json:"zip"`
// 	City     string `gorm:"size:50" json:"city"`
// 	Address  string `gorm:"size:200" json:"address"`
// 	Region   string `gorm:"size:50" json:"region"`
// 	Email    string `gorm:"size:100" json:"email"`
// }

// type Dog struct {}

// func (dog *Dog) WoofWoof() {
//     fmt.Println("Гав-Гав. Хозяин, дай пожрать")
// }

// type Cat struct {}

// func (cat *Cat) MeowMeow(isCall bool) {
//     if isCall {
//         fmt.Println("Где моя еда, раб? Ну так уж и быть... Мяу-мяу")
//     }
// }

// // целевой интерфейс - Target
// type AnimalAdapter interface {
//     Reaction()
// }

// // адаптер для собаки
// type DogAdapter struct{
//     *Dog
// }

// // реакция собаки
// func (adapter *DogAdapter) Reaction() {
//     adapter.WoofWoof()
// }

// // конструктор адаптера для собаки
// func NewDogAdapter(dog *Dog) AnimalAdapter {
//     return &DogAdapter{dog}
// }

// // конструктор адаптера для кота
// func NewCatAdapter(cat *Cat) AnimalAdapter {
//     return &CatAdapter{cat}
// }

// // адаптер для кошки
// type CatAdapter struct{
//     *Cat
// }

// // реакция кошки
// func (adapter *CatAdapter) Reaction() {
//     adapter.MeowMeow(true)
// }

// func Say[T AnimalAdapter](a T) {
// 	a.Reaction()
// }

// func main() {
// 	dog := Dog{}
// 	cat := Cat{}

// 	animal1 := NewDogAdapter(&dog)
// 	animal2 := NewCatAdapter(&cat)

// 	Say(animal1)
// 	Say(animal2)
// }
