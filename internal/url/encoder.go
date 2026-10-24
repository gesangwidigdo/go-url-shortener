package url

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Encode(id uint64, offset uint64) (string, error) {
	n := uint64(id) + offset
	if n == 0 {
		return string(alphabet[0]), nil
	}

	var buf [11]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = alphabet[n%62]
		n /= 62
	}

	return string(buf[i:]), nil
}
