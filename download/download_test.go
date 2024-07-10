package download

import (
	"net/http"
	"reflect"
	"testing"
)

func TestIsValidEndpoint(t *testing.T) {
	testCases := []struct {
		name     string
		endpoint string
		want     bool
	}{
		{"Valid HTTP URL", "http://example.com", true},
		{"Valid HTTPS URL", "https://example.com", true},
		{"Invalid URL", "htp:/example.com", false},
		{"Invalid URL", "htp:/example.ctdm", false},
		{"Empty string", "", false},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			got := isValidEndpoint(tc.endpoint)
			if got != tc.want {
				t.Errorf("isValidEndpoint(%q) = %v; want %v", tc.endpoint, got, tc.want)
			}
		})
	}
}

func TestReadResponse(t *testing.T) {
	res, err := http.Get("https://raw.githubusercontent.com/golang/go/master/README.md")

	if err != nil {
		t.Fatalf("http.Get() error = %v", err)
	}

	defer res.Body.Close()

	got, responseErr := readResponse(res)

	if responseErr != nil {
		t.Errorf("readResponse() error = %v", responseErr)
		return
	}

	if len(got) == 0 {
		t.Errorf("readResponse() returned empty content")
	}

}

func TestReverseSortUrlObjects(t *testing.T) {

	var testSlice = []UrlObject{
		{Id: 1, Content: "Content 1"},
		{Id: 2, Content: "Content 2"},
		{Id: 3, Content: "Content 3"},
		{Id: 4, Content: "Content 4"},
	}
	want := []UrlObject{
		{Id: 4, Content: "Content 4"},
		{Id: 3, Content: "Content 3"},
		{Id: 2, Content: "Content 2"},
		{Id: 1, Content: "Content 1"},
	}

	reverseSortUrlObjects(testSlice)

	if !reflect.DeepEqual(testSlice, want) {
		t.Errorf("Test failed, testSlice is not the same like")
	}

}

func TestConcatenateContent(t *testing.T) {

	TesturlObjects := []UrlObject{
		{Id: 1, Content: "Content 1"},
		{Id: 2, Content: "Content 2"},
		{Id: 3, Content: "Content 3"},
		{Id: 4, Content: "Content 4"},
	}

	result := concatenateContent(TesturlObjects)

	expected := "Content 1 " + "Content 2 " + "Content 3 " + "Content 4"
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

func TestReturnContentOrFail_Fail(t *testing.T) {

	testUrlsSlice := []string{
		"https://raw.githubusercontent.com/GoogleContainerTools/distroless/main/java/README.md",
		"https://raw.githubusercontent.com/golang/go/master/README.md",
		"https://pkg.go.dev/github.com/posener/goreadme",
		"http://thisisanabsolutellyinvalidurl.org.ee",
		"https://pkg.go.dev/go.jpap.org/godoc-readme-gen",
		"https://github.com/golang/example/blob/master/README.md",
	}

	result, err := ReturnContentOrFail(testUrlsSlice)

	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Unexpected length of  result : %v", len(result))
	}

}

func TestReturnContentOrFail_Pass(t *testing.T) {

	testUrlsSlice := []string{
		"https://raw.githubusercontent.com/GoogleContainerTools/distroless/main/java/README.md",
		"https://raw.githubusercontent.com/golang/go/master/README.md",
		"https://pkg.go.dev/github.com/posener/goreadme",
		"https://pkg.go.dev/go.jpap.org/godoc-readme-gen",
		"https://github.com/golang/example/blob/master/README.md",
	}

	result, err := ReturnContentOrFail(testUrlsSlice)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) == 0 {
		t.Errorf("unexpected content, got %v", len(result))
	}

}
