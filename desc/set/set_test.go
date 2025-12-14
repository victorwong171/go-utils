package set

import (
	"testing"
)

func TestSet(t *testing.T) {
	s := InitSet[int](12)
	if s.HasKey() {
		t.Error("HasKey() should return false when no keys are provided")
	}

	for i := 0; i < 12; i++ {
		s.Set(i + 1)
	}

	for i := 0; i < 6; i++ {
		s.Drop(i*2 + 1)
	}

	if s.Len() != 6 {
		t.Errorf("Expected 6 keys, got %d", s.Len())
	}

	expectedPresence := map[int]bool{
		1:  false,
		2:  true,
		3:  false,
		4:  true,
		5:  false,
		6:  true,
		7:  false,
		8:  true,
		9:  false,
		10: true,
		11: false,
		12: true,
	}

	for k, v := range expectedPresence {
		if s.HasKey(k) != v {
			t.Errorf("s.HasKey(%d) = %v, want %v", k, s.HasKey(k), v)
		}
	}

	s = s.DropAll()
	if s.Len() != 0 {
		t.Errorf("Expected 0 keys after DropAll, got %d", s.Len())
	}

	s.Set(1, 2, 3)
	slice := s.ToSlice()
	s2 := Setify[int](slice...)
	if s2.Len() != 3 {
		t.Errorf("Expected 3 keys in new set, got %d", s2.Len())
	}

	if !s.HasAny(1) {
		t.Error("s.HasAny(1) should return true")
	}

	if s.HasAny(4) {
		t.Error("s.HasAny(4) should return false")
	}
}
