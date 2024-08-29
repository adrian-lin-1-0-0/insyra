package insyra

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/HazelnutParadise/Go-Utils/conv"
	"github.com/HazelnutParadise/Go-Utils/sliceutil"
)

// DataList is a generic dynamic data list
type DataList struct {
	data []interface{}
}

// IDataList defines the behavior expected from a DataList
type IDataList interface {
	Append(value interface{})
	Get(index int) interface{}
	Pop() interface{}
	Len() int
	Sort(acending ...bool)
	Reverse()
	Max() interface{}
	Min() interface{}
	Mean() float64
	GMean() float64
	Median() float64
	Mode() interface{}
	Stdev() float64
	Variance() float64
	Range() float64
	Quartiles() (float64, float64)
	IQR() float64
	Skewness() float64
	Kurtosis() float64
	ToF64Slice() []float64
}

// NewDataList creates a new DataList, supporting both slice and variadic inputs,
// and flattens the input before storing it.
func NewDataList(values ...interface{}) *DataList {
	// Flatten the input values using sliceutil.Flatten with the specified generic type
	flatData, err := sliceutil.Flatten[interface{}](values)
	if err != nil {
		flatData = values
	}
	return &DataList{
		data: flatData,
	}
}

// Append adds a new value to the DataList.
// The value can be of any type.
// The value is appended to the end of the DataList.
func (dl *DataList) Append(value interface{}) {
	dl.data = append(dl.data, value)
}

// Get retrieves the value at the specified index in the DataList.
// Supports negative indexing.
// Returns nil if the index is out of bounds.
// Returns the value at the specified index.
func (dl *DataList) Get(index int) interface{} {
	if index < -len(dl.data) || index >= len(dl.data) {
		return nil
	}
	return dl.data[index]
}

// Pop removes and returns the last element from the DataList.
// Returns the last element.
// Returns nil if the DataList is empty.
func (dl *DataList) Pop() interface{} {
	n, err := sliceutil.Drt_PopFrom(&dl.data)
	if err != nil {
		return nil
	}
	return n
}

func (dl *DataList) Len() int {
	return len(dl.data)
}

// Sort sorts the DataList using a mixed sorting logic.
// It handles string, numeric (including all integer and float types), and time data types.
// If sorting fails, it restores the original order.
func (dl *DataList) Sort(ascending ...bool) {
	if len(dl.data) == 0 {
		return
	}

	// Save the original order
	originalData := make([]interface{}, len(dl.data))
	copy(originalData, dl.data)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Sorting failed, restoring original order:", r)
			dl.data = originalData
		}
	}()

	ascendingOrder := true
	if len(ascending) > 0 {
		ascendingOrder = ascending[0]
	}

	// Mixed sorting
	sort.SliceStable(dl.data, func(i, j int) bool {
		v1 := dl.data[i]
		v2 := dl.data[j]

		switch v1 := v1.(type) {
		case string:
			if v2, ok := v2.(string); ok {
				return (v1 < v2) == ascendingOrder
			}
		case int, int8, int16, int32, int64:
			v1Float := ToFloat64(v1)
			if v2Float, ok := ToFloat64Safe(v2); ok {
				return (v1Float < v2Float) == ascendingOrder
			}
		case uint, uint8, uint16, uint32, uint64:
			v1Float := ToFloat64(v1)
			if v2Float, ok := ToFloat64Safe(v2); ok {
				return (v1Float < v2Float) == ascendingOrder
			}
		case float32, float64:
			v1Float := ToFloat64(v1)
			if v2Float, ok := ToFloat64Safe(v2); ok {
				return (v1Float < v2Float) == ascendingOrder
			}
		case time.Time:
			if v2, ok := v2.(time.Time); ok {
				return v1.Before(v2) == ascendingOrder
			}
		}

		// Fallback: compare as strings
		return fmt.Sprint(v1) < fmt.Sprint(v2)
	})
}

// Reverse reverses the order of the elements in the DataList.
func (dl *DataList) Reverse() {
	sliceutil.Reverse(dl.data)
}

