package problemdetails

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		problemType string
		title       string
		statusCode  int
		detail      string
		instance    string
	}
	type test struct {
		name     string
		args     args
		expected *ProblemDetails
	}

	tests := []test{
		{
			name: "status_http_500_blank_type",
			args: args{
				problemType: "",
				title:       "",
				statusCode:  http.StatusInternalServerError,
				detail:      "",
				instance:    "",
			},
			expected: &ProblemDetails{
				Type:     defaultProblemType,
				Title:    http.StatusText(http.StatusInternalServerError),
				Status:   http.StatusInternalServerError,
				Detail:   "",
				Instance: "",
			},
		},
		{
			name: "status_http_404",
			args: args{
				problemType: "https://some-domain.com/not_found",
				title:       http.StatusText(http.StatusNotFound),
				statusCode:  http.StatusNotFound,
				detail:      "Object with id 1 was not found",
				instance:    "https://api.some-domain.com/xxx/1",
			},
			expected: &ProblemDetails{
				Type:     "https://some-domain.com/not_found",
				Title:    http.StatusText(http.StatusNotFound),
				Status:   http.StatusNotFound,
				Detail:   "Object with id 1 was not found",
				Instance: "https://api.some-domain.com/xxx/1",
			},
		},
		{
			name: "status_http_500_relative_uris",
			args: args{
				problemType: "/internal_server_error",
				title:       http.StatusText(http.StatusInternalServerError),
				statusCode:  http.StatusInternalServerError,
				detail:      "There was an error",
				instance:    "example-1",
			},
			expected: &ProblemDetails{
				Type:     "/internal_server_error",
				Title:    http.StatusText(http.StatusInternalServerError),
				Status:   http.StatusInternalServerError,
				Detail:   "There was an error",
				Instance: "example-1",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := New(test.args.problemType, test.args.title, test.args.statusCode, test.args.detail, test.args.instance)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("[actual]= %+v, [expected]= %+v", actual, test.expected)
			}
		})
	}
}

func TestFromHTTPStatus(t *testing.T) {
	type args struct {
		statusCode int
	}
	type test struct {
		name     string
		args     args
		expected *ProblemDetails
	}

	tests := []test{
		{
			name: "status_http_400",
			args: args{
				statusCode: http.StatusBadRequest,
			},
			expected: &ProblemDetails{
				Type:     defaultProblemType,
				Title:    http.StatusText(http.StatusBadRequest),
				Status:   http.StatusBadRequest,
				Detail:   "",
				Instance: "",
			},
		},
		{
			name: "status_http_500",
			args: args{
				statusCode: http.StatusInternalServerError,
			},
			expected: &ProblemDetails{
				Type:     defaultProblemType,
				Title:    http.StatusText(http.StatusInternalServerError),
				Status:   http.StatusInternalServerError,
				Detail:   "",
				Instance: "",
			},
		},
		{
			name: "unknown_status",
			args: args{
				statusCode: 0,
			},
			expected: &ProblemDetails{
				Type:     defaultProblemType,
				Title:    "",
				Status:   0,
				Detail:   "",
				Instance: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FromHTTPStatus(test.args.statusCode)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("[actual]= %+v, [expected]= %+v", actual, test.expected)
			}
		})
	}
}

func TestJSONMarshal(t *testing.T) {
	type test struct {
		name     string
		args     *ProblemDetails
		expected string
	}

	tests := []test{
		{
			name: "status_http_500_blank_type",
			args: &ProblemDetails{
				Type:     defaultProblemType,
				Title:    http.StatusText(http.StatusInternalServerError),
				Status:   http.StatusInternalServerError,
				Detail:   "",
				Instance: "",
			},
			expected: fmt.Sprintf(`{"type":"%s","title":"%s","status":%v}`, defaultProblemType, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
		},
		{
			name: "status_http_400_error",
			args: &ProblemDetails{
				Type:     "https://some-domain.com/validation-failed",
				Title:    "Validation error",
				Status:   http.StatusBadRequest,
				Detail:   "Your request parameters didn't validate",
				Instance: "https://api.some-domain.com/example",
			},
			expected: fmt.Sprintf(`{"type":"https://some-domain.com/validation-failed","title":"Validation error","status":%v,"detail":"Your request parameters didn't validate","instance":"https://api.some-domain.com/example"}`, http.StatusBadRequest),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := json.Marshal(test.args)
			if err != nil {
				t.Fatalf("error during json marshal %v", err)
			}
			if test.expected != string(actual) {
				t.Fatalf("[actual]= %s, [expected]= %+v", actual, test.expected)
			}
		})
	}
}

func TestXMLMarshal(t *testing.T) {
	type test struct {
		name     string
		args     *ProblemDetails
		expected string
	}

	tests := []test{
		{
			name: "status_http_500_blank_type",
			args: &ProblemDetails{
				Type:     defaultProblemType,
				Title:    http.StatusText(http.StatusInternalServerError),
				Status:   http.StatusInternalServerError,
				Detail:   "",
				Instance: "",
			},
			expected: fmt.Sprintf(`<problem xmlns="urn:ietf:rfc:7807"><type>%s</type><title>%s</title><status>%v</status></problem>`, defaultProblemType, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
		},
		{
			name: "status_http_400_error",
			args: &ProblemDetails{
				Type:     "https://some-domain.com/validation-failed",
				Title:    "Validation error",
				Status:   http.StatusBadRequest,
				Detail:   "Your request parameters didn't validate",
				Instance: "https://api.some-domain.com/example",
			},
			expected: fmt.Sprintf(`<problem xmlns="urn:ietf:rfc:7807"><type>https://some-domain.com/validation-failed</type><title>Validation error</title><status>%v</status><detail>Your request parameters didn&#39;t validate</detail><instance>https://api.some-domain.com/example</instance></problem>`, http.StatusBadRequest),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := xml.Marshal(test.args)
			if err != nil {
				t.Fatalf("error during xml marshal %v", err)
			}
			if test.expected != string(actual) {
				t.Fatalf("[actual]= %s, [expected]= %+v", actual, test.expected)
			}
		})
	}
}
