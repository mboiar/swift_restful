// Utils implements utility functions for processing SWIFT data.

package utils

import "testing"

func Test_isAlphanumeric(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"non_alphanumeric", args{"123qwe%^"}, false},
		{"alphanumeric", args{"123qwe"}, true},
		{"empty", args{""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAlphanumeric(tt.args.str); got != tt.want {
				t.Errorf("isAlphanumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidSwiftCode(t *testing.T) {
	type args struct {
		swiftCode string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid_swift_code", args{"QWERTY123tt"}, true},
		{"invalid_swift_code_length", args{"qwe"}, false},
		{"invalid_swift_code_alphanum", args{"qwerty123%%"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSwiftCode(tt.args.swiftCode); got != tt.want {
				t.Errorf("IsValidSwiftCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsHeadquarter(t *testing.T) {
	type args struct {
		swiftCode string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"valid_headquarter", args{"qewrty12XXX"}, true, false},
		{"valid_headquarter_lower", args{"12345678xxx"}, true, false},
		{"invalid_headquarter", args{"12345678eXX"}, false, false},
		{"invalid_swift_code", args{"12345678eX"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsHeadquarter(tt.args.swiftCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsHeadquarter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsHeadquarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidISO2Code(t *testing.T) {
	type args struct {
		iso2code string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidISO2Code(tt.args.iso2code); got != tt.want {
				t.Errorf("IsValidISO2Code() = %v, want %v", got, tt.want)
			}
		})
	}
}
