package valueobject

import "fmt"

type Address struct {
	Street     string
	City       string
	Province   string
	PostalCode string
}

func NewAddress(street, city, province, postalCode string) (Address, error) {
	if street == "" || city == "" || province == "" || postalCode == "" {
		return Address{}, fmt.Errorf("address fields can not be empty")
	}

	if len(postalCode) != 5 {
		return Address{}, fmt.Errorf("invalid postal code")
	}

	return Address{
		Street:     street,
		City:       city,
		Province:   province,
		PostalCode: postalCode,
	}, nil
}

func (a Address) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", a.Street, a.City, a.Province, a.PostalCode)
}
