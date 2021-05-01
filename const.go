package csrf

const (
	SaltsEnvKey = "CSRF_SALTS" // some flexibility in naming the container env key
	TokenLen    = 32
	/*
	   https://pkg.go.dev/golang.org/x/crypto/scrypt

	   The recommended parameters for interactive logins as of 2017 are N=32768, r=8 and p=1.
	   The parameters N, r, and p should be increased as memory latency and CPU parallelism increases;
	   consider setting N to the highest power of 2 you can derive within 100 milliseconds.
	   Remember to get a good random salt.
	*/

	N      = 1 << 15
	R      = 8
	P      = 1
	KeyLen = 32
)