// Max returns the maximum value in the DataList.
// Returns the maximum value.
// Returns nil if the DataList is empty.
// Max returns the maximum value in the DataList, or nil if the data types cannot be compared.
func (dl *DataList) Max() interface{} {
	if len(dl.data) == 0 {
		return nil
	}

	var max interface{}

	for _, v := range dl.data {
		if max == nil {
			max = v
			continue
		}

		switch maxVal := max.(type) {
		case int:
			if val, ok := v.(int); ok && val > maxVal {
				max = val
			} else if !ok {
				return nil
			}
		case float64:
			if val, ok := v.(float64); ok && val > maxVal {
				max = val
			} else if intVal, ok := v.(int); ok && float64(intVal) > maxVal {
				max = float64(intVal)
			} else if !ok {
				return nil
			}
		case string:
			if val, ok := v.(string); ok && conv.ParseF64(val) > conv.ParseF64(maxVal) {
				max = val
			} else if !ok {
				return nil
			}
		default:
			return nil
		}
	}

	return max
}

// Min returns the minimum value in the DataList.
// Returns the minimum value.
// Returns nil if the DataList is empty.
// Min returns the minimum value in the DataList, or nil if the data types cannot be compared.
func (dl *DataList) Min() interface{} {
	if len(dl.data) == 0 {
		return nil
	}

	var min interface{}

	for _, v := range dl.data {
		if min == nil {
			min = v
			continue
		}

		switch minVal := min.(type) {
		case int:
			if val, ok := v.(int); ok && val < minVal {
				min = val
			} else if !ok {
				return nil
			}
		case float64:
			if val, ok := v.(float64); ok && val < minVal {
				min = val
			} else if intVal, ok := v.(int); ok && float64(intVal) < minVal {
				min = float64(intVal)
			} else if !ok {
				return nil
			}
		case string:
			if val, ok := v.(string); ok && conv.ParseF64(val) < conv.ParseF64(minVal) {
				min = val
			} else if !ok {
				return nil
			}
		default:
			return nil
		}
	}

	return min
}

// Mean calculates the arithmetic mean of the DataList.
// Returns the arithmetic mean.
// Returns 0 if the DataList is empty.
// Mean returns the arithmetic mean of the DataList.
func (dl *DataList) Mean() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	var sum float64
	for _, v := range dl.data {
		if val, ok := ToFloat64Safe(v); ok {
			sum += val
		}
	}

	return sum / float64(len(dl.data))
}

// GMean calculates the geometric mean of the DataList.
// Returns the geometric mean.
// Returns 0 if the DataList is empty.
// GMean returns the geometric mean of the DataList.
func (dl *DataList) GMean() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	product := 1.0
	for _, v := range dl.data {
		if val, ok := ToFloat64Safe(v); ok {
			product *= val
		}
	}

	return math.Pow(product, 1.0/float64(len(dl.data)))
}

// Median calculates the median of the DataList.
// Returns the median.
// Returns 0 if the DataList is empty.
// Median returns the median of the DataList.
func (dl *DataList) Median() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	sortedData := make([]interface{}, len(dl.data))
	copy(sortedData, dl.data)
	dl.Sort()

	mid := len(sortedData) / 2
	if len(sortedData)%2 == 0 {
		return (ToFloat64(sortedData[mid-1]) + ToFloat64(sortedData[mid])) / 2
	}

	return ToFloat64(sortedData[mid])
}

// Mode calculates the mode of the DataList.
// Returns the mode.
// Returns nil if the DataList is empty.
// Mode returns the mode of the DataList.
func (dl *DataList) Mode() interface{} {
	if len(dl.data) == 0 {
		return nil
	}

	freqMap := make(map[interface{}]int)
	for _, v := range dl.data {
		freqMap[v]++
	}

	var mode interface{}
	maxFreq := 0
	for k, v := range freqMap {
		if v > maxFreq {
			mode = k
			maxFreq = v
		}
	}

	return mode
}

