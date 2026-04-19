package tokenize

import "testing"

func BenchmarkParse_ShallowPath(b *testing.B) {
	tz := New()
	for i := 0; i < b.N; i++ {
		tz.Parse("secret/app")
	}
}

func BenchmarkParse_DeepPath(b *testing.B) {
	tz := New()
	path := "secret/org/team/project/env/service/component/key"
	for i := 0; i < b.N; i++ {
		tz.Parse(path)
	}
}

func BenchmarkJoin_RoundTrip(b *testing.B) {
	tz := New()
	tok := tz.Parse("secret/org/team/project/env/service")
	for i := 0; i < b.N; i++ {
		tz.Join(tok)
	}
}

func BenchmarkParent_Repeated(b *testing.B) {
	tz := New()
	tok := tz.Parse("a/b/c/d/e")
	for i := 0; i < b.N; i++ {
		_, _ = tz.Parent(tok)
	}
}
