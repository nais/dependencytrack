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


type LicenseGroupAPI interface {

	/*
	AddLicenseToLicenseGroup Adds the license to the specified license group.

	<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param uuid A valid license group
	@param licenseUuid A valid license
	@return ApiAddLicenseToLicenseGroupRequest
	*/
	AddLicenseToLicenseGroup(ctx context.Context, uuid string, licenseUuid string) ApiAddLicenseToLicenseGroupRequest

	// AddLicenseToLicenseGroupExecute executes the request
	//  @return LicenseGroup
	AddLicenseToLicenseGroupExecute(r ApiAddLicenseToLicenseGroupRequest) (*LicenseGroup, *http.Response, error)

	/*
	CreateLicenseGroup Creates a new license group

	<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiCreateLicenseGroupRequest
	*/
	CreateLicenseGroup(ctx context.Context) ApiCreateLicenseGroupRequest

	// CreateLicenseGroupExecute executes the request
	//  @return LicenseGroup
	CreateLicenseGroupExecute(r ApiCreateLicenseGroupRequest) (*LicenseGroup, *http.Response, error)

	/*
	DeleteLicenseGroup Deletes a license group

	<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param uuid The UUID of the license group to delete
	@return ApiDeleteLicenseGroupRequest
	*/
	DeleteLicenseGroup(ctx context.Context, uuid string) ApiDeleteLicenseGroupRequest

	// DeleteLicenseGroupExecute executes the request
	DeleteLicenseGroupExecute(r ApiDeleteLicenseGroupRequest) (*http.Response, error)

	/*
	GetLicenseGroup Returns a specific license group

	<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param uuid The UUID of the license group to retrieve
	@return ApiGetLicenseGroupRequest
	*/
	GetLicenseGroup(ctx context.Context, uuid string) ApiGetLicenseGroupRequest

	// GetLicenseGroupExecute executes the request
	//  @return LicenseGroup
	GetLicenseGroupExecute(r ApiGetLicenseGroupRequest) (*LicenseGroup, *http.Response, error)

	/*
	GetLicenseGroups Returns a list of all license groups

	<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiGetLicenseGroupsRequest
	*/
	GetLicenseGroups(ctx context.Context) ApiGetLicenseGroupsRequest

	// GetLicenseGroupsExecute executes the request
	//  @return []LicenseGroup
	GetLicenseGroupsExecute(r ApiGetLicenseGroupsRequest) ([]LicenseGroup, *http.Response, error)

	/*
	RemoveLicenseFromLicenseGroup Removes the license from the license group.

	<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param uuid A valid license group
	@param licenseUuid A valid license
	@return ApiRemoveLicenseFromLicenseGroupRequest
	*/
	RemoveLicenseFromLicenseGroup(ctx context.Context, uuid string, licenseUuid string) ApiRemoveLicenseFromLicenseGroupRequest

	// RemoveLicenseFromLicenseGroupExecute executes the request
	//  @return LicenseGroup
	RemoveLicenseFromLicenseGroupExecute(r ApiRemoveLicenseFromLicenseGroupRequest) (*LicenseGroup, *http.Response, error)

	/*
	UpdateLicenseGroup Updates a license group

	<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiUpdateLicenseGroupRequest
	*/
	UpdateLicenseGroup(ctx context.Context) ApiUpdateLicenseGroupRequest

	// UpdateLicenseGroupExecute executes the request
	//  @return LicenseGroup
	UpdateLicenseGroupExecute(r ApiUpdateLicenseGroupRequest) (*LicenseGroup, *http.Response, error)
}

// LicenseGroupAPIService LicenseGroupAPI service
type LicenseGroupAPIService service

type ApiAddLicenseToLicenseGroupRequest struct {
	ctx context.Context
	ApiService LicenseGroupAPI
	uuid string
	licenseUuid string
}

func (r ApiAddLicenseToLicenseGroupRequest) Execute() (*LicenseGroup, *http.Response, error) {
	return r.ApiService.AddLicenseToLicenseGroupExecute(r)
}

/*
AddLicenseToLicenseGroup Adds the license to the specified license group.

<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param uuid A valid license group
 @param licenseUuid A valid license
 @return ApiAddLicenseToLicenseGroupRequest
*/
func (a *LicenseGroupAPIService) AddLicenseToLicenseGroup(ctx context.Context, uuid string, licenseUuid string) ApiAddLicenseToLicenseGroupRequest {
	return ApiAddLicenseToLicenseGroupRequest{
		ApiService: a,
		ctx: ctx,
		uuid: uuid,
		licenseUuid: licenseUuid,
	}
}

