package fastnum

var (
	int64Zeros [8 * 1024]int64
	int64Ones  = func() (a [seedElems]int64) {
		for i := range a {
			a[i] = 1
		}
		return a
	}()

	float64Zeros [8 * 1024]float64
	float64Ones  = func() (a [seedElems]float64) {
		for i := range a {
			a[i] = 1
		}
		return a
	}()
)
