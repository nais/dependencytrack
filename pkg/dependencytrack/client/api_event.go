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


type EventAPI interface {

	/*
	IsTokenBeingProcessed1 Determines if there are any tasks associated with the token that are being processed, or in the queue to be processed.

	<p>
  This endpoint is intended to be used in conjunction with other API calls which return a token for asynchronous tasks.
  The token can then be queried using this endpoint to determine if the task is complete:
  <ul>
    <li>A value of <code>true</code> indicates processing is occurring.</li>
    <li>A value of <code>false</code> indicates that no processing is occurring for the specified token.</li>
  </ul>
  However, a value of <code>false</code> also does not confirm the token is valid,
  only that no processing is associated with the specified token.
</p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param uuid The UUID of the token to query
	@return ApiIsTokenBeingProcessed1Request
	*/
	IsTokenBeingProcessed1(ctx context.Context, uuid string) ApiIsTokenBeingProcessed1Request

	// IsTokenBeingProcessed1Execute executes the request
	//  @return IsTokenBeingProcessedResponse
	IsTokenBeingProcessed1Execute(r ApiIsTokenBeingProcessed1Request) (*IsTokenBeingProcessedResponse, *http.Response, error)
}

// EventAPIService EventAPI service
type EventAPIService service

type ApiIsTokenBeingProcessed1Request struct {
	ctx context.Context
	ApiService EventAPI
	uuid string
}

func (r ApiIsTokenBeingProcessed1Request) Execute() (*IsTokenBeingProcessedResponse, *http.Response, error) {
	return r.ApiService.IsTokenBeingProcessed1Execute(r)
}

/*
IsTokenBeingProcessed1 Determines if there are any tasks associated with the token that are being processed, or in the queue to be processed.

<p>
  This endpoint is intended to be used in conjunction with other API calls which return a token for asynchronous tasks.
  The token can then be queried using this endpoint to determine if the task is complete:
  <ul>
    <li>A value of <code>true</code> indicates processing is occurring.</li>
    <li>A value of <code>false</code> indicates that no processing is occurring for the specified token.</li>
  </ul>
  However, a value of <code>false</code> also does not confirm the token is valid,
  only that no processing is associated with the specified token.
</p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param uuid The UUID of the token to query
 @return ApiIsTokenBeingProcessed1Request
*/
func (a *EventAPIService) IsTokenBeingProcessed1(ctx context.Context, uuid string) ApiIsTokenBeingProcessed1Request {
	return ApiIsTokenBeingProcessed1Request{
		ApiService: a,
		ctx: ctx,
		uuid: uuid,
	}
}

// Execute executes the request
//  @return IsTokenBeingProcessedResponse
func (a *EventAPIService) IsTokenBeingProcessed1Execute(r ApiIsTokenBeingProcessed1Request) (*IsTokenBeingProcessedResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *IsTokenBeingProcessedResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EventAPIService.IsTokenBeingProcessed1")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/event/token/{uuid}"
	localVarPath = strings.Replace(localVarPath, "{"+"uuid"+"}", url.PathEscape(parameterValueToString(r.uuid, "uuid")), -1)

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
