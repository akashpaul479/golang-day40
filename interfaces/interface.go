package interfaces

import "fmt"

// user storage defines what can be done
type UserStorage interface {
	SaveUser(name string) error
	GetUser() (string, error)
}

// filestorage is one implementation
type Filestorage struct {
	username string
}

func (f *Filestorage) SaveUser(name string) error {
	f.username = name
	fmt.Println("User saved in file storage")
	return nil

}
func (f *Filestorage) GetUser() (string, error) {
	return f.username, nil
}

// Memorystorage is another implementation
type MemoryStorage struct {
	username string
}

func (m *MemoryStorage) SaveUser(name string) error {
	m.username = name
	fmt.Println("user saved in memory ")
	return nil
}
func (m *MemoryStorage) GetUser() (string, error) {
	return m.username, nil
}

type Userservice struct {
	storage UserStorage
}

func (u Userservice) RegisterUser(name string) {
	err := u.storage.SaveUser(name)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	user, _ := u.storage.GetUser()
	fmt.Println("Fetch user:", user)
}

// main func
func Interfaces() {
	fmt.Println("-----Using File storage-----")
	file := &Filestorage{}
	Service1 := Userservice{storage: file}
	Service1.RegisterUser("Akash")

	fmt.Println("-----Using Memory storage-----")
	memory := &MemoryStorage{}
	service2 := Userservice{storage: memory}
	service2.RegisterUser("Biswas")
}
