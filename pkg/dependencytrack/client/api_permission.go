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


type PermissionAPI interface {

	/*
	AddPermissionToTeam Adds the permission to the specified team.

	<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param uuid A valid team uuid
	@param permission A valid permission
	@return ApiAddPermissionToTeamRequest
	*/
	AddPermissionToTeam(ctx context.Context, uuid string, permission string) ApiAddPermissionToTeamRequest

	// AddPermissionToTeamExecute executes the request
	//  @return Team
	AddPermissionToTeamExecute(r ApiAddPermissionToTeamRequest) (*Team, *http.Response, error)

	/*
	AddPermissionToUser Adds the permission to the specified username.

	<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param username A valid username
	@param permission A valid permission
	@return ApiAddPermissionToUserRequest
	*/
	AddPermissionToUser(ctx context.Context, username string, permission string) ApiAddPermissionToUserRequest

	// AddPermissionToUserExecute executes the request
	//  @return UserPrincipal
	AddPermissionToUserExecute(r ApiAddPermissionToUserRequest) (*UserPrincipal, *http.Response, error)

	/*
	GetAllPermissions Returns a list of all permissions

	<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiGetAllPermissionsRequest
	*/
	GetAllPermissions(ctx context.Context) ApiGetAllPermissionsRequest

	// GetAllPermissionsExecute executes the request
	//  @return string
	GetAllPermissionsExecute(r ApiGetAllPermissionsRequest) (string, *http.Response, error)

	/*
	RemovePermissionFromTeam Removes the permission from the team.

	<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param uuid A valid team uuid
	@param permission A valid permission
	@return ApiRemovePermissionFromTeamRequest
	*/
	RemovePermissionFromTeam(ctx context.Context, uuid string, permission string) ApiRemovePermissionFromTeamRequest

	// RemovePermissionFromTeamExecute executes the request
	//  @return Team
	RemovePermissionFromTeamExecute(r ApiRemovePermissionFromTeamRequest) (*Team, *http.Response, error)

	/*
	RemovePermissionFromUser Removes the permission from the user.

	<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param username A valid username
	@param permission A valid permission
	@return ApiRemovePermissionFromUserRequest
	*/
	RemovePermissionFromUser(ctx context.Context, username string, permission string) ApiRemovePermissionFromUserRequest

	// RemovePermissionFromUserExecute executes the request
	//  @return UserPrincipal
	RemovePermissionFromUserExecute(r ApiRemovePermissionFromUserRequest) (*UserPrincipal, *http.Response, error)
}

// PermissionAPIService PermissionAPI service
type PermissionAPIService service

type ApiAddPermissionToTeamRequest struct {
	ctx context.Context
	ApiService PermissionAPI
	uuid string
	permission string
}

func (r ApiAddPermissionToTeamRequest) Execute() (*Team, *http.Response, error) {
	return r.ApiService.AddPermissionToTeamExecute(r)
}

/*
AddPermissionToTeam Adds the permission to the specified team.

<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param uuid A valid team uuid
 @param permission A valid permission
 @return ApiAddPermissionToTeamRequest
*/
func (a *PermissionAPIService) AddPermissionToTeam(ctx context.Context, uuid string, permission string) ApiAddPermissionToTeamRequest {
	return ApiAddPermissionToTeamRequest{
		ApiService: a,
		ctx: ctx,
		uuid: uuid,
		permission: permission,
	}
}

