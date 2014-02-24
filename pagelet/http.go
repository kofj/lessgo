package pagelet

import (
    "bytes"
    "fmt"
    "net/http"
    "sort"
    "strconv"
    "strings"
)

type Request struct {
    *http.Request
    ContentType    string
    AcceptLanguage []AcceptLanguage
    Locale         string
}

type Response struct {
    Status      int
    ContentType string
    Out         http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
    return &Response{Out: w}
}

func NewRequest(r *http.Request) *Request {
    return &Request{
        Request:        r,
        ContentType:    "text/html",
        AcceptLanguage: resolveAcceptLanguage(r),
        Locale:         "",
    }
}

func (resp *Response) WriteHeader(status int, ctype string) {

    if resp.Status == 0 {
        resp.Status = status
        resp.Out.WriteHeader(resp.Status)
    }

    if resp.ContentType == "" {
        resp.ContentType = ctype
        resp.Out.Header().Set("Content-Type", resp.ContentType)
    }
}

// Resolve the Accept-Language header value.
//
// The results are sorted using the quality defined in the header for each language range with the
// most qualified language range as the first element in the slice.
//
// See the HTTP header fields specification
// (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4) for more details.
func resolveAcceptLanguage(r *http.Request) AcceptLanguages {

    als := strings.Split(r.Header.Get("Accept-Language"), ",")

    acceptLanguages := make(AcceptLanguages, len(als))

    for i, v := range als {

        if v2 := strings.Split(v, ";q="); len(v2) == 2 {
            quality, err := strconv.ParseFloat(v2[1], 32)
            if err != nil {
                quality = 1
            }
            acceptLanguages[i] = AcceptLanguage{v2[0], float32(quality)}
        } else if len(v2) == 1 {
            acceptLanguages[i] = AcceptLanguage{v, 1}
        }
    }

    sort.Sort(acceptLanguages)
    return acceptLanguages
}

// A single language from the Accept-Language HTTP header.
type AcceptLanguage struct {
    Language string
    Quality  float32
}

// A collection of sortable AcceptLanguage instances.
type AcceptLanguages []AcceptLanguage

func (al AcceptLanguages) Len() int           { return len(al) }
func (al AcceptLanguages) Swap(i, j int)      { al[i], al[j] = al[j], al[i] }
func (al AcceptLanguages) Less(i, j int) bool { return al[i].Quality > al[j].Quality }
func (al AcceptLanguages) String() string {
    output := bytes.NewBufferString("")
    for i, language := range al {
        output.WriteString(fmt.Sprintf("%s (%1.1f)", language.Language, language.Quality))
        if i != len(al)-1 {
            output.WriteString(", ")
        }
    }
    return output.String()
}