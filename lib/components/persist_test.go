package components

import (
	"testing"

	"github.com/x-hgg-x/arkanoid-go/lib/binconv"
)

func Test_Persist(t *testing.T) {

	p, err := NewPersist("arkanoid.db", "Pair")
	if err != nil {
		t.Error("error from creating", err)
	}
	defer p.Close()

	if p.Update([]byte("key"), []byte("value")) != nil {
		return
	}
	v, err := p.View([]byte("key"))
	if err != nil {
		return
	}
	t.Logf("[key]: [%s]", string(v))
}

func Test_PersistArray(t *testing.T) {

	p, err := NewPersist("arkanoid.db", "Array")
	if err != nil {
		t.Error("error from creating", err)
	}
	defer p.Close()

	list := [][]byte{
		[]byte("apples"),
		[]byte("balls"),
		[]byte("cats"),
		[]byte("dogs"),
		[]byte("eagles"),
	}

	if p.UpdateList([]byte("letters"), list) != nil {
		return
	}
	d, err := p.ViewList([]byte("letters"))
	if err != nil {
		return
	}
	for i, v := range d {
		t.Logf("data[%d]: %s", i, string(v))
	}
}

func Test_PersistAppend(t *testing.T) {
	p, err := NewPersist("arkanoid.db", "Test_Persist32")
	if err != nil {
		t.Error("error from creating", err)
		return
	}
	defer p.Close()

	key := []byte("scores")

	d, err := p.ViewList(key)
	if err != nil {
		t.Error("error from getting list", err)
		return
	}
	keyAppended := append(key, binconv.Itob(len(d))...)
	if p.Update(keyAppended, []byte("312")) != nil {
		t.Error("error from updating list", err)
		return
	}
	d, err = p.ViewList(key)
	if err != nil {
		t.Error("error from getting list", err)
		return
	}
	for i, v := range d {
		t.Logf("data[%d]: %s", i, string(v))
	}
}

func Test_PersistActual(t *testing.T) {

	p, err := NewPersist("arkanoid.db", "Score")
	if err != nil {
		t.Error("error from creating", err)
	}
	defer p.Close()

	d, err := p.ViewList([]byte("scores"))
	if err != nil {
		t.Error("error from getting list", err)
		return
	}
	t.Log("len d: ", len(d))
	for i, v := range d {
		t.Logf("data[%d]: %s", i, string(v))
	}
}
