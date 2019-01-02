package utils

import (
	"testing"
	"encoding/json"
)

func TestSlug(t *testing.T) {

	slug := Slug(" This    is a- 'test?")

	AssertEquals(t, "Slug", "this-is-a-test", slug)
}

type testStruct struct {
	StringVar1 string `json:"stringvar-1"`
	StringVar2 string `json:"stringvar-2,omitempty"`
	IntVar1    int    `json:"intvar-1"`
	IntVar2    int    `json:"intvar-2,omitempty"`
}

func TestJsonStringify(t *testing.T) {

	data1 := testStruct{
		StringVar1: "testing",
		IntVar1:    999,
	}

	string1 := JsonStringify(data1)
	desired1 := `{"stringvar-1":"testing","intvar-1":999}`

	AssertEquals(t, "Data1 JSON", desired1, string1)

	data2 := testStruct{
		StringVar1: "testing",
		IntVar1:    999,
		StringVar2: "testing2",
		IntVar2:    666,
	}

	string2 := JsonStringify(data2)
	desired2 := `{"stringvar-1":"testing","stringvar-2":"testing2","intvar-1":999,"intvar-2":666}`

	AssertEquals(t, "Data2 JSON", desired2, string2)
}

func TestJsonStack(t *testing.T) {

	mockStack := []byte(`line1
	line2
	line3
	line4`)

	out := JsonStack("panic", mockStack)

	traceData := struct {
		Panic string
		Stack []string
	}{}

	err := json.Unmarshal([]byte(out), &traceData)

	AssertNoError(t, "JsonStack parses without error", err)
	AssertEquals(t, "JsonStrack returns correct panic message:", "panic", traceData.Panic)
	AssertEquals(t, "JsonStrack returns correct number of stack lines:", 4, len(traceData.Stack))
	AssertEquals(t, "JsonStrack correctly strips lines", "line2", traceData.Stack[1])

	out2 := JsonStack(func(){}, mockStack)

	err2 := json.Unmarshal([]byte(out2), &traceData)
	AssertNoError(t, "JsonStack parses without error", err2)
	AssertEquals(t, "JsonStrack returns correct unprintable message:", "Unprintable", traceData.Panic)

}
