package structuredata

import "fmt"

func EqualStructureData(a StructureData, b StructureData) (bool, error) {
	aType, err := a.Type()
	if err != nil {
		return false, err
	}
	bType, err := b.Type()
	if err != nil {
		return false, err
	}
	if aType != bType {
		return false, nil
	}
	aKeys, err := a.Keys()
	if err != nil {
		return false, err
	}
	bKeys, err := b.Keys()
	if err != nil {
		return false, err
	}
	if len(aKeys) != len(bKeys) {
		return false, nil
	}
	aKeysInStruct := map[string]struct{}{}
	for _, aKey := range aKeys {
		aKeysInStruct[aKey] = struct{}{}
	}
	for _, bKey := range bKeys {
		if _, exist := aKeysInStruct[bKey]; !exist {
			return false, nil
		}
	}
	if aType == StructuredTypeScalar {
		aValue, err := a.Value("")
		if err != nil {
			return false, err
		}
		bValue, err := b.Value("")
		if err != nil {
			return false, err
		}
		return aValue == bValue, nil
	} else {
		for _, aKey := range aKeys {
			aValue, err := a.Value(aKey)
			if err != nil {
				return false, err
			}
			if aValueStructured, convertible := aValue.(StructureData); !convertible {
				return false, fmt.Errorf("failed to convert result to StructureData")
			} else {
				bValue, err := b.Value(aKey)
				if err != nil {
					return false, err
				}
				if bValueStructured, convertible := bValue.(StructureData); !convertible {
					return false, fmt.Errorf("failed to convert result to StructureData")
				} else {
					eq, err := EqualStructureData(aValueStructured, bValueStructured)
					if err != nil {
						return false, err
					}
					if !eq {
						return false, nil
					}
				}
			}
		}
		return true, nil
	}

}
