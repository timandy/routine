package routine

var unset = &object{}

type object struct {
	value Any
}

type threadLocalMap struct {
	table []Any
}

func (mp *threadLocalMap) get(index int) Any {
	lookup := mp.table
	if index < len(lookup) {
		return lookup[index]
	}
	return unset
}

func (mp *threadLocalMap) set(index int, value Any) {
	lookup := mp.table
	if index < len(lookup) {
		lookup[index] = value
		return
	}
	mp.expandAndSet(index, value)
}

func (mp *threadLocalMap) remove(index int) {
	lookup := mp.table
	if index < len(lookup) {
		lookup[index] = unset
	}
}

func (mp *threadLocalMap) expandAndSet(index int, value Any) {
	oldArray := mp.table
	oldCapacity := len(oldArray)
	newCapacity := index
	newCapacity |= newCapacity >> 1
	newCapacity |= newCapacity >> 2
	newCapacity |= newCapacity >> 4
	newCapacity |= newCapacity >> 8
	newCapacity |= newCapacity >> 16
	newCapacity++

	newArray := make([]Any, newCapacity)
	copy(newArray, oldArray)
	fill(newArray, oldCapacity, newCapacity, unset)
	newArray[index] = value
	mp.table = newArray
}

func createInheritedMap() *threadLocalMap {
	parent := currentThread(false)
	if parent == nil {
		return nil
	}
	parentMap := parent.inheritableThreadLocals
	if parentMap == nil {
		return nil
	}
	lookup := parentMap.table
	if lookup == nil {
		return nil
	}
	table := make([]Any, len(lookup))
	copy(table, lookup)
	for i := 0; i < len(table); i++ {
		if c, ok := table[i].(Cloneable); ok {
			table[i] = c.Clone()
		}
	}
	return &threadLocalMap{table: table}
}

func fill(a []Any, fromIndex int, toIndex int, val Any) {
	for i := fromIndex; i < toIndex; i++ {
		a[i] = val
	}
}