// Stdev calculates the standard deviation of the DataList.
// Returns the standard deviation.
// Returns 0 if the DataList is empty.
// Stdev returns the standard deviation of the DataList.
func (dl *DataList) Stdev() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	mean := dl.Mean()
	var sum float64
	for _, v := range dl.data {
		if val, ok := ToFloat64Safe(v); ok {
			sum += math.Pow(val-mean, 2)
		}
	}

	return math.Sqrt(sum / float64(len(dl.data)))
}

// Variance calculates the variance of the DataList.
// Returns the variance.
// Returns 0 if the DataList is empty.
// Variance returns the variance of the DataList.
func (dl *DataList) Variance() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	mean := dl.Mean()
	var sum float64
	for _, v := range dl.data {
		if val, ok := ToFloat64Safe(v); ok {
			sum += math.Pow(val-mean, 2)
		}
	}

	return sum / float64(len(dl.data))
}

// Range calculates the range of the DataList.
// Returns the range.
// Returns 0 if the DataList is empty.
// Range returns the range of the DataList.
func (dl *DataList) Range() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	max := ToFloat64(dl.Max())
	min := ToFloat64(dl.Min())

	return max - min
}

// Quartile calculates the quartile based on the input value (0 to 4).
// 0 corresponds to the minimum, 2 to the median (Q2), and 4 to the maximum.
func (dl *DataList) Quartile(q int) float64 {
	if len(dl.data) == 0 {
		return 0
	}
	if q < 0 || q > 4 {
		return 0
	}

	// Convert the DataList to a slice of float64 for numeric operations
	numericData := make([]float64, len(dl.data))
	for i, v := range dl.data {
		num, ok := ToFloat64Safe(v)
		if !ok {
			return 0
		}
		numericData[i] = num
	}

	// Sort the data
	sort.Float64s(numericData)

	n := len(numericData)
	switch q {
	case 0:
		return numericData[0] // Minimum value
	case 1:
		return NewDataList(numericData[:n/2]).Median() // First quartile (Q1)
	case 2:
		return NewDataList(numericData).Median() // Median (Q2)
	case 3:
		if n%2 == 0 {
			return NewDataList(numericData[n/2:]).Median() // Third quartile (Q3)
		}
		return NewDataList(numericData[n/2+1:]).Median()
	case 4:
		return numericData[n-1] // Maximum value
	}

	return 0
}

// IQR calculates the interquartile range of the DataList.
// Returns the interquartile range.
// Returns 0 if the DataList is empty.
// IQR returns the interquartile range of the DataList.
func (dl *DataList) IQR() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	sortedData := make([]interface{}, len(dl.data))
	copy(sortedData, dl.data)
	dl.Sort()

	return dl.Quartile(3) - dl.Quartile(1)
}

// Skewness calculates the skewness of the DataList.
// Returns the skewness.
// Returns 0 if the DataList is empty.
// Skewness returns the skewness of the DataList.
func (dl *DataList) Skewness() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	mean := dl.Mean()
	stdev := dl.Stdev()

	var sum float64
	for _, v := range dl.data {
		if val, ok := ToFloat64Safe(v); ok {
			sum += math.Pow(val-mean, 3)
		}
	}

	return sum / (float64(len(dl.data)) * math.Pow(stdev, 3))
}

// Kurtosis calculates the kurtosis of the DataList.
// Returns the kurtosis.
// Returns 0 if the DataList is empty.
// Kurtosis returns the kurtosis of the DataList.
func (dl *DataList) Kurtosis() float64 {
	if len(dl.data) == 0 {
		return 0
	}

	mean := dl.Mean()
	stdev := dl.Stdev()

	var sum float64
	for _, v := range dl.data {
		if val, ok := ToFloat64Safe(v); ok {
			sum += math.Pow(val-mean, 4)
		}
	}

	return sum / (float64(len(dl.data)) * math.Pow(stdev, 4))
}

// ToF64Slice converts the DataList to a float64 slice.
// Returns the float64 slice.
// Returns nil if the DataList is empty.
// ToF64Slice converts the DataList to a float64 slice.
func (dl *DataList) ToF64Slice() []float64 {
	if len(dl.data) == 0 {
		return nil
	}

	floatData := make([]float64, len(dl.data))
	for i, v := range dl.data {
		floatData[i] = ToFloat64(v)
	}

	return floatData
}
