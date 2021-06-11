package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	//replicas = 3
	hash := MapNew(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	//add three node : 2/4/6
	//virtual node : 02/12/22\04/14/24\06/16/26
	hash.MapAdd("6", "4", "2")

	//test: 2, 11, 23, 27  ----->  node : 02, 12, 24, 02 -----> real node : 02, 02, 04, 02
	testCases := map[string]string {
		"2": "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.MapGet(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	hash.MapAdd("8")
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.MapGet(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}