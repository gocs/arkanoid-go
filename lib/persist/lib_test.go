package persist

import "testing"

func Test_Persist(t *testing.T) {
	
	p, err := New("arkanoid.db", "Scores")
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
		t.Logf("data[%d]:%s", i, string(v))
	}
}