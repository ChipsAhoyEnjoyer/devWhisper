package server

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_hashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name              string
		args              args
		notWantedHashedPW string
		wantErr           bool
	}{
		{
			name: "Test hashing a valid password",
			args: args{
				password: "validPassword123!",
			},
			notWantedHashedPW: "validPassword123!", // We can't predict the hash, so we won't check it here
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHashedPW, err := hashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHashedPW == tt.notWantedHashedPW {
				t.Errorf("hashPassword() = %v, want %v", gotHashedPW, tt.notWantedHashedPW)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Check password hash",
			args: args{
				password: "validPassword1!",
				hash:     "validPassword1!",
			},
			wantErr: false,
		},
		{
			name: "Wrong password hash",
			args: args{
				password: "wrongPassword1!",
				hash:     "validPassword1!",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := hashPassword(tt.args.hash)
			if err != nil {
				t.Errorf("hashPassword error'd in CheckHashPassword test: %v", err)
			}
			if err := CheckPasswordHash(hash, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	type makeArgs struct {
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}
	type valArgs struct {
		tokenSecret string
	}
	tests := []struct {
		name        string
		makeArgs    makeArgs
		valArgs     valArgs
		wantMakeErr bool
		wantValErr  bool
	}{
		{
			name: "Validate jwt",
			makeArgs: makeArgs{
				userID:      uuid.New(),
				tokenSecret: "Secret",
				expiresIn:   time.Hour,
			},
			valArgs: valArgs{
				tokenSecret: "Secret",
			},
			wantMakeErr: false,
			wantValErr:  false,
		},
		{
			name: "Wrong secret",
			makeArgs: makeArgs{
				userID:      uuid.New(),
				tokenSecret: "wrong Secret",
				expiresIn:   time.Hour,
			},
			valArgs: valArgs{
				tokenSecret: "Secret",
			},
			wantMakeErr: false,
			wantValErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.makeArgs.userID, tt.makeArgs.tokenSecret, tt.makeArgs.expiresIn)
			if (err != nil) != tt.wantMakeErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantMakeErr)
				return
			}
			got, err := ValidateJWT(token, tt.valArgs.tokenSecret)
			if (err != nil) != tt.wantValErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantValErr)
				return
			}
			want := tt.makeArgs.userID
			if tt.wantMakeErr || tt.wantValErr {
				want = uuid.Nil
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ValidateJWT() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {

	type args struct {
		headers http.Header
		want    string
	}
	getBearerArg := args{
		headers: make(http.Header),
		want:    "secret",
	}
	noBearerArg := args{
		headers: make(http.Header),
		want:    "",
	}
	noHeaderArg := args{
		headers: make(http.Header),
		want:    "",
	}
	getBearerArg.headers.Add("Authorization", "Bearer secret")
	noBearerArg.headers.Add("Authorization", "1234567890 secret")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Get bearer",
			args:    getBearerArg,
			wantErr: false,
		},
		{
			name:    "No bearer",
			args:    noBearerArg,
			wantErr: true,
		},
		{
			name:    "No header found",
			args:    noHeaderArg,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.args.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.args.want)
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "make token",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MakeRefreshToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_isAllowedRune(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid - a",
			args: args{r: 'a'},
			want: true,
		},
		{
			name: "Valid - z",
			args: args{r: 'z'},
			want: true,
		},
		{
			name: "Valid - A",
			args: args{r: 'A'},
			want: true,
		},
		{
			name: "Valid - Z",
			args: args{r: 'Z'},
			want: true,
		},
		{
			name: "Valid - underscore",
			args: args{r: '_'},
			want: true,
		},
		{
			name: "Valid - space",
			args: args{r: ' '},
			want: true,
		},
		{
			name: "Invalid",
			args: args{r: '$'},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAllowedRune(tt.args.r); got != tt.want {
				t.Errorf("isAllowedRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ValidUser",
			args:    args{username: "Valid User_1"},
			wantErr: false,
		},
		{
			name:    "Invalid",
			args:    args{username: "$Invalid$"},
			wantErr: true,
		},
		{
			name:    "Short Length",
			args:    args{username: "2"},
			wantErr: true,
		},
		{
			name:    "Too long",
			args:    args{username: "UsernameIsTooDangLong1234567890"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateUsername(tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("validateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ValidUser",
			args:    args{password: "Valid Pass_1"},
			wantErr: false,
		},
		{
			name:    "Invalid",
			args:    args{password: "$Invalid$"},
			wantErr: true,
		},
		{
			name:    "Short Length",
			args:    args{password: "2"},
			wantErr: true,
		},
		{
			name:    "Too long",
			args:    args{password: "PasswordIsTooDangLong1234567890"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("validatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
