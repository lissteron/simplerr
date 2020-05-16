package simplerr

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestIs(t *testing.T) {
	var (
		e1 = errors.New("t1")
		e2 = errors.New("t2")
	)

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true1",
			args: args{
				err:    e1,
				target: e1,
			},
			want: true,
		},
		{
			name: "true2",
			args: args{
				err:    &withCode{err: e1},
				target: e1,
			},
			want: true,
		},
		{
			name: "true3",
			args: args{
				err:    e1,
				target: &withCode{err: e1},
			},
			want: true,
		},
		{
			name: "true4",
			args: args{
				err:    &withCode{err: &withCode{err: e1}},
				target: e1,
			},
			want: true,
		},
		{
			name: "false1",
			args: args{
				err:    e1,
				target: e2,
			},
			want: false,
		},
		{
			name: "false2",
			args: args{
				err:    &withCode{err: e2},
				target: e1,
			},
			want: false,
		},
		{
			name: "false3",
			args: args{
				err:    e1,
				target: &withCode{err: e2},
			},
			want: false,
		},
		{
			name: "false4",
			args: args{
				err:    e1,
				target: &withCode{err: &withCode{err: e1}},
			},
			want: false,
		},
		{
			name: "nil",
			args: args{
				err:    nil,
				target: &withCode{err: &withCode{err: e1}},
			},
			want: false,
		},
		{
			name: "nil2",
			args: args{
				err:    e1,
				target: nil,
			},
			want: false,
		},
		{
			name: "nil3",
			args: args{
				err:    nil,
				target: nil,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withCode_Error(t *testing.T) {
	type fields struct {
		err  error
		msg  string
		code int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "pass1",
			fields: fields{
				err:  errors.New("t1"),
				msg:  "t2",
				code: 42,
			},
			want: "t2: t1",
		},
		{
			name: "pass2",
			fields: fields{
				err:  &withCode{err: errors.New("t1"), msg: "t2"},
				msg:  "t3",
				code: 42,
			},
			want: "t3: t2: t1",
		},
		{
			name: "nil",
			fields: fields{
				err:  &withCode{err: nil, msg: "t2"},
				msg:  "t3",
				code: 42,
			},
			want: "t3: ",
		},
		{
			name: "nil2",
			fields: fields{
				err:  nil,
				msg:  "t3",
				code: 42,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &withCode{
				err:  tt.fields.err,
				msg:  tt.fields.msg,
				code: tt.fields.code,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("withCode.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasCode(t *testing.T) {
	type args struct {
		err  error
		code int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				err:  &withCode{err: errors.New("t1"), msg: "t2", code: 42},
				code: 42,
			},
			want: true,
		},
		{
			name: "true2",
			args: args{
				err:  &withCode{err: &withCode{err: errors.New("t1"), msg: "t2", code: 42}, msg: "t2", code: 24},
				code: 42,
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				err:  errors.New("t1"),
				code: 42,
			},
			want: false,
		},
		{
			name: "false2",
			args: args{
				err:  nil,
				code: 42,
			},
			want: false,
		},
		{
			name: "false3",
			args: args{
				err:  &withCode{err: &withCode{err: errors.New("t1"), msg: "t2", code: 424}, msg: "t2", code: 24},
				code: 42,
			},
			want: false,
		},
		{
			name: "false4",
			args: args{
				err:  &withCode{err: nil, msg: "t2", code: 24},
				code: 42,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasCode(tt.args.err, tt.args.code); got != tt.want {
				t.Errorf("HasCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "pass",
			args: args{
				err: errors.New("t1"),
				msg: "t2",
			},
			want: "t2: t1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Wrap(tt.args.err, tt.args.msg)

			if got := e.Error(); got != tt.want {
				t.Errorf("Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapf(t *testing.T) {
	type args struct {
		err  error
		tmpl string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "pass",
			args: args{
				err:  errors.New("t1"),
				tmpl: "this is %s",
				args: []interface{}{"error"},
			},
			want: "this is error: t1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Wrapf(tt.args.err, tt.args.tmpl, tt.args.args...)

			if got := e.Error(); got != tt.want {
				t.Errorf("Wrapf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStack(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want []Call
	}{
		{
			name: "pass",
			args: args{
				err: &withCode{
					stack: []Call{
						{
							Line:     12,
							File:     "main.go",
							FuncName: "main.main",
						},
						{
							Line:     24,
							File:     "main.go",
							FuncName: "main.test",
						},
					},
				},
			},
			want: []Call{
				{
					Line:     12,
					File:     "main.go",
					FuncName: "main.main",
				},
				{
					Line:     24,
					File:     "main.go",
					FuncName: "main.test",
				},
			},
		},
		{
			name: "nil",
			args: args{
				err: errors.New("t1"),
			},
			want: nil,
		},
		{
			name: "nil2",
			args: args{
				err: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStack(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "pass",
			args: args{
				err: &withCode{code: 42},
			},
			want: 42,
		},
		{
			name: "zero",
			args: args{
				err: errors.New("t1"),
			},
			want: 0,
		},
		{
			name: "zero2",
			args: args{
				err: nil,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCode(tt.args.err); got != tt.want {
				t.Errorf("GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkHasCode(b *testing.B) {
	err := WrapWithCode(errors.New("t1"), 12, "asd")

	for i := 0; i < b.N; i++ {
		HasCode(err, 12)
	}
}

func BenchmarkHasCode100(b *testing.B) {
	err := WrapWithCode(errors.New("t1"), 9999, "-1")

	for i := int64(0); i < 100; i++ {
		err = WrapWithCode(err, i, strconv.FormatInt(i, 10))
	}

	for i := 0; i < b.N; i++ {
		HasCode(err, 9999)
	}
}

func BenchmarkHasCode1000(b *testing.B) {
	err := WrapWithCode(errors.New("t1"), 9999, "-1")

	for i := int64(0); i < 1000; i++ {
		err = WrapWithCode(err, i, strconv.FormatInt(i, 10))
	}

	for i := 0; i < b.N; i++ {
		HasCode(err, 9999)
	}
}

func BenchmarkError(b *testing.B) {
	err := WrapWithCode(errors.New("t1"), 9999, "-1")

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkError100(b *testing.B) {
	err := WrapWithCode(errors.New("t1"), 9999, "-1")

	for i := int64(0); i < 100; i++ {
		err = WrapWithCode(err, i, strconv.FormatInt(i, 10))
	}

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkError1000(b *testing.B) {
	err := WrapWithCode(errors.New("t1"), 9999, "-1")

	for i := int64(0); i < 1000; i++ {
		err = WrapWithCode(err, i, strconv.FormatInt(i, 10))
	}

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
