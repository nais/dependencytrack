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

// checks if the TaggedCollectionProjectListResponseItem type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TaggedCollectionProjectListResponseItem{}

// TaggedCollectionProjectListResponseItem struct for TaggedCollectionProjectListResponseItem
type TaggedCollectionProjectListResponseItem struct {
	// UUID of the collection project
	Uuid string `json:"uuid"`
	// Name of the collection project
	Name string `json:"name"`
	// Version of the collection project
	Version *string `json:"version,omitempty"`
}

type _TaggedCollectionProjectListResponseItem TaggedCollectionProjectListResponseItem

// NewTaggedCollectionProjectListResponseItem instantiates a new TaggedCollectionProjectListResponseItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaggedCollectionProjectListResponseItem(uuid string, name string) *TaggedCollectionProjectListResponseItem {
	this := TaggedCollectionProjectListResponseItem{}
	this.Uuid = uuid
	this.Name = name
	return &this
}

// NewTaggedCollectionProjectListResponseItemWithDefaults instantiates a new TaggedCollectionProjectListResponseItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaggedCollectionProjectListResponseItemWithDefaults() *TaggedCollectionProjectListResponseItem {
	this := TaggedCollectionProjectListResponseItem{}
	return &this
}

// GetUuid returns the Uuid field value
func (o *TaggedCollectionProjectListResponseItem) GetUuid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value
// and a boolean to check if the value has been set.
func (o *TaggedCollectionProjectListResponseItem) GetUuidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uuid, true
}

// SetUuid sets field value
func (o *TaggedCollectionProjectListResponseItem) SetUuid(v string) {
	o.Uuid = v
}

// GetName returns the Name field value
func (o *TaggedCollectionProjectListResponseItem) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *TaggedCollectionProjectListResponseItem) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *TaggedCollectionProjectListResponseItem) SetName(v string) {
	o.Name = v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *TaggedCollectionProjectListResponseItem) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaggedCollectionProjectListResponseItem) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *TaggedCollectionProjectListResponseItem) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *TaggedCollectionProjectListResponseItem) SetVersion(v string) {
	o.Version = &v
}

func (o TaggedCollectionProjectListResponseItem) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TaggedCollectionProjectListResponseItem) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["uuid"] = o.Uuid
	toSerialize["name"] = o.Name
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	return toSerialize, nil
}

func (o *TaggedCollectionProjectListResponseItem) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"uuid",
		"name",
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

	varTaggedCollectionProjectListResponseItem := _TaggedCollectionProjectListResponseItem{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varTaggedCollectionProjectListResponseItem)

	if err != nil {
		return err
	}

	*o = TaggedCollectionProjectListResponseItem(varTaggedCollectionProjectListResponseItem)

	return err
}

type NullableTaggedCollectionProjectListResponseItem struct {
	value *TaggedCollectionProjectListResponseItem
	isSet bool
}

func (v NullableTaggedCollectionProjectListResponseItem) Get() *TaggedCollectionProjectListResponseItem {
	return v.value
}

func (v *NullableTaggedCollectionProjectListResponseItem) Set(val *TaggedCollectionProjectListResponseItem) {
	v.value = val
	v.isSet = true
}

func (v NullableTaggedCollectionProjectListResponseItem) IsSet() bool {
	return v.isSet
}

func (v *NullableTaggedCollectionProjectListResponseItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaggedCollectionProjectListResponseItem(val *TaggedCollectionProjectListResponseItem) *NullableTaggedCollectionProjectListResponseItem {
	return &NullableTaggedCollectionProjectListResponseItem{value: val, isSet: true}
}

func (v NullableTaggedCollectionProjectListResponseItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaggedCollectionProjectListResponseItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


