package structure

import "fmt"

type ReaderFilter interface {
	Filter(reader *Reader) (bool, error)
}

type FieldEqualityFilter struct {
	fieldName string
	equalTo   any
}

// Filter implements ReaderFilter.
func (f *FieldEqualityFilter) Filter(reader *Reader) (bool, error) {
	readers, err := reader.Reader(f.fieldName)
	if err != nil {
		return false, nil
	}
	if len(readers) > 1 {
		return false, fmt.Errorf("filter target is expected to be a single value")
	}
	if len(readers) == 0 {
		return false, nil
	}
	value, err := readers[0].readScalar()
	if err != nil {
		return false, err
	}
	return value == f.equalTo, nil
}

func EqualFilter(fieldName string, equalTo any) ReaderFilter {
	return &FieldEqualityFilter{
		fieldName: fieldName,
		equalTo:   equalTo,
	}
}
