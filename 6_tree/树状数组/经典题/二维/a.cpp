

//�������ŁB�K�؂ɍċA�̃O���[�o�����E�A�����[�������邱��
struct MultiDimensionalFenwickTree {
	typedef char Index;
	typedef ll Val;
	vector<Index> dim;
	vector<Val> data;
	MultiDimensionalFenwickTree(const vector<Index> &dimension): dim(dimension) {
		int n = 1;
		rep(i, dim.size()) n *= dim[i];
		data.assign(n, 0);
	}
	inline void add(const vector<Index> &indices, Val x) {
		add_rec(indices, 0, 0, x);
	}
	inline Val sum(const vector<Index> &indices) const {
		return sum_rec(indices, 0, 0);
	}
	inline Val sum2(const vector<Index> &a, const vector<Index> &b) const {
		vector<Index> t(a.size());
		return sum2_rec(0, a, b, t);
	}
private:
	void add_rec(const vector<Index> &indices, int k, int t, Val x) {
		int d = dim[k];
		t *= d;
		if(k+1 == dim.size()) {
			Val *p = &data[t];
			for(int i = indices[k]; i < d; i |= i+1)
				p[i] += x;
		}else
			for(int i = indices[k]; i < d; i |= i+1)
				add_rec(indices, k+1, t + i, x);
	}
	Val sum_rec(const vector<Index> &indices, int k, int t) const {
		int d = dim[k];
		t *= d;
		Val res = 0;
		if(k+1 == dim.size()) {
			const Val *p = &data[t];
			for(int i = indices[k]; i > 0; i -= i & -i)
				res += p[i-1];
		}else
			for(int i = indices[k]; i > 0; i -= i & -i)
				res += sum_rec(indices, k+1, t + i - 1);
		return res;
	}
	Val sum2_rec(int d, const vector<Index> &a, const vector<Index> &b, vector<Index> &t) const {
		if(d == dim.size())
			return sum(t);
		Val r = 0;
		t[d] = b[d];
		r += sum2_rec(d+1, a, b, t);
		t[d] = a[d];
		r -= sum2_rec(d+1, a, b, t);
		return r;
	}
};


