
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]