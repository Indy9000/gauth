package storage

import (
	"testing"
	"time"
)

type MockUserProfile struct {
	Name string
	ID   string
}

func TestSessionTokenGet1(t *testing.T) {
	s := NewSessionCache(time.Second * 3)
	s.Set("key1", &MockUserProfile{Name: "BOB", ID: "001"})

	k1, _ := s.Get("A") //must fail
	if k1 != nil {
		t.Errorf("Expected session key to not exist, but did")
	}

	k2, _ := s.Get("key1") //must succeed
	if k2 == nil {
		t.Error("Expected session key to exist, but it didn`t")
	}
	up := k2.(*MockUserProfile)
	if up.Name != "BOB" {
		t.Errorf("Expected user to be BOB but it was %s", up.Name)
	}
	if up.ID != "001" {
		t.Errorf("Expected user to be BOB but it was %s", up.Name)
	}

}

func TestScanAndRemoveExpiredTokens(t *testing.T) {
	sc := &SessionCache{
		keys: map[string]*CacheItem{
			"K": &CacheItem{
				timestamp: time.Now().UTC(), value: &MockUserProfile{
					Name: "ALICE", ID: "002",
				},
			},
		},
	}
	expiry := time.Second * 3
	time.Sleep(time.Second * 2)
	sc.scanAndRemoveExpiredTokens(expiry)
	if _, ok := sc.keys["K"]; !ok {
		t.Errorf("Expected session key to exist, but didn't")
	}
	time.Sleep(time.Second * 2)
	sc.scanAndRemoveExpiredTokens(expiry)
	if _, ok := sc.keys["K"]; ok {
		t.Errorf("Expected session key to not exist, but did")
	}
}