// Execute executes the request
//  @return LicenseGroup
func (a *LicenseGroupAPIService) AddLicenseToLicenseGroupExecute(r ApiAddLicenseToLicenseGroupRequest) (*LicenseGroup, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *LicenseGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LicenseGroupAPIService.AddLicenseToLicenseGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/licenseGroup/{uuid}/license/{licenseUuid}"
	localVarPath = strings.Replace(localVarPath, "{"+"uuid"+"}", url.PathEscape(parameterValueToString(r.uuid, "uuid")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"licenseUuid"+"}", url.PathEscape(parameterValueToString(r.licenseUuid, "licenseUuid")), -1)

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

type ApiCreateLicenseGroupRequest struct {
	ctx context.Context
	ApiService LicenseGroupAPI
	licenseGroup *LicenseGroup
}

func (r ApiCreateLicenseGroupRequest) LicenseGroup(licenseGroup LicenseGroup) ApiCreateLicenseGroupRequest {
	r.licenseGroup = &licenseGroup
	return r
}

func (r ApiCreateLicenseGroupRequest) Execute() (*LicenseGroup, *http.Response, error) {
	return r.ApiService.CreateLicenseGroupExecute(r)
}

/*
CreateLicenseGroup Creates a new license group

<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiCreateLicenseGroupRequest
*/
func (a *LicenseGroupAPIService) CreateLicenseGroup(ctx context.Context) ApiCreateLicenseGroupRequest {
	return ApiCreateLicenseGroupRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return LicenseGroup
func (a *LicenseGroupAPIService) CreateLicenseGroupExecute(r ApiCreateLicenseGroupRequest) (*LicenseGroup, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *LicenseGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LicenseGroupAPIService.CreateLicenseGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/licenseGroup"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	// body params
	localVarPostBody = r.licenseGroup
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

type ApiDeleteLicenseGroupRequest struct {
	ctx context.Context
	ApiService LicenseGroupAPI
	uuid string
}

func (r ApiDeleteLicenseGroupRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteLicenseGroupExecute(r)
}

/*
DeleteLicenseGroup Deletes a license group

<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param uuid The UUID of the license group to delete
 @return ApiDeleteLicenseGroupRequest
*/
func (a *LicenseGroupAPIService) DeleteLicenseGroup(ctx context.Context, uuid string) ApiDeleteLicenseGroupRequest {
	return ApiDeleteLicenseGroupRequest{
		ApiService: a,
		ctx: ctx,
		uuid: uuid,
	}
}

// Execute executes the request
func (a *LicenseGroupAPIService) DeleteLicenseGroupExecute(r ApiDeleteLicenseGroupRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LicenseGroupAPIService.DeleteLicenseGroup")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/licenseGroup/{uuid}"
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
	localVarHTTPHeaderAccepts := []string{}

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
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiGetLicenseGroupRequest struct {
	ctx context.Context
	ApiService LicenseGroupAPI
	uuid string
}

func (r ApiGetLicenseGroupRequest) Execute() (*LicenseGroup, *http.Response, error) {
	return r.ApiService.GetLicenseGroupExecute(r)
}

/*
GetLicenseGroup Returns a specific license group

<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param uuid The UUID of the license group to retrieve
 @return ApiGetLicenseGroupRequest
*/
func (a *LicenseGroupAPIService) GetLicenseGroup(ctx context.Context, uuid string) ApiGetLicenseGroupRequest {
	return ApiGetLicenseGroupRequest{
		ApiService: a,
		ctx: ctx,
		uuid: uuid,
	}
}

// Execute executes the request
//  @return LicenseGroup
func (a *LicenseGroupAPIService) GetLicenseGroupExecute(r ApiGetLicenseGroupRequest) (*LicenseGroup, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *LicenseGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LicenseGroupAPIService.GetLicenseGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/licenseGroup/{uuid}"
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

type ApiGetLicenseGroupsRequest struct {
	ctx context.Context
	ApiService LicenseGroupAPI
	pageNumber *string
	pageSize *string
	offset *string
	limit *string
	sortName *string
	sortOrder *string
}

// The page to return. To be used in conjunction with &lt;code&gt;pageSize&lt;/code&gt;.
func (r ApiGetLicenseGroupsRequest) PageNumber(pageNumber string) ApiGetLicenseGroupsRequest {
	r.pageNumber = &pageNumber
	return r
}

// Number of elements to return per page. To be used in conjunction with &lt;code&gt;pageNumber&lt;/code&gt;.
func (r ApiGetLicenseGroupsRequest) PageSize(pageSize string) ApiGetLicenseGroupsRequest {
	r.pageSize = &pageSize
	return r
}

// Offset to start returning elements from. To be used in conjunction with &lt;code&gt;limit&lt;/code&gt;.
func (r ApiGetLicenseGroupsRequest) Offset(offset string) ApiGetLicenseGroupsRequest {
	r.offset = &offset
	return r
}

// Number of elements to return per page. To be used in conjunction with &lt;code&gt;offset&lt;/code&gt;.
func (r ApiGetLicenseGroupsRequest) Limit(limit string) ApiGetLicenseGroupsRequest {
	r.limit = &limit
	return r
}

// Name of the resource field to sort on.
func (r ApiGetLicenseGroupsRequest) SortName(sortName string) ApiGetLicenseGroupsRequest {
	r.sortName = &sortName
	return r
}

// Ordering of items when sorting with &lt;code&gt;sortName&lt;/code&gt;.
func (r ApiGetLicenseGroupsRequest) SortOrder(sortOrder string) ApiGetLicenseGroupsRequest {
	r.sortOrder = &sortOrder
	return r
}

func (r ApiGetLicenseGroupsRequest) Execute() ([]LicenseGroup, *http.Response, error) {
	return r.ApiService.GetLicenseGroupsExecute(r)
}

/*
GetLicenseGroups Returns a list of all license groups

<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiGetLicenseGroupsRequest
*/
func (a *LicenseGroupAPIService) GetLicenseGroups(ctx context.Context) ApiGetLicenseGroupsRequest {
	return ApiGetLicenseGroupsRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return []LicenseGroup
func (a *LicenseGroupAPIService) GetLicenseGroupsExecute(r ApiGetLicenseGroupsRequest) ([]LicenseGroup, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  []LicenseGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LicenseGroupAPIService.GetLicenseGroups")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/licenseGroup"

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

type ApiRemoveLicenseFromLicenseGroupRequest struct {
	ctx context.Context
	ApiService LicenseGroupAPI
	uuid string
	licenseUuid string
}

func (r ApiRemoveLicenseFromLicenseGroupRequest) Execute() (*LicenseGroup, *http.Response, error) {
	return r.ApiService.RemoveLicenseFromLicenseGroupExecute(r)
}

/*
RemoveLicenseFromLicenseGroup Removes the license from the license group.

<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param uuid A valid license group
 @param licenseUuid A valid license
 @return ApiRemoveLicenseFromLicenseGroupRequest
*/
func (a *LicenseGroupAPIService) RemoveLicenseFromLicenseGroup(ctx context.Context, uuid string, licenseUuid string) ApiRemoveLicenseFromLicenseGroupRequest {
	return ApiRemoveLicenseFromLicenseGroupRequest{
		ApiService: a,
		ctx: ctx,
		uuid: uuid,
		licenseUuid: licenseUuid,
	}
}

// Execute executes the request
//  @return LicenseGroup
func (a *LicenseGroupAPIService) RemoveLicenseFromLicenseGroupExecute(r ApiRemoveLicenseFromLicenseGroupRequest) (*LicenseGroup, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *LicenseGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LicenseGroupAPIService.RemoveLicenseFromLicenseGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/licenseGroup/{uuid}/license/{licenseUuid}"
	localVarPath = strings.Replace(localVarPath, "{"+"uuid"+"}", url.PathEscape(parameterValueToString(r.uuid, "uuid")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"licenseUuid"+"}", url.PathEscape(parameterValueToString(r.licenseUuid, "licenseUuid")), -1)

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

type ApiUpdateLicenseGroupRequest struct {
	ctx context.Context
	ApiService LicenseGroupAPI
	licenseGroup *LicenseGroup
}

func (r ApiUpdateLicenseGroupRequest) LicenseGroup(licenseGroup LicenseGroup) ApiUpdateLicenseGroupRequest {
	r.licenseGroup = &licenseGroup
	return r
}

func (r ApiUpdateLicenseGroupRequest) Execute() (*LicenseGroup, *http.Response, error) {
	return r.ApiService.UpdateLicenseGroupExecute(r)
}

/*
UpdateLicenseGroup Updates a license group

<p>Requires permission <strong>POLICY_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiUpdateLicenseGroupRequest
*/
func (a *LicenseGroupAPIService) UpdateLicenseGroup(ctx context.Context) ApiUpdateLicenseGroupRequest {
	return ApiUpdateLicenseGroupRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return LicenseGroup
func (a *LicenseGroupAPIService) UpdateLicenseGroupExecute(r ApiUpdateLicenseGroupRequest) (*LicenseGroup, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *LicenseGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LicenseGroupAPIService.UpdateLicenseGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/licenseGroup"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	// body params
	localVarPostBody = r.licenseGroup
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
