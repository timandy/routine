package routine

type entry struct {
	value Any
}

type threadLocalMap struct {
	table []*entry
}

func (mp *threadLocalMap) getEntry(key ThreadLocal) *entry {
	index := key.Id()
	if index < len(mp.table) {
		return mp.table[index]
	}
	return nil
}

func (mp *threadLocalMap) set(key ThreadLocal, value Any) {
	index := key.Id()
	if index < len(mp.table) {
		e := mp.table[index]
		if e == nil {
			mp.table[index] = &entry{value: value}
			// try restart gc timer if Set for the first time
			gcTimerStart()
		} else {
			e.value = value
		}
		return
	}

	newCapacity := index
	newCapacity |= newCapacity >> 1
	newCapacity |= newCapacity >> 2
	newCapacity |= newCapacity >> 4
	newCapacity |= newCapacity >> 8
	newCapacity |= newCapacity >> 16
	newCapacity++

	newEntries := make([]*entry, newCapacity)
	copy(newEntries, mp.table)
	newEntries[index] = &entry{value: value}
	mp.table = newEntries
	// try restart gc timer if Set for the first time
	gcTimerStart()
}

func (mp *threadLocalMap) remove(key ThreadLocal) {
	index := key.Id()
	if index < len(mp.table) {
		mp.table[index] = nil
	}
}

// BackupContext copy all local table into an threadLocalMap instance.
func createInheritedMap() *threadLocalMap {
	parent := currentThread(false)
	if parent == nil {
		return nil
	}
	mp := parent.inheritableThreadLocals
	if mp == nil || mp.table == nil {
		return nil
	}
	table := make([]*entry, len(mp.table))
	copy(table, mp.table)
	return &threadLocalMap{table: table}
}
