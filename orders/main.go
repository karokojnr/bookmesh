package main

func main() {

	store := NewStorage()
	svc := NewService(store)

	svc.CreateOrder()

}
