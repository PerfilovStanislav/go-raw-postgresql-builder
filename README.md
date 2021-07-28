## Usage

```text
go get github.com/PerfilovStanislav/go-raw-postgresql-builder@v1.0.0
```

## EXAMPLES
- [Simple](#simple-example)
- [Bulk insert](#bulk-insert-example)
- [Sql in sql](#sql-in-sql-example)

## Simple Example

```go
package main

import (
	"fmt"
    ps "github.com/PerfilovStanislav/go-raw-postgresql-builder"
)

func main() {
	type User struct {
		Lastname  string
		Firstname string
		IsAuthor  bool
		FakeField int
	}
	user := User{"Perfilov", "Stanislav", true, 123}
	sql := ps.Sql{
		`SELECT * FROM USERS WHERE firstname = $Firstname AND lastname = $Lastname AND is_author = $IsAuthor`,
		user,
	}
	fmt.Println(sql.String()) // SELECT * FROM USERS WHERE firstname = 'Stanislav' AND lastname = 'Perfilov' AND is_author = TRUE
}
```

## Bulk insert example

```go
/* create table post
(
    id int,
    title text,
    rating float,
    keywords text[],
    author_ids int[],
    is_published bool
); */
package main

import (
	"fmt"
    ps "github.com/PerfilovStanislav/go-raw-postgresql-builder"
)

func main() {
	type Post struct {
		Id          int
		Title       string
		Rating      float64
		AuthorIds   []int
		IsPublished bool
	}
	type Catalog struct {
		Posts []Post
	}
	post1 := Post{10, "Perfilov", 9.95, []int{1001, 1002},false}
	post2 := Post{20, "Stanislav", 9.9, []int{1001}, true}
	catalog := Catalog{[]Post{post1, post2}}

	sql := ps.Sql{
		"INSERT INTO post(id, title, rating, author_ids, is_published) VALUES $Values",
		struct{ Values ps.Sql }{
			ps.Sql{Query: "($Id, $Title, $Rating, '{$AuthorIds}', $IsPublished)", Data: catalog.Posts},
		},
	}

	fmt.Println(sql.String()) /* Result
	    INSERT INTO post(id, title, rating, author_ids, is_published) 
	    VALUES 
                (10, 'Perfilov', 9.950000, '{1001,1002}', FALSE),
                (20, 'Stanislav', 9.900000, '{1001}', TRUE) 
	*/

}
```


## Sql in sql example

```go
package main

import (
	"fmt"
    ps "github.com/PerfilovStanislav/go-raw-postgresql-builder"
)

func main() {
	type Params struct {
		Rating      int
		PeriodName  string
	}
	mainSql := `
		WITH _get_post_admins(admin_id) AS (
			$AdminSql
		), _insert_to_stat AS (
			$StatSql
		)
		SELECT id, first_name, last_name 
		FROM admins
		INNER JOIN _get_post_admins 
			ON _get_post_admins.admin_id = admins.id
	`

	params := Params{9, "day"}
	adminSql := ps.Sql{`
			SELECT DISTINCT admin_id 
			FROM posts 
			WHERE rating > $Rating 
			ORDER BY rating DESC 
			LIMiT 10`, params}
	statSql := ps.Sql{`
			INSERT INTO dayly_stats(admin_id, day)
			SELECT admin_id, date_trunc($PeriodName, now())
			FROM _get_post_admins`, params}

	sql := ps.Sql{mainSql,
		struct{
			AdminSql ps.Sql
			StatSql ps.Sql
		}{
			adminSql,
			statSql,
		},
	}

	fmt.Println(sql.String()) /*
		WITH _get_post_admins(admin_id) AS (
			SELECT DISTINCT admin_id
			FROM posts
			WHERE rating > 9
			ORDER BY rating DESC
			LIMiT 10
		), _insert_to_stat AS (
			INSERT INTO dayly_stats(admin_id, day)
			SELECT admin_id, date_trunc('day', now())
			FROM _get_post_admins
		)
		SELECT id, first_name, last_name
		FROM admins
		INNER JOIN _get_post_admins
			ON _get_post_admins.admin_id = admins.id
	*/

}
```

## Creators

**Perfilov Stanislav**

- <https://t.me/PerfilovStanislav>
