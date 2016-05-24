// Copyright 2012-2016 Aerospike, Inc.
//
// Portions may be licensed to Aerospike, Inc. under one or more contributor
// license agreements WHICH ARE COMPATIBLE WITH THE APACHE LICENSE, VERSION 2.0.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package aerospike

const (
	_CDT_MAP_SET_TYPE                 = 64
	_CDT_MAP_ADD                      = 65
	_CDT_MAP_ADD_ITEMS                = 66
	_CDT_MAP_PUT                      = 67
	_CDT_MAP_PUT_ITEMS                = 68
	_CDT_MAP_REPLACE                  = 69
	_CDT_MAP_REPLACE_ITEMS            = 70
	_CDT_MAP_INCREMENT                = 73
	_CDT_MAP_DECREMENT                = 74
	_CDT_MAP_CLEAR                    = 75
	_CDT_MAP_REMOVE_BY_KEY            = 76
	_CDT_MAP_REMOVE_BY_INDEX          = 77
	_CDT_MAP_REMOVE_BY_RANK           = 79
	_CDT_MAP_REMOVE_KEY_LIST          = 81
	_CDT_MAP_REMOVE_BY_VALUE          = 82
	_CDT_MAP_REMOVE_VALUE_LIST        = 83
	_CDT_MAP_REMOVE_BY_KEY_INTERVAL   = 84
	_CDT_MAP_REMOVE_BY_INDEX_RANGE    = 85
	_CDT_MAP_REMOVE_BY_VALUE_INTERVAL = 86
	_CDT_MAP_REMOVE_BY_RANK_RANGE     = 87
	_CDT_MAP_SIZE                     = 96
	_CDT_MAP_GET_BY_KEY               = 97
	_CDT_MAP_GET_BY_INDEX             = 98
	_CDT_MAP_GET_BY_RANK              = 100
	_CDT_MAP_GET_BY_VALUE             = 102
	_CDT_MAP_GET_BY_KEY_INTERVAL      = 103
	_CDT_MAP_GET_BY_INDEX_RANGE       = 104
	_CDT_MAP_GET_BY_VALUE_INTERVAL    = 105
	_CDT_MAP_GET_BY_RANK_RANGE        = 106
)

type mapOrderType int

// Map storage order.
var MapOrder = struct {
	// Map is not ordered. This is the default.
	UNORDERED mapOrderType // 0

	// Order map by key.
	KEY_ORDERED mapOrderType // 1

	// Order map by key, then value.
	KEY_VALUE_ORDERED mapOrderType // 3
}{0, 1, 3}

type mapReturnType int

// Map return type. Type of data to return when selecting or removing items from the map.
var MapReturnType = struct {
	// Do not return a result.
	NONE mapReturnType

	// Return key index order.
	//
	// 0 = first key
	// N = Nth key
	// -1 = last key
	INDEX mapReturnType

	// Return reverse key order.
	//
	// 0 = last key
	// -1 = first key
	REVERSE_INDEX mapReturnType

	// Return value order.
	//
	// 0 = smallest value
	// N = Nth smallest value
	// -1 = largest value
	RANK mapReturnType

	// Return reserve value order.
	//
	// 0 = largest value
	// N = Nth largest value
	// -1 = smallest value
	REVERSE_RANK mapReturnType

	// Return count of items selected.
	COUNT mapReturnType

	// Return key for single key read and key list for range read.
	KEY mapReturnType

	// Return value for single key read and value list for range read.
	VALUE mapReturnType

	// Return key/value items. The possible return types are:
	//
	// map[interface{}]interface{} : Returned for unordered maps
	// []MapPair : Returned for range results where range order needs to be preserved.
	KEY_VALUE mapReturnType
}{
	0, 1, 2, 3, 4, 5, 6, 7, 8,
}

// Unique key map write type.
type mapWriteMode struct {
	itemCommand  int
	itemsCommand int
}

