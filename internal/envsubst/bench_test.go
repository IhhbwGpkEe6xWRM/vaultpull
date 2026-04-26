package envsubst

import "testing"

func BenchmarkApply_NoReferences(b *testing.B) {
	s := New()
	src := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_NAME": "mydb",
		"APP_ENV": "production",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Apply(src)
	}
}

func BenchmarkApply_WithReferences(b *testing.B) {
	s := New()
	src := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"DSN": "postgres://${HOST}:${PORT}/db",
		"URL": "http://${HOST}:8080",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Apply(src)
	}
}

func BenchmarkContainsReferences(b *testing.B) {
	val := "postgres://${DB_HOST}:${DB_PORT}/${DB_NAME}"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ContainsReferences(val)
	}
}
