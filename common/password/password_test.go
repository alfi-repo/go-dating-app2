package password

import "testing"

func TestHash(t *testing.T) {
	longPassword := "Loremipsumdolorsitamet,consecteturadipiscingelit,seddoeiusmodtemporincididuntutlaboreetdoloremagnaaliqua"
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid password", args: args{s: "password"}},
		{name: "invalid password", args: args{s: longPassword}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Hash(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestVerify(t *testing.T) {
	validPassword := "validpassword"
	validPasswordHash, _ := Hash(validPassword)
	invalidPasswordHash := "12345678"

	type args struct {
		plain string
		hash  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "valid password", args: args{plain: validPassword, hash: validPasswordHash}, want: true},
		{name: "invalid password", args: args{plain: validPassword, hash: invalidPasswordHash}, want: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Verify(tt.args.plain, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Verify() got = %v, want %v", got, tt.want)
			}
		})
	}
}
