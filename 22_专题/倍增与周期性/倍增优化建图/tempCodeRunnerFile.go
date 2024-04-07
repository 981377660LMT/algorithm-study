
	n, log := d.n, d.log
	for k := log; k >= 0; k-- {
		if cur == -1 {
			return
		}
		if len&(1<<k) != 0 {
			f(k*n + cur)
			cur = d.jump[k*n+cur]
		}
	}
	f(cur)