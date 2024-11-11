package handle

import "testing"

func Test_isHTTP(t *testing.T) {
	tests := []struct {
		peek []byte
		want bool
	}{
		{
			peek: []byte("GET /path  "), // 10 bytes
			want: true,
		},
		{
			peek: []byte("POST /path "), // 10 bytes
			want: true,
		},
		{
			peek: []byte("put /path  "), // 10 bytes
			want: true,
		},
		{
			peek: []byte("          "), // 10 spaces
			want: false,
		},
		{
			peek: []byte("INVALID /p "), // 10 bytes
			want: false,
		},
		{
			peek: []byte("GET       "), // 10 bytes
			want: true,
		},
		{
			peek: []byte("GET /p?k=v"), // 10 bytes
			want: true,
		},
		{
			peek: []byte("DELETE /pa"), // 10 bytes
			want: true,
		},
		{
			peek: []byte("PROPPATCH 123123"), // 10 bytes
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run("isHttp", func(t *testing.T) {
			if got := isHTTP(tt.peek); got != tt.want {
				t.Errorf("isHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
