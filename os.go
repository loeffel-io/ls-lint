package main

func normalizeConfig(bytes []byte, unixSep byte, sep byte) []byte {
	var forbidden = false

	for i, b := range bytes {
		switch b {
		case unixSep:
			if !forbidden && unixSep != sep {
				bytes[i] = sep
			}
		case byte(':'):
			forbidden = true
		case byte('\n'):
			forbidden = false
		}
	}

	return bytes
}
