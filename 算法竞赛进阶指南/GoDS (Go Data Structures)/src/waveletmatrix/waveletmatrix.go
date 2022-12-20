// waveletmatrix
package waveletmatrix

func main() {

}

func NewWaveletMatrix(data []int) *WaveletMatrix {
	wm := &WaveletMatrix{}
	wm.data = data
	wm.n = len(data)
	wm.max = 0
	for _, v := range data {
		if v > wm.max {
			wm.max = v
		}
	}
	wm.max++
	wm.build()
	return wm
}
