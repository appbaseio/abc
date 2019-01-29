// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package env

import "testing"
import "time"

func mockedGetter(env map[string]string) envGetter {
	return func(key string) string {
		if v, ok := env[key]; ok {
			return v
		}
		return ""
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		Env      map[string]string
		Default  string
		Keys     []string
		Expected string
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  "default",
			Keys:     []string{"KEY"},
			Expected: "default",
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": "value",
			},
			Default:  "",
			Keys:     []string{"KEY"},
			Expected: "value",
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": "value",
			},
			Default:  "default",
			Keys:     []string{"KEY"},
			Expected: "default",
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": "value",
			},
			Default:  "",
			Keys:     []string{"MISSING", "KEY"},
			Expected: "value",
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := String(tt.Default, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %q, want %q", i, have, want)
		}
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		Env      map[string]string
		Default  int
		Keys     []string
		Expected int
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  -1,
			Keys:     []string{"KEY"},
			Expected: -1,
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": "42",
			},
			Default:  0,
			Keys:     []string{"KEY"},
			Expected: 42,
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": "42",
			},
			Default:  -1,
			Keys:     []string{"KEY"},
			Expected: -1,
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": "42",
			},
			Default:  -1,
			Keys:     []string{"MISSING", "KEY"},
			Expected: 42,
		},
		// #4: Ignore conversion errors
		{
			Env: map[string]string{
				"KEY": "not-an-integer",
			},
			Default:  -1,
			Keys:     []string{"KEY"},
			Expected: -1,
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := Int(tt.Default, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %v, want %v", i, have, want)
		}
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		Env      map[string]string
		Default  int64
		Keys     []string
		Expected int64
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  -1,
			Keys:     []string{"KEY"},
			Expected: -1,
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": "42",
			},
			Default:  0,
			Keys:     []string{"KEY"},
			Expected: 42,
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": "42",
			},
			Default:  -1,
			Keys:     []string{"KEY"},
			Expected: -1,
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": "42",
			},
			Default:  -1,
			Keys:     []string{"MISSING", "KEY"},
			Expected: 42,
		},
		// #4: Ignore conversion errors
		{
			Env: map[string]string{
				"KEY": "not-an-integer",
			},
			Default:  -1,
			Keys:     []string{"KEY"},
			Expected: -1,
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := Int64(tt.Default, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %v, want %v", i, have, want)
		}
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		Env      map[string]string
		Default  float32
		Keys     []string
		Expected float32
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  -1.0,
			Keys:     []string{"KEY"},
			Expected: -1.0,
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": "42.5",
			},
			Default:  0,
			Keys:     []string{"KEY"},
			Expected: 42.5,
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": "42.5",
			},
			Default:  -1.0,
			Keys:     []string{"KEY"},
			Expected: -1.0,
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": "42.5",
			},
			Default:  -1.0,
			Keys:     []string{"MISSING", "KEY"},
			Expected: 42.5,
		},
		// #4: Ignore conversion errors
		{
			Env: map[string]string{
				"KEY": "not-an-integer",
			},
			Default:  -1.2,
			Keys:     []string{"KEY"},
			Expected: -1.2,
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := Float32(tt.Default, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %v, want %v", i, have, want)
		}
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		Env      map[string]string
		Default  float64
		Keys     []string
		Expected float64
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  -1.0,
			Keys:     []string{"KEY"},
			Expected: -1.0,
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": "42.5",
			},
			Default:  0,
			Keys:     []string{"KEY"},
			Expected: 42.5,
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": "42.5",
			},
			Default:  -1.0,
			Keys:     []string{"KEY"},
			Expected: -1.0,
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": "42.5",
			},
			Default:  -1.0,
			Keys:     []string{"MISSING", "KEY"},
			Expected: 42.5,
		},
		// #4: Ignore conversion errors
		{
			Env: map[string]string{
				"KEY": "not-an-integer",
			},
			Default:  -1.2,
			Keys:     []string{"KEY"},
			Expected: -1.2,
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := Float64(tt.Default, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %v, want %v", i, have, want)
		}
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		Env      map[string]string
		Default  bool
		Keys     []string
		Expected bool
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  true,
			Keys:     []string{"KEY"},
			Expected: true,
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": "false",
			},
			Default:  true,
			Keys:     []string{"KEY"},
			Expected: false,
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": "false",
			},
			Default:  true,
			Keys:     []string{"KEY"},
			Expected: true,
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": "false",
			},
			Default:  true,
			Keys:     []string{"MISSING", "KEY"},
			Expected: false,
		},
		// #4: Ignore conversion errors
		{
			Env: map[string]string{
				"KEY": "not-a-bool",
			},
			Default:  true,
			Keys:     []string{"KEY"},
			Expected: true,
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := Bool(tt.Default, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %v, want %v", i, have, want)
		}
	}
}

