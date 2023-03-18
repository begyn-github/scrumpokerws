package myutils

func fibo(i int) int {
	var resp int

	switch i {
	case 0:
		resp = 0
	case 1:
		resp = 1
	default:
		resp = fibo(i-1) + fibo(i-2)
	}

	return resp
}

func Phi(i int) float32 {
	var a float32 = float32(fibo(i))
	var b float32 = float32(fibo(i + 1))
	return a / b
}
