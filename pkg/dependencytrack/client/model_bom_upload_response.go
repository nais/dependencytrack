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

// checks if the BomUploadResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BomUploadResponse{}

// BomUploadResponse struct for BomUploadResponse
type BomUploadResponse struct {
	// Token used to check task progress
	Token string `json:"token"`
}

type _BomUploadResponse BomUploadResponse

// NewBomUploadResponse instantiates a new BomUploadResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBomUploadResponse(token string) *BomUploadResponse {
	this := BomUploadResponse{}
	this.Token = token
	return &this
}

// NewBomUploadResponseWithDefaults instantiates a new BomUploadResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBomUploadResponseWithDefaults() *BomUploadResponse {
	this := BomUploadResponse{}
	return &this
}

// GetToken returns the Token field value
func (o *BomUploadResponse) GetToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Token
}

// GetTokenOk returns a tuple with the Token field value
// and a boolean to check if the value has been set.
func (o *BomUploadResponse) GetTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Token, true
}

// SetToken sets field value
func (o *BomUploadResponse) SetToken(v string) {
	o.Token = v
}

func (o BomUploadResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BomUploadResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["token"] = o.Token
	return toSerialize, nil
}

func (o *BomUploadResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"token",
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

	varBomUploadResponse := _BomUploadResponse{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varBomUploadResponse)

	if err != nil {
		return err
	}

	*o = BomUploadResponse(varBomUploadResponse)

	return err
}

type NullableBomUploadResponse struct {
	value *BomUploadResponse
	isSet bool
}

func (v NullableBomUploadResponse) Get() *BomUploadResponse {
	return v.value
}

func (v *NullableBomUploadResponse) Set(val *BomUploadResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableBomUploadResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableBomUploadResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBomUploadResponse(val *BomUploadResponse) *NullableBomUploadResponse {
	return &NullableBomUploadResponse{value: val, isSet: true}
}

func (v NullableBomUploadResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBomUploadResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