var MapWriteMode = struct {
	// If the key already exists, the item will be overwritten.
	// If the key does not exist, a new item will be created.
	Update *mapWriteMode

	// If the key already exists, the item will be overwritten.
	// If the key does not exist, the write will fail.
	UpdateOnly *mapWriteMode

	// If the key already exists, the write will fail.
	// If the key does not exist, a new item will be created.
	CreateOnly *mapWriteMode
}{
	&mapWriteMode{_CDT_MAP_PUT, _CDT_MAP_PUT_ITEMS},
	&mapWriteMode{_CDT_MAP_REPLACE, _CDT_MAP_REPLACE_ITEMS},
	&mapWriteMode{_CDT_MAP_ADD, _CDT_MAP_ADD_ITEMS},
}

// MapPolicy directives when creating a map and writing map items.
type MapPolicy struct {
	attributes   mapOrderType
	itemCommand  int
	itemsCommand int
}

// Create unique key map with specified order when map does not exist.
// Use specified write mode when writing map items.
func NewMapPolicy(order mapOrderType, writeMode *mapWriteMode) *MapPolicy {
	return &MapPolicy{
		attributes:   order,
		itemCommand:  writeMode.itemCommand,
		itemsCommand: writeMode.itemsCommand,
	}
}

func newMapSetPolicy(binName string, attributes mapOrderType) *Operation {
	packer := newPacker()
	packer.PackShortRaw(_CDT_MAP_SET_TYPE)
	packer.PackArrayBegin(1)
	packer.PackAInt(int(attributes))
	return &Operation{
		OpType:   MAP_MODIFY,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes())}
}

