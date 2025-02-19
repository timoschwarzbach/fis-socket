// https://github.com/yurisasuke/golang-decode-string-json-to-int/blob/main/main.go
package sequences

import (
	"encoding/json"
	"strconv"
)

// StringInt create a type alias for type int
type StringInt int

// UnmarshalJSON create a custom unmarshal for the StringInt
/// this helps us check the type of our value before unmarshalling it

func (st *StringInt) UnmarshalJSON(b []byte) error {
	//convert the bytes into an interface
	//this will help us check the type of our value
	//if it is a string that can be converted into a int we convert it
	///otherwise we return an error
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		///here convert the string into
		///an integer
		i, err := strconv.Atoi(v)
		if err != nil {
			// the string might not be of integer type
			// set default
			*st = StringInt(int(0))
			return nil
		}

		*st = StringInt(i)

	}
	return nil
}
