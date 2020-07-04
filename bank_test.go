package bank

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestNewAccount(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Account
	}{
		{"Griesemer", args{"Griesemer"}, &Account{Name: "Griesemer", Bal: 0, Hist: nil}},
		{"Pike", args{"Pike"}, &Account{Name: "Pike", Bal: 0, Hist: nil}},
		{"Thompson", args{"Thompson"}, &Account{Name: "Thompson", Bal: 0, Hist: nil}},
	}
	accounts = map[string]*Account{}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure that the correct account is created
			if got := NewAccount(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccount() = %v, want %v", got, tt.want)
			}
			// Ensure each account gets inserted into accounts
			if i+1 != len(accounts) || !reflect.DeepEqual(accounts[tt.name], tt.want) {
				t.Errorf("len(accounts) = %v, want %v\naccounts[\"%v\"] is not %v\n", len(accounts), i+1, tt.name, tt.want)
			}
		})
	}
}

func ExampleNewAccount() {
	a := NewAccount("Test")
	fmt.Println(Name(a))
	// Output: Test
}

func TestName(t *testing.T) {
	type args struct {
		a *Account
	}

	pike := &Account{"Pike", 100, nil}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Pike", args{pike}, "Pike"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Name(tt.args.a); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBalance(t *testing.T) {
	type args struct {
		a *Account
	}
	pike := &Account{"Pike", 100, nil}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"Pike 100", args{pike}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Balance(tt.args.a); got != tt.want {
				t.Errorf("Balance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	type args struct {
		a *Account
		m int
	}

	griesemer := &Account{"Griesemer", 100, nil}
	pike := &Account{"Pike", 0, nil}
	thompson := &Account{"Thompson", 0, nil}

	tests := []struct {
		name    string
		args    args
		want    int
		hist    []history
		wantErr bool
	}{
		{"Griesemer deposits 100", args{griesemer, 100}, 200, []history{{100, 200}}, false},
		{"Pike deposits 42", args{pike, 42}, 42, []history{{42, 42}}, false},
		{"Pike deposits -1", args{pike, -1}, 42, []history{{42, 42}}, true},
		{"Thompson deposits 60", args{thompson, 60}, 60, []history{{60, 60}}, false},
		{"Thompson deposits 99", args{thompson, 39}, 99, []history{{60, 60}, {39, 99}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Deposit(tt.args.a, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Deposit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	type args struct {
		a *Account
		m int
	}

	griesemer := &Account{"Griesemer", 100, nil}
	pike := &Account{"Pike", 100, nil}
	thompson := &Account{"Thompson", 100, nil}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"Griesemer withdraws 100", args{griesemer, 100}, 0, false},
		{"Pike withdraws 42", args{pike, 42}, 58, false},
		{"Pike withdraws -1", args{pike, -1}, 58, true},
		{"Thompson withdraws 60", args{thompson, 101}, 100, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Withdraw(tt.args.a, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Withdraw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	type args struct {
		a *Account
		b *Account
		m int
	}
	griesemer := &Account{"Griesemer", 100, nil}
	pike := &Account{"Pike", 100, nil}
	thompson := &Account{"Thompson", 100, nil}

	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{"Griesemer transfers 100 to Pike", args{griesemer, pike, 100}, 0, 200, false},
		{"Griesemer transfers 100 to Pike again", args{griesemer, pike, 100}, 0, 200, true},
		{"Pike transfers 300 to Thompson", args{pike, thompson, 300}, 200, 100, true},
		{"Pike transfers -100 to Thompson", args{pike, thompson, -100}, 200, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Transfer(tt.args.a, tt.args.b, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Transfer() got = %v, want %v", got, tt.want)
			}
			if !tt.wantErr && got1 != tt.want1 {
				t.Errorf("Transfer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHistory(t *testing.T) {
	type args struct {
		a Account
	}

	pike := Account{"Pike", 100, nil}
	pike.Hist = []history{
		{100, 100},
		{10, 110},
		{-40, 70},
		{23, 93},
	}

	tests := []struct {
		name     string
		args     args
		wantAmt  []int
		wantBal  []int
		wantMore []bool
	}{
		{"Pike's account history", args{pike}, []int{100, 10, -40, 23}, []int{100, 110, 70, 93}, []bool{true, true, true, false}},
	}
	for _, tt := range tests {
		h := History(&pike)
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < len(pike.Hist); i++ {
				amt, bal, more := h()
				if amt != tt.wantAmt[i] || bal != tt.wantBal[i] || more != tt.wantMore[i] {
					t.Errorf("History() = %v, %v, %v, want %v, %v, %v", amt, bal, more, tt.wantAmt[i], tt.wantBal[i], tt.wantMore[i])
				}
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"SaveAndLoad"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewAccount("Hiasl")
			if err := Save(); err != nil {
				t.Errorf("Save() error = %v, stack = %v", err, errors.WithStack(err))
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			accounts = nil
			if err := Load(); err != nil {
				t.Errorf("Load() error = %v, stack = %v", err, errors.WithStack(err))
			}
			if accounts == nil {
				t.Errorf("accounts not restored: %v", accounts)
			}
			hiasl, ok := accounts["Hiasl"]
			if !ok {
				t.Errorf("accounts = %v, hiasl = %v", accounts, hiasl)
			}
		})
	}
}