func newMapCreatePut(command int, attributes mapOrderType, binName string, value1 interface{}, value2 interface{}) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(command))

	if command == _CDT_MAP_REPLACE {
		// Replace doesn't allow map attributes because it does not create on non-existing key.
		packer.PackArrayBegin(2)
		NewValue(value1).pack(packer)
		NewValue(value2).pack(packer)
	} else {
		packer.PackArrayBegin(3)
		NewValue(value1).pack(packer)
		NewValue(value2).pack(packer)
		packer.PackAInt(int(attributes))
	}
	return &Operation{
		OpType:   MAP_MODIFY,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

func newMapCreateOperationValues2(command int, attributes mapOrderType, binName string, value1 interface{}, value2 interface{}) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(command))
	packer.PackArrayBegin(3)
	NewValue(value1).pack(packer)
	NewValue(value2).pack(packer)
	packer.PackAInt(int(attributes))
	return &Operation{
		OpType:   MAP_MODIFY,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

func newMapCreateOperationValues0(command int, typ OperationType, binName string) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(command))
	return &Operation{
		OpType:   typ,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

func newMapCreateOperationValuesN(command int, typ OperationType, binName string, values []interface{}, returnType mapReturnType) *Operation {
	list := ToValueSlice(values)

	packer := newPacker()
	packer.PackShortRaw(int16(command))
	packer.PackArrayBegin(2)
	packer.PackAInt(int(returnType))
	packer.packValueArray(list)
	return &Operation{
		OpType:   typ,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

func newMapCreateOperationValue1(command int, typ OperationType, binName string, value interface{}, returnType mapReturnType) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(command))
	packer.PackArrayBegin(2)
	packer.PackAInt(int(returnType))
	NewValue(value).pack(packer)
	return &Operation{
		OpType:   typ,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

func newMapCreateOperationIndex(command int, typ OperationType, binName string, index int, returnType mapReturnType) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(command))
	packer.PackArrayBegin(2)
	packer.PackAInt(int(returnType))
	packer.PackAInt(index)
	return &Operation{
		OpType:   typ,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

func newMapCreateOperationIndexCount(command int, typ OperationType, binName string, index int, count int, returnType mapReturnType) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(command))
	packer.PackArrayBegin(3)
	packer.PackAInt(int(returnType))
	packer.PackAInt(index)
	packer.PackAInt(count)
	return &Operation{
		OpType:   typ,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

func newMapCreateRangeOperation(command int, typ OperationType, binName string, begin interface{}, end interface{}, returnType mapReturnType) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(command))

	if begin == nil {
		begin = NewNullValue()
	}

	if end == nil {
		packer.PackArrayBegin(2)
		packer.PackAInt(int(returnType))
		NewValue(begin).pack(packer)
	} else {
		packer.PackArrayBegin(3)
		packer.PackAInt(int(returnType))
		NewValue(begin).pack(packer)
		NewValue(end).pack(packer)
	}
	return &Operation{
		OpType:   typ,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

/////////////////////////

// Unique key map bin operations. Create map operations used by the client operate command.
// The default unique key map is unordered.
//
// All maps maintain an index and a rank.  The index is the item offset from the start of the map,
// for both unordered and ordered maps.  The rank is the sorted index of the value component.
// Map supports negative indexing for index and rank.
//
// Index examples:
//
// Index 0: First item in map.
// Index 4: Fifth item in map.
// Index -1: Last item in map.
// Index -3: Third to last item in map.
// Index 1 Count 2: Second and third items in map.
// Index -3 Count 3: Last three items in map.
// Index -5 Count 4: Range between fifth to last item to second to last item inclusive.
//
// Rank examples:
//
// Rank 0: Item with lowest value rank in map.
// Rank 4: Fifth lowest ranked item in map.
// Rank -1: Item with highest ranked value in map.
// Rank -3: Item with third highest ranked value in map.
// Rank 1 Count 2: Second and third lowest ranked items in map.
// Rank -3 Count 3: Top three ranked items in map.

// Create set map policy operation.
// Server sets map policy attributes.  Server returns null.
//
// The required map policy attributes can be changed after the map is created.
func MapSetPolicyOp(policy *MapPolicy, binName string) *Operation {
	return newMapSetPolicy(binName, policy.attributes)
}

// Create map put operation.
// Server writes key/value item to map bin and returns map size.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
// See policy {@link com.aerospike.client.cdt.MapPolicy} and write mode
// {@link com.aerospike.client.cdt.mapWriteMode}.
func MapPutOp(policy *MapPolicy, binName string, key interface{}, value interface{}) *Operation {
	return newMapCreatePut(policy.itemCommand, policy.attributes, binName, key, value)
}

// Create map put items operation
// Server writes each map item to map bin and returns map size.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
// See policy {@link com.aerospike.client.cdt.MapPolicy} and write mode
// {@link com.aerospike.client.cdt.mapWriteMode}.
func MapPutItemsOp(policy *MapPolicy, binName string, amap map[interface{}]interface{}) *Operation {
	packer := newPacker()
	packer.PackShortRaw(int16(policy.itemsCommand))

	if policy.itemsCommand == int(_CDT_MAP_REPLACE_ITEMS) {
		// Replace doesn't allow map attributes because it does not create on non-existing key.
		packer.PackArrayBegin(1)
		packer.PackMap(amap)
	} else {
		packer.PackArrayBegin(2)
		packer.PackMap(amap)
		packer.PackAInt(int(policy.attributes))
	}
	return &Operation{
		OpType:   MAP_MODIFY,
		BinName:  binName,
		BinValue: NewValue(packer.buffer.Bytes()),
	}
}

// Create map increment operation.
// Server increments values by incr for all items identified by key and returns final result.
// Valid only for numbers.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
// See policy {@link com.aerospike.client.cdt.MapPolicy} and write mode
// {@link com.aerospike.client.cdt.mapWriteMode}.
func MapIncrementOp(policy *MapPolicy, binName string, key interface{}, incr interface{}) *Operation {
	return newMapCreateOperationValues2(_CDT_MAP_INCREMENT, policy.attributes, binName, key, incr)
}

// Create map decrement operation.
// Server decrements values by decr for all items identified by key and returns final result.
// Valid only for numbers.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
// See policy {@link com.aerospike.client.cdt.MapPolicy} and write mode
// {@link com.aerospike.client.cdt.mapWriteMode}.
func MapDecrementOp(policy *MapPolicy, binName string, key interface{}, decr interface{}) *Operation {
	return newMapCreateOperationValues2(_CDT_MAP_DECREMENT, policy.attributes, binName, key, decr)
}

// Create map clear operation.
// Server removes all items in map.  Server returns null.
func MapClearOp(binName string) *Operation {
	return newMapCreateOperationValues0(_CDT_MAP_CLEAR, MAP_MODIFY, binName)
}

// Create map remove operation.
// Server removes map item identified by key and returns removed data specified by returnType.
func MapRemoveByKeyOp(binName string, key interface{}, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_REMOVE_BY_KEY, MAP_MODIFY, binName, key, returnType)
}

// Create map remove operation.
// Server removes map items identified by keys and returns removed data specified by returnType.
func MapRemoveByKeyListOp(binName string, keys []interface{}, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_REMOVE_KEY_LIST, MAP_MODIFY, binName, keys, returnType)
}

// Create map remove operation.
// Server removes map items identified by key range (keyBegin inclusive, keyEnd exclusive).
// If keyBegin is null, the range is less than keyEnd.
// If keyEnd is null, the range is greater than equal to keyBegin.
//
// Server returns removed data specified by returnType.
func MapRemoveByKeyRangeOp(binName string, keyBegin interface{}, keyEnd interface{}, returnType mapReturnType) *Operation {
	return newMapCreateRangeOperation(_CDT_MAP_REMOVE_BY_KEY_INTERVAL, MAP_MODIFY, binName, keyBegin, keyEnd, returnType)
}

// Create map remove operation.
// Server removes map items identified by value and returns removed data specified by returnType.
func MapRemoveByValueOp(binName string, value interface{}, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_REMOVE_BY_VALUE, MAP_MODIFY, binName, value, returnType)
}

