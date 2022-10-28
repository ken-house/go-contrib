package requester

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
)

type RequestClient interface {
	PostJSON(ctx context.Context, endpoint string, payloadData map[string]interface{}, responseStruct interface{}, querystring map[string]string) (*http.Response, error)
	Post(ctx context.Context, endpoint string, payloadData map[string]interface{}, responseStruct interface{}, querystring map[string]string) (*http.Response, error)
	PostFiles(ctx context.Context, endpoint string, payload io.Reader, responseStruct interface{}, querystring map[string]string, files []string) (*http.Response, error)
	PostXML(ctx context.Context, endpoint string, xml string, responseStruct interface{}, querystring map[string]string) (*http.Response, error)
	Get(ctx context.Context, endpoint string, responseStruct interface{}, querystring map[string]string) (*http.Response, error)
	GetXML(ctx context.Context, endpoint string, responseStruct interface{}, querystring map[string]string) (*http.Response, error)
	Do(ctx context.Context, ar *APIRequest, responseStruct interface{}, options ...interface{}) (*http.Response, error)
	ReadRawResponse(response *http.Response, responseStruct interface{}) (*http.Response, error)
	ReadJSONResponse(response *http.Response, responseStruct interface{}) (*http.Response, error)
}

type requestClient struct {
	Base              string
	Client            *http.Client
	CACert            []byte
	SslVerify         bool
	CustomerHeaderArr map[string]string
}

func NewRequestClient(base string, sslVerify bool, customerHeaderArr map[string]string, cACert string) RequestClient {
	return &requestClient{
		Base:              base,
		CustomerHeaderArr: customerHeaderArr,
		SslVerify:         sslVerify,
		CACert:            []byte(cACert),
		Client:            http.DefaultClient,
	}
}

func (r *requestClient) PostJSON(ctx context.Context, endpoint string, payloadData map[string]interface{}, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	data, err := json.Marshal(payloadData)
	if err != nil {
		return nil, errors.New("json Marshal failed")
	}
	payload := bytes.NewBufferString(string(data))
	ar := NewAPIRequest("POST", endpoint, payload)
	ar.SetHeader("Content-Type", "application/json")
	return r.Do(ctx, ar, &responseStruct, querystring)
}

func (r *requestClient) Post(ctx context.Context, endpoint string, payloadData map[string]interface{}, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	data := url.Values{}
	for key, value := range payloadData {
		data.Set(key, cast.ToString(value))
	}
	payload := bytes.NewBufferString(data.Encode())
	ar := NewAPIRequest("POST", endpoint, payload)
	ar.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return r.Do(ctx, ar, &responseStruct, querystring)
}

func (r *requestClient) PostFiles(ctx context.Context, endpoint string, payload io.Reader, responseStruct interface{}, querystring map[string]string, files []string) (*http.Response, error) {
	ar := NewAPIRequest("POST", endpoint, payload)
	return r.Do(ctx, ar, &responseStruct, querystring, files)
}

func (r *requestClient) PostXML(ctx context.Context, endpoint string, xml string, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	payload := bytes.NewBuffer([]byte(xml))
	ar := NewAPIRequest("POST", endpoint, payload)
	ar.SetHeader("Content-Type", "application/xml")
	return r.Do(ctx, ar, &responseStruct, querystring)
}

func (r *requestClient) GetXML(ctx context.Context, endpoint string, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("GET", endpoint, nil)
	ar.SetHeader("Content-Type", "application/xml")
	return r.Do(ctx, ar, responseStruct, querystring)
}

func (r *requestClient) Get(ctx context.Context, endpoint string, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("GET", endpoint, nil)
	return r.Do(ctx, ar, responseStruct, querystring)
}

func (r *requestClient) Do(ctx context.Context, ar *APIRequest, responseStruct interface{}, options ...interface{}) (*http.Response, error) {
	if !strings.HasSuffix(ar.Endpoint, "/") && ar.Method != "POST" {
		ar.Endpoint += "/"
	}

	fileUpload := false
	var files []string
	URL, err := url.Parse(r.Base + ar.Endpoint)

	if err != nil {
		return nil, err
	}

	for _, o := range options {
		switch v := o.(type) {
		case map[string]string:
			querystring := make(url.Values)
			for key, val := range v {
				querystring.Set(key, val)
			}

			URL.RawQuery = querystring.Encode()
		case []string:
			fileUpload = true
			files = v
		}
	}

	var req *http.Request
	if fileUpload {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for _, file := range files {
			fileData, err := os.Open(file)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}

			part, err := writer.CreateFormFile("file", filepath.Base(file))
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			if _, err = io.Copy(part, fileData); err != nil {
				return nil, err
			}
			defer fileData.Close()
		}
		var params map[string]string
		json.NewDecoder(ar.Payload).Decode(&params)
		for key, val := range params {
			if err = writer.WriteField(key, val); err != nil {
				return nil, err
			}
		}
		if err = writer.Close(); err != nil {
			return nil, err
		}
		req, err = http.NewRequest(ar.Method, URL.String(), body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest(ar.Method, URL.String(), ar.Payload)
		if err != nil {
			return nil, err
		}
	}

	if r.CustomerHeaderArr != nil {
		for k, v := range r.CustomerHeaderArr {
			req.Header.Add(k, v)
		}
	}

	for k := range ar.Headers {
		req.Header.Add(k, ar.Headers.Get(k))
	}

	if response, err := r.Client.Do(req); err != nil {
		return nil, err
	} else {
		if v := ctx.Value("debug"); v != nil {
			dump, err := httputil.DumpResponse(response, true)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("DEBUG %q\n", dump)
		}
		errorText := response.Header.Get("X-Error")
		if errorText != "" {
			return nil, errors.New(errorText)
		}
		switch responseStruct.(type) {
		case *string:
			return r.ReadRawResponse(response, responseStruct)
		default:
			return r.ReadJSONResponse(response, responseStruct)
		}
	}
}

func (r *requestClient) ReadRawResponse(response *http.Response, responseStruct interface{}) (*http.Response, error) {
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if str, ok := responseStruct.(*string); ok {
		*str = string(content)
	} else {
		return nil, fmt.Errorf("Could not cast responseStruct to *string")
	}

	return response, nil
}

func (r *requestClient) ReadJSONResponse(response *http.Response, responseStruct interface{}) (*http.Response, error) {
	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(responseStruct)
	return response, nil
}
