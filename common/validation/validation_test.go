package validation

import (
	"go-dating-app/app/dto"
	"os"
	"testing"

	"github.com/labstack/gommon/log"
)

func TestMain(m *testing.M) {
	if err := NewValidation(); err != nil {
		log.Fatal("Failed to load validation: ", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestFormatStructErrors(t *testing.T) {
	type args struct {
		Username string `validate:"required,min=6"`
	}
	tests := []struct {
		name       string
		args       args
		wantQtyErr int
	}{
		{name: "struct validation errors can be formatted", args: args{Username: "uname"}, wantQtyErr: 1},
		{name: "struct validation errors can not be formatted", args: args{Username: "validusername"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validationErrors := ValidateStruct(tt.args)
			formattedErrors := FormatStructErrors(validationErrors)
			feLen := len(formattedErrors)
			if feLen != tt.wantQtyErr {
				t.Errorf("FormatStructErrors() = %v, want %v", feLen, tt.wantQtyErr)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	type args struct {
		s any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid struct", args: args{dto.AuthLoginRequest{Email: "validemail@example.com", Password: "validpassword"}}},
		{name: "invalid struct", args: args{dto.AuthLoginRequest{Email: "invalidemailexample", Password: "validpassword"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateStruct(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVar(t *testing.T) {
	type args struct {
		s   string
		tag string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid var", args: args{"validemail@example.com", "email"}},
		{name: "invalid var", args: args{"invalidemailexample", "email"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateVar(tt.args.s, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("ValidateVar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
