
// Source: https://codeforces.com/blog/entry/83248
// T must be a commutative ring
// Given some function op:\mathbb{Z}^2->T,
// each query asks for \sum_{l<=i<j<r} op(i,j)
//
// solve() initializes a multiset S to be empty.
// Each insert(x) call should insert x to S.
// Each query_left(x) should return \sum_{y \in S} op(x, y).
// Each query_right(x) should return \sum_{y \in S} op(y, x).
//
// O(n) insert() calls
// O(n + qn*BX) query_left() calls, and
// O(n + n^2/BX) query_right() calls.
//
// Set BX = n/sqrt(qn) to achieve
// O(n) insert() calls, and
// O(n * sqrt(qn)) query_left() and query_right() calls.
template<int BX, class T, class F = plus<>, class I = minus<>>
struct mo_sweepline_base{
	int n;
	mutable vector<array<int, 3>> q;
	F TT;
	T T_id;
	I TinvT;
	mo_sweepline_base(int n, F TT = plus<>(), T T_id = 0, I TinvT = minus<>()): n(n), TT(TT), T_id(T_id), TinvT(TinvT){ }
	void query(int qi, int ql, int qr){
		assert(0 <= ql && ql <= qr && qr <= n);
		q.push_back({ql, qr, qi});
	}
	vector<T> solve(auto insert, auto query_left, auto query_right) const{
		sort(q.begin(), q.end(), [&](auto x, auto y){ return x[0] / BX != y[0] / BX ? x[0] < y[0] : x[0] / BX & 1 ? x[1] > y[1] : x[1] < y[1]; });
		vector<vector<array<int, 4>>> update(n + 1);
		int l = 0, r = 0;
		// [0-bit] 0: query_left, 1: query_right
		// [1-bit] 0: add       , 1: subtract
		for(auto [ql, qr, qi]: q){
			if(ql < l){
				update[r].push_back({ql, l, 2, qi});
				l = ql;
			}
			if(r < qr){
				update[l].push_back({r, qr, 1, qi});
				r = qr;
			}
			if(l < ql){
				update[r].push_back({l, ql, 0, qi});
				l = ql;
			}
			if(qr < r){
				update[l].push_back({qr, r, 3, qi});
				r = qr;
			}
		}
		vector<T> pref_exc(n + 1), pref_inc(n + 1);
		// pref_exc[r] = \sum_{0<=i< j<r} op(i,j)
		// pref_inc[r] = \sum_{0<=i<=j<r} op(j,i)
		int qn = (int)q.size();
		vector<T> res(qn, T_id);
		for(auto x = 0; x <= n; ++ x){
			for(auto [from, to, coef, qi]: update[x]){
				T sum = T_id;
				if(coef & 1) for(auto i = from; i < to; ++ i) sum = TT(sum, query_right(i));
				else for(auto i = from; i < to; ++ i) sum = TT(sum, query_left(i));
				res[qi] = coef & 2 ? TinvT(res[qi], sum) : TT(res[qi], sum);
			}
			if(x < n){
				pref_exc[x + 1] = TT(pref_exc[x], query_right(x));
				insert(x);
				pref_inc[x + 1] = TT(pref_inc[x], query_left(x));
			}
		}
		//  \sum_{0<=i<=j<l} op(j,i) +
		//  \sum_{0<=i<l,i<j<r} op(i,j)
		T aux = T_id;
		for(auto [ql, qr, qi]: q){
			aux = TT(aux, res[qi]);
			res[qi] = TinvT(TT(pref_inc[ql], pref_exc[qr]), aux);
		}
		return res;
	}
};
template<class T, class F = plus<>, class I = minus<>>
using mo_sweepline = mo_sweepline_base<500, T, F, I>;


// Requires sqrt_decomposition_heavy_point_update_light_range_query_commutative_group and mo_sweepline
template<int BX>
struct range_inversion_query_solver_offline{
	int n;
	vector<int> data;
	mo_sweepline_base<BX, long long, plus<>, minus<>> mo;
	template<class T, class Compare = less<>>
	range_inversion_query_solver_offline(const vector<T> &a, Compare cmp = less<>()): n((int)a.size()), data((int)a.size()), mo(n, plus<>(), 0LL, minus<>()){
		vector<T> temp = a;
		sort(temp.begin(), temp.end(), cmp);
		for(auto i = 0; i < n; ++ i) data[i] = lower_bound(temp.begin(), temp.end(), a[i]) - temp.begin();
	}
	void query(int qi, int ql, int qr){
		mo.query(qi, ql, qr);
	}
	// O(n*BX + qn*BX + n^2/BX)
	vector<long long> solve() const{
		sqrt_decomposition_heavy_point_update_light_range_query_commutative_group<BX, int, plus<>, minus<>> sqrtdecomp(n, plus<>(), 0, minus<>());
		auto insert = [&](int i)->void{
			sqrtdecomp.update(data[i], 1);
		};
		auto query_left = [&](int i)->int{
			return sqrtdecomp.query(0, data[i]);
		};
		auto query_right = [&](int i)->int{
			return sqrtdecomp.query(data[i] + 1, n);
		};
		return mo.solve(insert, query_left, query_right);
	}
};

int main(){
	cin.tie(0)->sync_with_stdio(0);
	cin.exceptions(ios::badbit | ios::failbit);
	int n, qn;
	cin >> n >> qn;
	vector<int> a(n);
	copy_n(istream_iterator<int>(cin), n, a.begin());
	range_inversion_query_solver_offline<320> rinvq(a);
	for(auto qi = 0; qi < qn; ++ qi){
		int l, r;
		cin >> l >> r;
		rinvq.query(qi, l, r);
	}
	for(auto x: rinvq.solve()){
		cout << x << "\n";
	}
	return 0;
}
