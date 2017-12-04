package util

func SubString(str string, begin, end int) string {
	bytes := []byte(str)
	lenght := len(bytes)
	if begin < 0 {
		begin = 0
	}
	if end > lenght {
		end = lenght
	}
	return string(bytes[begin:end])
}