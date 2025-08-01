/*
OWASP Dependency-Track

REST API of OWASP Dependency-Track

API version: 4.13.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the Repository type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Repository{}

// Repository struct for Repository
type Repository struct {
	Type string `json:"type"`
	Identifier *string `json:"identifier,omitempty"`
	Url *string `json:"url,omitempty"`
	ResolutionOrder int32 `json:"resolutionOrder"`
	Enabled bool `json:"enabled"`
	Internal bool `json:"internal"`
	AuthenticationRequired *bool `json:"authenticationRequired,omitempty"`
	Username *string `json:"username,omitempty"`
	Uuid string `json:"uuid"`
	Password *string `json:"password,omitempty"`
}

type _Repository Repository

// NewRepository instantiates a new Repository object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRepository(type_ string, resolutionOrder int32, enabled bool, internal bool, uuid string) *Repository {
	this := Repository{}
	this.Type = type_
	this.ResolutionOrder = resolutionOrder
	this.Enabled = enabled
	this.Internal = internal
	this.Uuid = uuid
	return &this
}

// NewRepositoryWithDefaults instantiates a new Repository object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRepositoryWithDefaults() *Repository {
	this := Repository{}
	return &this
}

// GetType returns the Type field value
func (o *Repository) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *Repository) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *Repository) SetType(v string) {
	o.Type = v
}

// GetIdentifier returns the Identifier field value if set, zero value otherwise.
func (o *Repository) GetIdentifier() string {
	if o == nil || IsNil(o.Identifier) {
		var ret string
		return ret
	}
	return *o.Identifier
}

// GetIdentifierOk returns a tuple with the Identifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Repository) GetIdentifierOk() (*string, bool) {
	if o == nil || IsNil(o.Identifier) {
		return nil, false
	}
	return o.Identifier, true
}

// HasIdentifier returns a boolean if a field has been set.
func (o *Repository) HasIdentifier() bool {
	if o != nil && !IsNil(o.Identifier) {
		return true
	}

	return false
}

// SetIdentifier gets a reference to the given string and assigns it to the Identifier field.
func (o *Repository) SetIdentifier(v string) {
	o.Identifier = &v
}

// GetUrl returns the Url field value if set, zero value otherwise.
func (o *Repository) GetUrl() string {
	if o == nil || IsNil(o.Url) {
		var ret string
		return ret
	}
	return *o.Url
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Repository) GetUrlOk() (*string, bool) {
	if o == nil || IsNil(o.Url) {
		return nil, false
	}
	return o.Url, true
}

// HasUrl returns a boolean if a field has been set.
func (o *Repository) HasUrl() bool {
	if o != nil && !IsNil(o.Url) {
		return true
	}

	return false
}

// SetUrl gets a reference to the given string and assigns it to the Url field.
func (o *Repository) SetUrl(v string) {
	o.Url = &v
}

// GetResolutionOrder returns the ResolutionOrder field value
func (o *Repository) GetResolutionOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ResolutionOrder
}

// GetResolutionOrderOk returns a tuple with the ResolutionOrder field value
// and a boolean to check if the value has been set.
func (o *Repository) GetResolutionOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ResolutionOrder, true
}

// SetResolutionOrder sets field value
func (o *Repository) SetResolutionOrder(v int32) {
	o.ResolutionOrder = v
}

// GetEnabled returns the Enabled field value
func (o *Repository) GetEnabled() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value
// and a boolean to check if the value has been set.
func (o *Repository) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Enabled, true
}

// SetEnabled sets field value
func (o *Repository) SetEnabled(v bool) {
	o.Enabled = v
}

// GetInternal returns the Internal field value
func (o *Repository) GetInternal() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Internal
}

// GetInternalOk returns a tuple with the Internal field value
// and a boolean to check if the value has been set.
func (o *Repository) GetInternalOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Internal, true
}

// SetInternal sets field value
func (o *Repository) SetInternal(v bool) {
	o.Internal = v
}

// GetAuthenticationRequired returns the AuthenticationRequired field value if set, zero value otherwise.
func (o *Repository) GetAuthenticationRequired() bool {
	if o == nil || IsNil(o.AuthenticationRequired) {
		var ret bool
		return ret
	}
	return *o.AuthenticationRequired
}

// GetAuthenticationRequiredOk returns a tuple with the AuthenticationRequired field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Repository) GetAuthenticationRequiredOk() (*bool, bool) {
	if o == nil || IsNil(o.AuthenticationRequired) {
		return nil, false
	}
	return o.AuthenticationRequired, true
}

// HasAuthenticationRequired returns a boolean if a field has been set.
func (o *Repository) HasAuthenticationRequired() bool {
	if o != nil && !IsNil(o.AuthenticationRequired) {
		return true
	}

	return false
}

// SetAuthenticationRequired gets a reference to the given bool and assigns it to the AuthenticationRequired field.
func (o *Repository) SetAuthenticationRequired(v bool) {
	o.AuthenticationRequired = &v
}

// GetUsername returns the Username field value if set, zero value otherwise.
func (o *Repository) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Repository) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}
	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *Repository) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *Repository) SetUsername(v string) {
	o.Username = &v
}

// GetUuid returns the Uuid field value
func (o *Repository) GetUuid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value
// and a boolean to check if the value has been set.
func (o *Repository) GetUuidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uuid, true
}

// SetUuid sets field value
func (o *Repository) SetUuid(v string) {
	o.Uuid = v
}

// GetPassword returns the Password field value if set, zero value otherwise.
func (o *Repository) GetPassword() string {
	if o == nil || IsNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Repository) GetPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.Password) {
		return nil, false
	}
	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *Repository) HasPassword() bool {
	if o != nil && !IsNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *Repository) SetPassword(v string) {
	o.Password = &v
}

func (o Repository) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Repository) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["type"] = o.Type
	if !IsNil(o.Identifier) {
		toSerialize["identifier"] = o.Identifier
	}
	if !IsNil(o.Url) {
		toSerialize["url"] = o.Url
	}
	toSerialize["resolutionOrder"] = o.ResolutionOrder
	toSerialize["enabled"] = o.Enabled
	toSerialize["internal"] = o.Internal
	if !IsNil(o.AuthenticationRequired) {
		toSerialize["authenticationRequired"] = o.AuthenticationRequired
	}
	if !IsNil(o.Username) {
		toSerialize["username"] = o.Username
	}
	toSerialize["uuid"] = o.Uuid
	if !IsNil(o.Password) {
		toSerialize["password"] = o.Password
	}
	return toSerialize, nil
}

func (o *Repository) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"type",
		"resolutionOrder",
		"enabled",
		"internal",
		"uuid",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varRepository := _Repository{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varRepository)

	if err != nil {
		return err
	}

	*o = Repository(varRepository)

	return err
}

type NullableRepository struct {
	value *Repository
	isSet bool
}

func (v NullableRepository) Get() *Repository {
	return v.value
}

func (v *NullableRepository) Set(val *Repository) {
	v.value = val
	v.isSet = true
}

func (v NullableRepository) IsSet() bool {
	return v.isSet
}

func (v *NullableRepository) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRepository(val *Repository) *NullableRepository {
	return &NullableRepository{value: val, isSet: true}
}

func (v NullableRepository) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRepository) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