// Execute executes the request
//  @return Team
func (a *PermissionAPIService) AddPermissionToTeamExecute(r ApiAddPermissionToTeamRequest) (*Team, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Team
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PermissionAPIService.AddPermissionToTeam")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/permission/{permission}/team/{uuid}"
	localVarPath = strings.Replace(localVarPath, "{"+"uuid"+"}", url.PathEscape(parameterValueToString(r.uuid, "uuid")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"permission"+"}", url.PathEscape(parameterValueToString(r.permission, "permission")), -1)

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

type ApiAddPermissionToUserRequest struct {
	ctx context.Context
	ApiService PermissionAPI
	username string
	permission string
}

func (r ApiAddPermissionToUserRequest) Execute() (*UserPrincipal, *http.Response, error) {
	return r.ApiService.AddPermissionToUserExecute(r)
}

/*
AddPermissionToUser Adds the permission to the specified username.

<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param username A valid username
 @param permission A valid permission
 @return ApiAddPermissionToUserRequest
*/
func (a *PermissionAPIService) AddPermissionToUser(ctx context.Context, username string, permission string) ApiAddPermissionToUserRequest {
	return ApiAddPermissionToUserRequest{
		ApiService: a,
		ctx: ctx,
		username: username,
		permission: permission,
	}
}

// Execute executes the request
//  @return UserPrincipal
func (a *PermissionAPIService) AddPermissionToUserExecute(r ApiAddPermissionToUserRequest) (*UserPrincipal, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *UserPrincipal
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PermissionAPIService.AddPermissionToUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/permission/{permission}/user/{username}"
	localVarPath = strings.Replace(localVarPath, "{"+"username"+"}", url.PathEscape(parameterValueToString(r.username, "username")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"permission"+"}", url.PathEscape(parameterValueToString(r.permission, "permission")), -1)

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

type ApiGetAllPermissionsRequest struct {
	ctx context.Context
	ApiService PermissionAPI
}

func (r ApiGetAllPermissionsRequest) Execute() (string, *http.Response, error) {
	return r.ApiService.GetAllPermissionsExecute(r)
}

/*
GetAllPermissions Returns a list of all permissions

<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiGetAllPermissionsRequest
*/
func (a *PermissionAPIService) GetAllPermissions(ctx context.Context) ApiGetAllPermissionsRequest {
	return ApiGetAllPermissionsRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return string
func (a *PermissionAPIService) GetAllPermissionsExecute(r ApiGetAllPermissionsRequest) (string, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  string
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PermissionAPIService.GetAllPermissions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/permission"

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

type ApiRemovePermissionFromTeamRequest struct {
	ctx context.Context
	ApiService PermissionAPI
	uuid string
	permission string
}

func (r ApiRemovePermissionFromTeamRequest) Execute() (*Team, *http.Response, error) {
	return r.ApiService.RemovePermissionFromTeamExecute(r)
}

/*
RemovePermissionFromTeam Removes the permission from the team.

<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param uuid A valid team uuid
 @param permission A valid permission
 @return ApiRemovePermissionFromTeamRequest
*/
func (a *PermissionAPIService) RemovePermissionFromTeam(ctx context.Context, uuid string, permission string) ApiRemovePermissionFromTeamRequest {
	return ApiRemovePermissionFromTeamRequest{
		ApiService: a,
		ctx: ctx,
		uuid: uuid,
		permission: permission,
	}
}

// Execute executes the request
//  @return Team
func (a *PermissionAPIService) RemovePermissionFromTeamExecute(r ApiRemovePermissionFromTeamRequest) (*Team, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Team
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PermissionAPIService.RemovePermissionFromTeam")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/permission/{permission}/team/{uuid}"
	localVarPath = strings.Replace(localVarPath, "{"+"uuid"+"}", url.PathEscape(parameterValueToString(r.uuid, "uuid")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"permission"+"}", url.PathEscape(parameterValueToString(r.permission, "permission")), -1)

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

type ApiRemovePermissionFromUserRequest struct {
	ctx context.Context
	ApiService PermissionAPI
	username string
	permission string
}

func (r ApiRemovePermissionFromUserRequest) Execute() (*UserPrincipal, *http.Response, error) {
	return r.ApiService.RemovePermissionFromUserExecute(r)
}

/*
RemovePermissionFromUser Removes the permission from the user.

<p>Requires permission <strong>ACCESS_MANAGEMENT</strong></p>

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param username A valid username
 @param permission A valid permission
 @return ApiRemovePermissionFromUserRequest
*/
func (a *PermissionAPIService) RemovePermissionFromUser(ctx context.Context, username string, permission string) ApiRemovePermissionFromUserRequest {
	return ApiRemovePermissionFromUserRequest{
		ApiService: a,
		ctx: ctx,
		username: username,
		permission: permission,
	}
}

// Execute executes the request
//  @return UserPrincipal
func (a *PermissionAPIService) RemovePermissionFromUserExecute(r ApiRemovePermissionFromUserRequest) (*UserPrincipal, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *UserPrincipal
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PermissionAPIService.RemovePermissionFromUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/permission/{permission}/user/{username}"
	localVarPath = strings.Replace(localVarPath, "{"+"username"+"}", url.PathEscape(parameterValueToString(r.username, "username")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"permission"+"}", url.PathEscape(parameterValueToString(r.permission, "permission")), -1)

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
