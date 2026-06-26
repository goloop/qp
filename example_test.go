package qp_test

import (
	"fmt"
	"net/url"

	"github.com/goloop/qp/v2"
)

// ExampleNew shows the holder reading several parameters from one decoded query
// — the typical pagination/listing case.
func ExampleNew() {
	u, _ := url.Parse("http://x/?page=2&limit=50&sort=name&tag=go,web")

	q := qp.New(u)
	page := q.Int("page", qp.Default(1), qp.Min(1))
	limit := q.Int("limit", qp.Default(20), qp.Between(1, 100))
	sort := q.String("sort", qp.Default("created_at"))
	tags := q.StringSlice("tag")

	fmt.Println(page, limit, sort, tags)
	// Output: 2 50 name [go web]
}

// ExampleQuery_GetString shows validation against an allowed set.
func ExampleQuery_GetString() {
	u, _ := url.Parse("http://x/?role=admin")

	role, ok := qp.New(u).GetString("role", qp.OneOf("admin", "user"))
	fmt.Printf("%q %v\n", role, ok)

	role, ok = qp.New(u).GetString("role", qp.OneOf("guest", "user"))
	fmt.Printf("%q %v\n", role, ok) // not allowed -> zero value, ok=false
	// Output:
	// "admin" true
	// "" false
}

// ExampleParseInt shows the validation contract: an out-of-range value yields
// the default plus a non-nil error.
func ExampleParseInt() {
	u, _ := url.Parse("http://x/?age=99")

	r := qp.ParseInt(u, "age", qp.Default(18), qp.Between(18, 65))
	fmt.Println(r.Value, r.Error != nil)
	// Output: 18 true
}

// ExamplePullBool shows tri-state optional filters: an absent parameter pulls
// nil, so only supplied filters take effect.
func ExamplePullBool() {
	u, _ := url.Parse("http://x/?is_active=true&is_staff=false")

	q := qp.New(u)
	for _, name := range []string{"is_active", "is_staff", "is_superuser"} {
		if v := q.PullBool(name); v != nil {
			fmt.Printf("%s=%t\n", name, *v)
		} else {
			fmt.Printf("%s=<unset>\n", name)
		}
	}
	// Output:
	// is_active=true
	// is_staff=false
	// is_superuser=<unset>
}

// ExampleQuery_IntSlice shows element validation and the unified handling of the
// comma and repeated forms.
func ExampleQuery_IntSlice() {
	u, _ := url.Parse("http://x/?ids=1,2&ids=3")

	ids := qp.New(u).IntSlice("ids", qp.Each(qp.Between(1, 100)))
	fmt.Println(ids)
	// Output: [1 2 3]
}