// Create map remove operation.
// Server removes map items identified by values and returns removed data specified by returnType.
func MapRemoveByValueListOp(binName string, values []interface{}, returnType mapReturnType) *Operation {
	return newMapCreateOperationValuesN(_CDT_MAP_REMOVE_VALUE_LIST, MAP_MODIFY, binName, values, returnType)
}

// Create map remove operation.
// Server removes map items identified by value range (valueBegin inclusive, valueEnd exclusive).
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
//
// Server returns removed data specified by returnType.
func MapRemoveByValueRangeOp(binName string, valueBegin interface{}, valueEnd interface{}, returnType mapReturnType) *Operation {
	return newMapCreateRangeOperation(_CDT_MAP_REMOVE_BY_VALUE_INTERVAL, MAP_MODIFY, binName, valueBegin, valueEnd, returnType)
}

// Create map remove operation.
// Server removes map item identified by index and returns removed data specified by returnType.
func MapRemoveByIndexOp(binName string, index int, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_REMOVE_BY_INDEX, MAP_MODIFY, binName, index, returnType)
}

// Create map remove operation.
// Server removes map items starting at specified index to the end of map and returns removed
// data specified by returnType.
func MapRemoveByIndexRangeOp(binName string, index int, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_REMOVE_BY_INDEX_RANGE, MAP_MODIFY, binName, index, returnType)
}

// Create map remove operation.
// Server removes "count" map items starting at specified index and returns removed data specified by returnType.
func MapRemoveByIndexRangeCountOp(binName string, index int, count int, returnType mapReturnType) *Operation {
	return newMapCreateOperationIndexCount(_CDT_MAP_REMOVE_BY_INDEX_RANGE, MAP_MODIFY, binName, index, count, returnType)
}

// Create map remove operation.
// Server removes map item identified by rank and returns removed data specified by returnType.
func MapRemoveByRankOp(binName string, rank int, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_REMOVE_BY_RANK, MAP_MODIFY, binName, rank, returnType)
}

// Create map remove operation.
// Server removes map items starting at specified rank to the last ranked item and returns removed
// data specified by returnType.
func MapRemoveByRankRangeOp(binName string, rank int, returnType mapReturnType) *Operation {
	return newMapCreateOperationIndex(_CDT_MAP_REMOVE_BY_RANK_RANGE, MAP_MODIFY, binName, rank, returnType)
}

