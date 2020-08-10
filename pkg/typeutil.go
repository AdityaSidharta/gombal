package pkg

func Find(values []string, val string) (int, bool){
	for i, value := range values {
		if value == val {
			return i, true
		}
	}
	return -1, false
}