func TestTime(t *testing.T) {
	t2000 := time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC).Truncate(time.Second)
	t2011 := time.Date(2011, 12, 31, 0, 0, 0, 0, time.UTC).Truncate(time.Second)

	tests := []struct {
		Env      map[string]string
		Default  time.Time
		Format   string
		Keys     []string
		Expected time.Time
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  t2000,
			Format:   time.RFC3339,
			Keys:     []string{"KEY"},
			Expected: t2000,
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": t2011.Format(time.RFC3339),
			},
			Default:  t2000,
			Format:   time.RFC3339,
			Keys:     []string{"KEY"},
			Expected: t2011,
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": t2011.Format(time.RFC3339),
			},
			Default:  t2000,
			Format:   time.RFC3339,
			Keys:     []string{"KEY"},
			Expected: t2000,
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": t2011.Format(time.RFC3339),
			},
			Default:  t2000,
			Format:   time.RFC3339,
			Keys:     []string{"MISSING", "KEY"},
			Expected: t2011,
		},
		// #4: Ignore conversion errors
		{
			Env: map[string]string{
				"KEY": "not-a-date",
			},
			Default:  t2000,
			Format:   time.RFC3339,
			Keys:     []string{"KEY"},
			Expected: t2000,
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := Time(tt.Default, tt.Format, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %v, want %v", i, have.Format(time.RFC3339), want.Format(time.RFC3339))
		}
	}
}

func TestDuration(t *testing.T) {
	d1 := time.Duration(381 * time.Second)
	d2, _ := time.ParseDuration("1h12m36s")

	tests := []struct {
		Env      map[string]string
		Default  time.Duration
		Keys     []string
		Expected time.Duration
	}{
		// #0: Return default
		{
			Env:      map[string]string{},
			Default:  d1,
			Keys:     []string{"KEY"},
			Expected: d1,
		},
		// #1: Find key in environment
		{
			Env: map[string]string{
				"KEY": d2.String(),
			},
			Default:  d1,
			Keys:     []string{"KEY"},
			Expected: d2,
		},
		// #2: Honor case-insensitivity
		{
			Env: map[string]string{
				"key": d2.String(),
			},
			Default:  d1,
			Keys:     []string{"KEY"},
			Expected: d1,
		},
		// #3: Multiple keys passed
		{
			Env: map[string]string{
				"KEY": d2.String(),
			},
			Default:  d1,
			Keys:     []string{"MISSING", "KEY"},
			Expected: d2,
		},
		// #4: Ignore conversion errors
		{
			Env: map[string]string{
				"KEY": "not-a-duration",
			},
			Default:  d1,
			Keys:     []string{"KEY"},
			Expected: d1,
		},
	}

	for i, tt := range tests {
		getter = mockedGetter(tt.Env)
		if have, want := Duration(tt.Default, tt.Keys...), tt.Expected; have != want {
			t.Errorf("#%d: have %v, want %v", i, have, want)
		}
	}
}
