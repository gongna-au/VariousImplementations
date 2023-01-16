package TokenBucket

import "testing"

func TestTokenBucket(t *testing.T) {
	bucket := NewcommonBucket(4)
	bucket.Run()
}