// Create map remove operation.
// Server removes "count" map items starting at specified rank and returns removed data specified by returnType.
func MapRemoveByRankRangeCountOp(binName string, rank int, count int, returnType mapReturnType) *Operation {
	return newMapCreateOperationIndexCount(_CDT_MAP_REMOVE_BY_RANK_RANGE, MAP_MODIFY, binName, rank, count, returnType)
}

// Create map size operation.
// Server returns size of map.
func MapSizeOp(binName string) *Operation {
	return newMapCreateOperationValues0(_CDT_MAP_SIZE, MAP_READ, binName)
}

// Create map get by key operation.
// Server selects map item identified by key and returns selected data specified by returnType.
func MapGetByKeyOp(binName string, key interface{}, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_GET_BY_KEY, MAP_READ, binName, key, returnType)
}

// Create map get by key range operation.
// Server selects map items identified by key range (keyBegin inclusive, keyEnd exclusive).
// If keyBegin is null, the range is less than keyEnd.
// If keyEnd is null, the range is greater than equal to keyBegin.
//
// Server returns selected data specified by returnType.
func MapGetByKeyRangeOp(binName string, keyBegin interface{}, keyEnd interface{}, returnType mapReturnType) *Operation {
	return newMapCreateRangeOperation(_CDT_MAP_GET_BY_KEY_INTERVAL, MAP_READ, binName, keyBegin, keyEnd, returnType)
}

// Create map get by value operation.
// Server selects map items identified by value and returns selected data specified by returnType.
func MapGetByValueOp(binName string, value interface{}, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_GET_BY_VALUE, MAP_READ, binName, value, returnType)
}

// Create map get by value range operation.
// Server selects map items identified by value range (valueBegin inclusive, valueEnd exclusive)
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
//
// Server returns selected data specified by returnType.
func MapGetByValueRangeOp(binName string, valueBegin interface{}, valueEnd interface{}, returnType mapReturnType) *Operation {
	return newMapCreateRangeOperation(_CDT_MAP_GET_BY_VALUE_INTERVAL, MAP_READ, binName, valueBegin, valueEnd, returnType)
}

// Create map get by index operation.
// Server selects map item identified by index and returns selected data specified by returnType.
func MapGetByIndexOp(binName string, index int, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_GET_BY_INDEX, MAP_READ, binName, index, returnType)
}

// Create map get by index range operation.
// Server selects map items starting at specified index to the end of map and returns selected
// data specified by returnType.
func MapGetByIndexRangeOp(binName string, index int, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_GET_BY_INDEX_RANGE, MAP_READ, binName, index, returnType)
}

// Create map get by index range operation.
// Server selects "count" map items starting at specified index and returns selected data specified by returnType.
func MapGetByIndexRangeCountOp(binName string, index int, count int, returnType mapReturnType) *Operation {
	return newMapCreateOperationIndexCount(_CDT_MAP_GET_BY_INDEX_RANGE, MAP_READ, binName, index, count, returnType)
}

// Create map get by rank operation.
// Server selects map item identified by rank and returns selected data specified by returnType.
func MapGetByRankOp(binName string, rank int, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_GET_BY_RANK, MAP_READ, binName, rank, returnType)
}

// Create map get by rank range operation.
// Server selects map items starting at specified rank to the last ranked item and returns selected
// data specified by returnType.
func MapGetByRankRangeOp(binName string, rank int, returnType mapReturnType) *Operation {
	return newMapCreateOperationValue1(_CDT_MAP_GET_BY_RANK_RANGE, MAP_READ, binName, rank, returnType)
}

// Create map get by rank range operation.
// Server selects "count" map items starting at specified rank and returns selected data specified by returnType.
func MapGetByRankRangeCountOp(binName string, rank int, count int, returnType mapReturnType) *Operation {
	return newMapCreateOperationIndexCount(_CDT_MAP_GET_BY_RANK_RANGE, MAP_READ, binName, rank, count, returnType)
}
