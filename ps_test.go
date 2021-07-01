package ps

import "testing"

func TestSimpleStruct(t *testing.T) {
	want := "SELECT * FROM USERS WHERE firstname = 'Stanislav' AND lastname = 'Perfilov' AND is_author = TRUE"
	type User struct {
		Lastname  string
		Firstname string
		IsAuthor  bool
		FakeField int
	}
	user := User{"Perfilov", "Stanislav", true, 123}
	sql := Sql{
		"SELECT * FROM USERS WHERE firstname = $Firstname AND lastname = $Lastname AND is_author = $IsAuthor",
		user,
	}
	if got := sql.String(); got != want {
		t.Errorf("\nGot: %q\nWant:%q", got, want)
	}
}
