package gxjpath

import (
	"testing"

	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestThreeSegmentPath(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	defer SetDefaultFailureMode(FailureHalts)
	Convey("Should return compiled path", t, func() {
		path := "k1.k2.@last"
		compiled, err := CompilePath(path)

		Convey("Error should not appear", func() { So(err, ShouldBeNil) })

		Convey("First element must be a simple key with a container of map type and any value type", func() {
			So(compiled[0].index, ShouldEqual, 0)
			So(compiled[0].key, ShouldEqual, "k1")
			So(compiled[0].containerType, ShouldEqual, GXJMapContainer)
		})

		Convey("Second element must be a simple key with a container of map type and any value type", func() {
			So(compiled[1].index, ShouldEqual, 0)
			So(compiled[1].key, ShouldEqual, "k2")
			So(compiled[1].containerType, ShouldEqual, GXJMapContainer)
		})

		Convey("Third element must be an index key for the last slice element of any value type", func() {
			So(compiled[2].index, ShouldEqual, -1)
			So(compiled[2].key, ShouldEqual, "last")
			So(compiled[2].containerType, ShouldEqual, GXJSliceContainer)
		})
	})
}

func TestTwoSegmentPath(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	defer SetDefaultFailureMode(FailureHalts)
	Convey("Should return compiled path with last element of array type", t, func() {
		path := "k.@last"
		compiled, err := CompilePath(path)
		Convey("Error should not appear", func() { So(err, ShouldBeNil) })

		Convey("First element must be a simple key with a container of map type and any value type", func() {
			So(compiled[0].index, ShouldEqual, 0)
			So(compiled[0].key, ShouldEqual, "k")
			So(compiled[0].containerType, ShouldEqual, GXJMapContainer)
		})
		Convey("Second element points to the last element of the slice of map element type as -1", func() {
			So(compiled[1].index, ShouldEqual, -1)
			So(compiled[1].key, ShouldEqual, "last")
			So(compiled[1].containerType, ShouldEqual, GXJSliceContainer)
		})
	})
}

func TestKeyCompileArrayIndex(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	defer SetDefaultFailureMode(FailureHalts)
	Convey("Array reference", t, func() {
		Convey("Must be first array element", func() {
			compiled, err := CompilePath("@first")
			So(err, ShouldBeNil)
			So(compiled[0].index, ShouldEqual, 0)
			So(compiled[0].key, ShouldEqual, "first")
			So(compiled[0].containerType, ShouldEqual, GXJSliceContainer)
		})
		Convey("Must be last array element", func() {
			compiled, err := CompilePath("@last")
			So(err, ShouldBeNil)
			So(compiled[0].index, ShouldEqual, -1)
			So(compiled[0].key, ShouldEqual, "last")
			So(compiled[0].containerType, ShouldEqual, GXJSliceContainer)
		})
		Convey("Must be -100", func() {
			compiled, err := CompilePath("@-100")
			So(err, ShouldBeNil)
			So(compiled[0].index, ShouldEqual, -100)
			So(compiled[0].key, ShouldEqual, "-100")
			So(compiled[0].containerType, ShouldEqual, GXJSliceContainer)
		})
		Convey("Must be 100", func() {
			compiled, err := CompilePath("@111")
			So(err, ShouldBeNil)
			So(compiled[0].index, ShouldEqual, 111)
			So(compiled[0].key, ShouldEqual, "111")
			So(compiled[0].containerType, ShouldEqual, GXJSliceContainer)
		})
	})
}

func TestHowEscapingWorks(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	defer SetDefaultFailureMode(FailureHalts)
	Convey("Unescape should work properly", t, func() {
		Convey("Key must be k1@", func() {
			compiled, err := CompilePath("k1\\.k2")
			So(err, ShouldBeNil)
			So(compiled[0].index, ShouldEqual, 0)
			So(compiled[0].key, ShouldEqual, "k1.k2")
			So(compiled[0].containerType, ShouldEqual, GXJMapContainer)
		})

		Convey("Key must be @k", func() {
			compiled, err := CompilePath("\\@k1")
			So(err, ShouldBeNil)
			So(compiled[0].index, ShouldEqual, 0)
			So(compiled[0].key, ShouldEqual, "@k1")
			So(compiled[0].containerType, ShouldEqual, GXJMapContainer)
		})
	})
}

func TestKeyCompileError(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	defer SetDefaultFailureMode(FailureHalts)
	Convey("Error should appear", t, func() {
		Convey("Invalid empty path: .", func() {
			_, err := CompilePath(".")
			So(err, ShouldEqual, ErrWrongPath)
		})
	})
	Convey("Invalid array index", t, func() {
		Convey("Invalid empty path: .", func() {
			_, err := CompilePath("@art")
			So(err, ShouldEqual, ErrInvalidIndex)
		})
		Convey("Invalid empty path: k1.", func() {
			_, err := CompilePath("k1.")
			So(err, ShouldEqual, ErrWrongPath)
		})
	})
}

func getTestJson() interface{} {
	testJsonStructure := `
	{"intkey": 123,
	 "strkey": "str",
	 "floatkey": 1.1,
	 "anyarraykey": [0.1, 1.1, 2, "somestr", {"k1":"v", "k2":1, "k3": {"ik": 222}}, 5],
	 "anymapkey": {"k1":"v", "k2":1, "k3": {"ik": 2}}
	}
	`
	var parsed interface{}
	err := json.Unmarshal([]byte(testJsonStructure), &parsed)
	if err != nil {
		panic(err)
	}
	return parsed
}

func TestLookValues(t *testing.T) {
	Convey("Extract value", t, func() {
		Convey("Extract int value", func() {
			data := getTestJson()
			v, err := LookupRawPath("intkey", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 123)
		})
		Convey("Extract str value", func() {
			data := getTestJson()
			v, err := LookupRawPath("strkey", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "str")
		})
		Convey("Extract float value", func() {
			data := getTestJson()
			v, err := LookupRawPath("floatkey", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 1.1)
		})
		Convey("Extract slice value", func() {
			data := getTestJson()
			v, err := LookupRawPath("anyarraykey", data)
			So(err, ShouldBeNil)
			So(v, ShouldHaveLength, 6)
		})
		Convey("Extract map value", func() {
			data := getTestJson()
			v, err := LookupRawPath("anymapkey", data)
			So(err, ShouldBeNil)
			So(v, ShouldHaveLength, 3)
		})
		Convey("Extract sub first array element", func() {
			data := getTestJson()
			v, err := LookupRawPath("anyarraykey.@first", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 0.1)
		})
		Convey("Extract sub second array element", func() {
			data := getTestJson()
			v, err := LookupRawPath("anyarraykey.@1", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 1.1)
		})
		Convey("Extract sub third array element", func() {
			data := getTestJson()
			v, err := LookupRawPath("anyarraykey.@2", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 2)
		})
		Convey("Extract last array element", func() {
			data := getTestJson()
			v, err := LookupRawPath("anyarraykey.@last", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 5)
		})
		Convey("Extract last array element again", func() {
			data := getTestJson()
			v, err := LookupRawPath("anyarraykey.@-1", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 5)
		})
		Convey("Extract very deep value", func() {
			data := getTestJson()
			v, err := LookupRawPath("anyarraykey.@4.k3.ik", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 222)
		})
	})
}

func TestCompileCache(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	defer SetDefaultFailureMode(FailureHalts)
	Convey("All lookups should occuer", t, func() {
		Convey("Lookup value first time", func() {
			data := getTestJson()
			v, err := CachedLookup("anyarraykey.@4.k3.ik", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 222)
		})
		Convey("Lookup value second time", func() {
			data := getTestJson()
			v, err := CachedLookup("anyarraykey.@4.k3.ik", data)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 222)
		})
	})
}
