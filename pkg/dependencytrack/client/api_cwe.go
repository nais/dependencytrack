/*
OWASP Dependency-Track

REST API of OWASP Dependency-Track

API version: 4.13.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)


type CweAPI interface {

	/*
	GetCwe Returns a specific CWE

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param cweId The CWE ID of the CWE to retrieve
	@return ApiGetCweRequest
	*/
	GetCwe(ctx context.Context, cweId int32) ApiGetCweRequest

	// GetCweExecute executes the request
	//  @return Cwe
	GetCweExecute(r ApiGetCweRequest) (*Cwe, *http.Response, error)

	/*
	GetCwes Returns a list of all CWEs

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiGetCwesRequest
	*/
	GetCwes(ctx context.Context) ApiGetCwesRequest

	// GetCwesExecute executes the request
	//  @return []Cwe
	GetCwesExecute(r ApiGetCwesRequest) ([]Cwe, *http.Response, error)
}

// CweAPIService CweAPI service
type CweAPIService service

type ApiGetCweRequest struct {
	ctx context.Context
	ApiService CweAPI
	cweId int32
}

func (r ApiGetCweRequest) Execute() (*Cwe, *http.Response, error) {
	return r.ApiService.GetCweExecute(r)
}

/*
GetCwe Returns a specific CWE

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param cweId The CWE ID of the CWE to retrieve
 @return ApiGetCweRequest
*/
func (a *CweAPIService) GetCwe(ctx context.Context, cweId int32) ApiGetCweRequest {
	return ApiGetCweRequest{
		ApiService: a,
		ctx: ctx,
		cweId: cweId,
	}
}

// Execute executes the request
//  @return Cwe
func (a *CweAPIService) GetCweExecute(r ApiGetCweRequest) (*Cwe, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Cwe
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CweAPIService.GetCwe")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/cwe/{cweId}"
	localVarPath = strings.Replace(localVarPath, "{"+"cweId"+"}", url.PathEscape(parameterValueToString(r.cweId, "cweId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["ApiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["X-Api-Key"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetCwesRequest struct {
	ctx context.Context
	ApiService CweAPI
	pageNumber *string
	pageSize *string
	offset *string
	limit *string
	sortName *string
	sortOrder *string
}

// The page to return. To be used in conjunction with &lt;code&gt;pageSize&lt;/code&gt;.
func (r ApiGetCwesRequest) PageNumber(pageNumber string) ApiGetCwesRequest {
	r.pageNumber = &pageNumber
	return r
}

// Number of elements to return per page. To be used in conjunction with &lt;code&gt;pageNumber&lt;/code&gt;.
func (r ApiGetCwesRequest) PageSize(pageSize string) ApiGetCwesRequest {
	r.pageSize = &pageSize
	return r
}

// Offset to start returning elements from. To be used in conjunction with &lt;code&gt;limit&lt;/code&gt;.
func (r ApiGetCwesRequest) Offset(offset string) ApiGetCwesRequest {
	r.offset = &offset
	return r
}

// Number of elements to return per page. To be used in conjunction with &lt;code&gt;offset&lt;/code&gt;.
func (r ApiGetCwesRequest) Limit(limit string) ApiGetCwesRequest {
	r.limit = &limit
	return r
}

// Name of the resource field to sort on.
func (r ApiGetCwesRequest) SortName(sortName string) ApiGetCwesRequest {
	r.sortName = &sortName
	return r
}

// Ordering of items when sorting with &lt;code&gt;sortName&lt;/code&gt;.
func (r ApiGetCwesRequest) SortOrder(sortOrder string) ApiGetCwesRequest {
	r.sortOrder = &sortOrder
	return r
}

func (r ApiGetCwesRequest) Execute() ([]Cwe, *http.Response, error) {
	return r.ApiService.GetCwesExecute(r)
}

/*
GetCwes Returns a list of all CWEs

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiGetCwesRequest
*/
func (a *CweAPIService) GetCwes(ctx context.Context) ApiGetCwesRequest {
	return ApiGetCwesRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return []Cwe
func (a *CweAPIService) GetCwesExecute(r ApiGetCwesRequest) ([]Cwe, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  []Cwe
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CweAPIService.GetCwes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/cwe"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.pageNumber != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNumber", r.pageNumber, "form", "")
	} else {
		var defaultValue string = "1"
		r.pageNumber = &defaultValue
	}
	if r.pageSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageSize", r.pageSize, "form", "")
	} else {
		var defaultValue string = "100"
		r.pageSize = &defaultValue
	}
	if r.offset != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "offset", r.offset, "form", "")
	}
	if r.limit != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "limit", r.limit, "form", "")
	}
	if r.sortName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "sortName", r.sortName, "form", "")
	}
	if r.sortOrder != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "sortOrder", r.sortOrder, "form", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["ApiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["X-Api-Key"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
