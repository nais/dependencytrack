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

// checks if the AnalysisComment type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AnalysisComment{}

// AnalysisComment struct for AnalysisComment
type AnalysisComment struct {
	// UNIX epoch timestamp in milliseconds
	Timestamp int64 `json:"timestamp"`
	Comment string `json:"comment"`
	Commenter *string `json:"commenter,omitempty"`
}

type _AnalysisComment AnalysisComment

// NewAnalysisComment instantiates a new AnalysisComment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAnalysisComment(timestamp int64, comment string) *AnalysisComment {
	this := AnalysisComment{}
	this.Timestamp = timestamp
	this.Comment = comment
	return &this
}

// NewAnalysisCommentWithDefaults instantiates a new AnalysisComment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAnalysisCommentWithDefaults() *AnalysisComment {
	this := AnalysisComment{}
	return &this
}

// GetTimestamp returns the Timestamp field value
func (o *AnalysisComment) GetTimestamp() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value
// and a boolean to check if the value has been set.
func (o *AnalysisComment) GetTimestampOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Timestamp, true
}

// SetTimestamp sets field value
func (o *AnalysisComment) SetTimestamp(v int64) {
	o.Timestamp = v
}

// GetComment returns the Comment field value
func (o *AnalysisComment) GetComment() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Comment
}

// GetCommentOk returns a tuple with the Comment field value
// and a boolean to check if the value has been set.
func (o *AnalysisComment) GetCommentOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Comment, true
}

// SetComment sets field value
func (o *AnalysisComment) SetComment(v string) {
	o.Comment = v
}

// GetCommenter returns the Commenter field value if set, zero value otherwise.
func (o *AnalysisComment) GetCommenter() string {
	if o == nil || IsNil(o.Commenter) {
		var ret string
		return ret
	}
	return *o.Commenter
}

// GetCommenterOk returns a tuple with the Commenter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AnalysisComment) GetCommenterOk() (*string, bool) {
	if o == nil || IsNil(o.Commenter) {
		return nil, false
	}
	return o.Commenter, true
}

// HasCommenter returns a boolean if a field has been set.
func (o *AnalysisComment) HasCommenter() bool {
	if o != nil && !IsNil(o.Commenter) {
		return true
	}

	return false
}

// SetCommenter gets a reference to the given string and assigns it to the Commenter field.
func (o *AnalysisComment) SetCommenter(v string) {
	o.Commenter = &v
}

func (o AnalysisComment) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AnalysisComment) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["timestamp"] = o.Timestamp
	toSerialize["comment"] = o.Comment
	if !IsNil(o.Commenter) {
		toSerialize["commenter"] = o.Commenter
	}
	return toSerialize, nil
}

func (o *AnalysisComment) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"timestamp",
		"comment",
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

	varAnalysisComment := _AnalysisComment{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAnalysisComment)

	if err != nil {
		return err
	}

	*o = AnalysisComment(varAnalysisComment)

	return err
}

type NullableAnalysisComment struct {
	value *AnalysisComment
	isSet bool
}

func (v NullableAnalysisComment) Get() *AnalysisComment {
	return v.value
}

func (v *NullableAnalysisComment) Set(val *AnalysisComment) {
	v.value = val
	v.isSet = true
}

func (v NullableAnalysisComment) IsSet() bool {
	return v.isSet
}

func (v *NullableAnalysisComment) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAnalysisComment(val *AnalysisComment) *NullableAnalysisComment {
	return &NullableAnalysisComment{value: val, isSet: true}
}

func (v NullableAnalysisComment) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAnalysisComment) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


