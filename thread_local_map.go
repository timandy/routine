package routine

var unset entry = &object{}

type object struct {
	none bool //nolint:unused
}

type threadLocalMap struct {
	table []entry
}

func (mp *threadLocalMap) get(index int) entry {
	lookup := mp.table
	if index < len(lookup) {
		return lookup[index]
	}
	return unset
}

func (mp *threadLocalMap) set(index int, value entry) {
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

func (mp *threadLocalMap) expandAndSet(index int, value entry) {
	oldArray := mp.table
	oldCapacity := len(oldArray)
	newCapacity := index
	newCapacity |= newCapacity >> 1
	newCapacity |= newCapacity >> 2
	newCapacity |= newCapacity >> 4
	newCapacity |= newCapacity >> 8
	newCapacity |= newCapacity >> 16
	newCapacity++

	newArray := make([]entry, newCapacity)
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
	table := make([]entry, len(lookup))
	copy(table, lookup)
	for i := 0; i < len(table); i++ {
		if c, ok := entryAssert[Cloneable](table[i]); ok && !isNil(c) {
			table[i] = entry(c.Clone())
		}
	}
	return &threadLocalMap{table: table}
}

func fill[T any](a []T, fromIndex int, toIndex int, val T) {
	for i := fromIndex; i < toIndex; i++ {
		a[i] = val
	}
}
