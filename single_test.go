package single

import (
	"testing"
)

func TestSingle(t *testing.T) {
	type fields struct {
		name string
	}
	tests := []struct {
		name          string
		fields        fields
		want          bool
		wantErr       bool
		wantErrUnlock bool
	}{
		{
			name: "test1",
			fields: fields{
				name: "single_test1",
			},
			want:          false,
			wantErr:       false,
			wantErrUnlock: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				busy bool
				err  error
			)

			s := New(tt.fields.name)

			busy, err = s.Lock()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if busy != tt.want {
				t.Errorf("Lock() got = %v, want %v", busy, tt.want)
				return
			}

			err = s.Unlock()
			if (err != nil) != tt.wantErrUnlock {
				t.Errorf("TryUnlock() error = %v, wantErrUnlock %v", err, tt.wantErr)
				return
			}
		})
	}
}
