package discovery

import (
	"fmt"
	"testing"
)

func TestContainerGet(t *testing.T) {
	c := newContainer()
	c.add("k", "v")

	v := c.get("k")

	if v != "v" {
		t.Errorf("expected 'v' but got: %s", v)
	}
}

func TestContainerRemove(t *testing.T) {
	c := newContainer()
	c.add("name", "fatrbaby")
	t.Log("cursor:", c.cursor, "\n")
	fmt.Printf("keys: %v\n", c.keys)
	fmt.Printf("values: %v\n", c.values)
	fmt.Printf("indexes: %v\n", c.indexes)

	c.remove("name")

	fmt.Printf("keys: %v\n", c.keys)
	fmt.Printf("values: %v\n", c.values)
	fmt.Printf("indexes: %v\n", c.indexes)

	t.Log("cursor:", c.cursor, "\n")

	v := c.get("name")

	if v != "" {
		t.Errorf("expected empty, got: %s", v)
	}
}
