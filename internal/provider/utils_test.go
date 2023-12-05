package provider

import "testing"

func TestConvertToMongoIndexTypeAsc(t *testing.T) {
	val := convertToMongoIndexType("asc")
	want := 1
	if want != val {
		t.Fatalf("Expected %v, got %v", want, val)
	}
}

func TestConvertToMongoIndexTypeDesc(t *testing.T) {
	val := convertToMongoIndexType("desc")
	want := -1
	if want != val {
		t.Fatalf("Expected %v, got %v", want, val)
	}
}

func TestConvertToMongoIndexTypeStr(t *testing.T) {
	val := convertToMongoIndexType("text")
	want := "text"
	if want != val {
		t.Fatalf("Expected %v, got %v", want, val)
	}
}

func TestConvertToTfIndexTypeAsc(t *testing.T) {
	val, err := convertToTfIndexType(int32(1))
	want := "asc"
	if want != val {
		t.Fatalf("Expected %v, got %v, err %v", want, val, err)
	}
}

func TestConvertToTfIndexTypeDesc(t *testing.T) {
	val, err := convertToTfIndexType(int32(-1))
	want := "desc"
	if want != val {
		t.Fatalf("Expected %v, got %v, err %v", want, val, err)
	}
}

func TestConvertToTfIndexTypeStr(t *testing.T) {
	val, err := convertToTfIndexType("2dsphere")
	want := "2dsphere"
	if want != val {
		t.Fatalf("Expected %v, got %v, err %v", want, val, err)
	}
}

func TestConvertToTfIndexTypeTypeError(t *testing.T) {
	val, err := convertToTfIndexType(float32(1))
	if err == nil {
		t.Fatalf("Expected an error, got %v, ", val)
	}
}

func TestParseValidIndexId(t *testing.T) {
	id, err := parseIndexId("db.collec.index_name")

	want := indexId{
		database:   "db",
		collection: "collec",
		indexName:  "index_name",
	}

	if want != *id {
		t.Fatalf("Expected %v, got %v, err %v", want, id, err)
	}
}

func TestParseInvalidIndex(t *testing.T) {
	_, err := parseIndexId("db.collec")

	if err == nil {
		t.Fatalf("Should have failed")
	}
}
