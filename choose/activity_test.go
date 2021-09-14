package choose

import "testing"

func TestMergeImage2(t *testing.T) {
	MergeImage2()
}

func BenchmarkMergeImage2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MergeImage()
	}
}
