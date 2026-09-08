package root

import (
	"fmt"
	"github.com/72sevenzy2/in-memory-database/db"
	"testing"
)

func TestDbInt(t *testing.T) {
	b := *db.NewDB()
	b.SetInt("name", 100000)
	b.Del("name")

	b.SetInt("test1", 100000000)
	b.SetInt("test2", 1434242324)

	all, _ := b.GetAllInt()

	// print all
	for k, v := range all {
		fmt.Println(k, v)
	}

	val, ok := b.GetInt("name")
	if !ok {
		fmt.Println("not found", val)
	}
	fmt.Println(val)
}
