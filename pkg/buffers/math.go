package buffers

func Modulus(num, mod int) int {
	num = num % mod
	if num < 0 {
		num += mod
	} else if num >= mod {
		num -= mod
	}
	return num
}
