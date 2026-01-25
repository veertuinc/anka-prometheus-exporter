package envflag

import (
	"flag"
	"testing"
	"time"
)

func TestStringVar(t *testing.T) {
	// Reset flags for clean test
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value string
	StringVar(&value, "TEST_STRING", "default", "test string flag")

	f := EnvironmentFlags.Lookup("TEST_STRING")
	if f == nil {
		t.Fatal("Flag TEST_STRING not registered")
	}
	if f.DefValue != "default" {
		t.Errorf("Expected default value 'default', got '%s'", f.DefValue)
	}
}

func TestString(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	value := String("TEST_STRING_PTR", "default_ptr", "test string pointer flag")

	if value == nil {
		t.Fatal("String() returned nil")
	}
	if *value != "default_ptr" {
		t.Errorf("Expected default value 'default_ptr', got '%s'", *value)
	}
}

func TestIntVar(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value int
	IntVar(&value, "TEST_INT", 42, "test int flag")

	f := EnvironmentFlags.Lookup("TEST_INT")
	if f == nil {
		t.Fatal("Flag TEST_INT not registered")
	}
	if f.DefValue != "42" {
		t.Errorf("Expected default value '42', got '%s'", f.DefValue)
	}
}

func TestInt(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	value := Int("TEST_INT_PTR", 100, "test int pointer flag")

	if value == nil {
		t.Fatal("Int() returned nil")
	}
	if *value != 100 {
		t.Errorf("Expected default value 100, got %d", *value)
	}
}

func TestBoolVar(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value bool
	BoolVar(&value, "TEST_BOOL", true, "test bool flag")

	f := EnvironmentFlags.Lookup("TEST_BOOL")
	if f == nil {
		t.Fatal("Flag TEST_BOOL not registered")
	}
	if f.DefValue != "true" {
		t.Errorf("Expected default value 'true', got '%s'", f.DefValue)
	}
}

func TestBool(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	value := Bool("TEST_BOOL_PTR", false, "test bool pointer flag")

	if value == nil {
		t.Fatal("Bool() returned nil")
	}
	if *value != false {
		t.Errorf("Expected default value false, got %v", *value)
	}
}

func TestInt64Var(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value int64
	Int64Var(&value, "TEST_INT64", 9223372036854775807, "test int64 flag")

	f := EnvironmentFlags.Lookup("TEST_INT64")
	if f == nil {
		t.Fatal("Flag TEST_INT64 not registered")
	}
}

func TestUintVar(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value uint
	UintVar(&value, "TEST_UINT", 123, "test uint flag")

	f := EnvironmentFlags.Lookup("TEST_UINT")
	if f == nil {
		t.Fatal("Flag TEST_UINT not registered")
	}
}

func TestUint64Var(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value uint64
	Uint64Var(&value, "TEST_UINT64", 18446744073709551615, "test uint64 flag")

	f := EnvironmentFlags.Lookup("TEST_UINT64")
	if f == nil {
		t.Fatal("Flag TEST_UINT64 not registered")
	}
}

func TestFloat64Var(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value float64
	Float64Var(&value, "TEST_FLOAT64", 3.14159, "test float64 flag")

	f := EnvironmentFlags.Lookup("TEST_FLOAT64")
	if f == nil {
		t.Fatal("Flag TEST_FLOAT64 not registered")
	}
}

func TestDurationVar(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value time.Duration
	DurationVar(&value, "TEST_DURATION", 5*time.Second, "test duration flag")

	f := EnvironmentFlags.Lookup("TEST_DURATION")
	if f == nil {
		t.Fatal("Flag TEST_DURATION not registered")
	}
}

func TestLookup(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	StringVar(new(string), "LOOKUP_TEST", "value", "test lookup")

	f := Lookup("LOOKUP_TEST")
	if f == nil {
		t.Error("Lookup returned nil for existing flag")
	}

	f = Lookup("NONEXISTENT")
	if f != nil {
		t.Error("Lookup should return nil for non-existent flag")
	}
}

func TestSet(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value string
	StringVar(&value, "SET_TEST", "default", "test set")

	err := Set("SET_TEST", "new_value")
	if err != nil {
		t.Errorf("Set returned error: %v", err)
	}

	// After parsing, the value should be updated
	f := Lookup("SET_TEST")
	if f.Value.String() != "new_value" {
		t.Errorf("Expected 'new_value', got '%s'", f.Value.String())
	}
}

func TestVisitAll(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	StringVar(new(string), "VISIT_TEST_1", "value1", "test 1")
	StringVar(new(string), "VISIT_TEST_2", "value2", "test 2")

	count := 0
	VisitAll(func(f *flag.Flag) {
		count++
	})

	if count != 2 {
		t.Errorf("Expected to visit 2 flags, visited %d", count)
	}
}

func TestVisit(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	StringVar(new(string), "VISIT_SET_TEST", "default", "test")
	_ = Set("VISIT_SET_TEST", "changed")

	count := 0
	Visit(func(f *flag.Flag) {
		count++
	})

	if count != 1 {
		t.Errorf("Expected to visit 1 set flag, visited %d", count)
	}
}

func TestParsed(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	if Parsed() {
		t.Error("Parsed() should return false before parsing")
	}

	Parse()

	if !Parsed() {
		t.Error("Parsed() should return true after parsing")
	}
}

func TestParsePrefix(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value string
	StringVar(&value, "TEST_VAR", "default", "test var")

	// Set environment variable with prefix
	t.Setenv("MY_PREFIX_TEST_VAR", "from_env")

	ParsePrefix("MY_PREFIX_")

	if value != "from_env" {
		t.Errorf("Expected 'from_env', got '%s'", value)
	}
}

func TestParsePrefix_NoMatch(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value string
	StringVar(&value, "UNMATCHED_VAR", "default", "test var")

	// Set environment variable without matching prefix
	t.Setenv("OTHER_PREFIX_UNMATCHED_VAR", "should_not_match")

	ParsePrefix("MY_PREFIX_")

	if value != "default" {
		t.Errorf("Expected 'default' (no match), got '%s'", value)
	}
}

func TestParsePrefix_CaseInsensitive(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value string
	StringVar(&value, "CASE_TEST", "default", "test var")

	// Set environment variable with different case prefix
	t.Setenv("my_prefix_CASE_TEST", "case_insensitive")

	ParsePrefix("MY_PREFIX_")

	if value != "case_insensitive" {
		t.Errorf("Expected 'case_insensitive', got '%s'", value)
	}
}

func TestParse(t *testing.T) {
	EnvironmentFlags = flag.NewFlagSet("test", flag.ContinueOnError)

	var value string
	StringVar(&value, "PARSE_TEST", "default", "test var")

	// Set environment variable without prefix
	t.Setenv("PARSE_TEST", "no_prefix_value")

	Parse() // Parse without prefix

	if value != "no_prefix_value" {
		t.Errorf("Expected 'no_prefix_value', got '%s'", value)
	}
}
