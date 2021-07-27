package ps

import (
	"fmt"
	"testing"
)

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

/* create table post
(
    id int,
    title text,
    rating float,
    keywords text[],
    author_ids int[],
    data json,
    is_published bool
); */

func TestBulkInsert(t *testing.T) {
	want := "INSERT INTO post(id, title, rating, author_ids, data, is_published) VALUES " +
		"(10, 'Perfilov', 9.950000, '{1001,1002}', '{\"Result\":\"is the best\"}', FALSE)," +
		"(20, 'Stanislav', 9.900000, '{1001}', '{\"Result\":\"success\"}', TRUE)"

	type PostData struct {
		Result string
	}
	type Post struct {
		Id          int
		Title       string
		Rating      float64
		AuthorIds   []int
		Data        PostData
		IsPublished bool
	}
	type Catalog struct {
		Posts []Post
	}
	post1 := Post{10, "Perfilov", 9.95, []int{1001, 1002},
		PostData{"is the best"}, false,
	}
	post2 := Post{20, "Stanislav", 9.9, []int{1001},
		PostData{"success"}, true,
	}
	catalog := Catalog{[]Post{post1, post2}}

	posts := Sql{Query: "($Id, $Title, $Rating, '{$AuthorIds}', $Data, $IsPublished)", Data: catalog.Posts}

	sql := Sql{
		"INSERT INTO post(id, title, rating, author_ids, data, is_published) VALUES $Values",
		struct{ Values Sql }{posts},
	}

	if got := sql.String(); got != want {
		t.Errorf("\nGot: %q\nWant:%q", got, want)
	}

	fmt.Println(sql)
}
