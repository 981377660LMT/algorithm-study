// // https://codeforces.com/contest/1638/problem/E
// // 持つ値のタイプ T、座標タイプ X
// // コンストラクタでは T none_val を指定する
// template <typename T, typename X = ll>
// struct Intervals {
//   static constexpr X LLIM = -infty<X>;
//   static constexpr X RLIM = infty<X>;
//   const T none_val;
//   // none_val でない区間の個数と長さ合計
//   int total_num;
//   X total_len;
//   map<X, T> dat;

import { SortedList } from '../../22_专题/离线查询/根号分治/SortedList/SortedList'

//   Intervals(T none_val) : none_val(none_val), total_num(0), total_len(0) {
//     dat[LLIM] = none_val;
//     dat[RLIM] = none_val;
//   }

//   // x を含む区間の情報の取得
//   tuple<X, X, T> get(X x, bool ERASE) {
//     auto it2 = dat.upper_bound(x);
//     auto it1 = prev(it2);
//     auto [l, tl] = *it1;
//     auto [r, tr] = *it2;
//     if (tl != none_val && ERASE) {
//       --total_num, total_len -= r - l;
//       dat[l] = none_val;
//       merge_at(l);
//       merge_at(r);
//     }
//     return {l, r, tl};
//   }

//   // [L, R) 内の全データの取得
//   template <typename F>
//   void enumerate_range(X L, X R, F f, bool ERASE) {
//     assert(LLIM <= L && L <= R && R <= RLIM);
//     if (!ERASE) {
//       auto it = prev(dat.upper_bound(L));
//       while ((*it).fi < R) {
//         auto it2 = next(it);
//         f(max((*it).fi, L), min((*it2).fi, R), (*it).se);
//         it = it2;
//       }
//       return;
//     }
//     // 半端なところの分割
//     auto p = prev(dat.upper_bound(L));
//     if ((*p).fi < L) {
//       dat[L] = (*p).se;
//       if (dat[L] != none_val) ++total_num;
//     }
//     p = dat.lower_bound(R);
//     if (R < (*p).fi) {
//       T t = (*prev(p)).se;
//       dat[R] = t;
//       if (t != none_val) ++total_num;
//     }
//     p = dat.lower_bound(L);
//     while (1) {
//       if ((*p).fi >= R) break;
//       auto q = next(p);
//       T t = (*p).se;
//       f((*p).fi, (*q).fi, t);
//       if (t != none_val) --total_num, total_len -= (*q).fi - (*p).fi;
//       p = dat.erase(p);
//     }
//     dat[L] = none_val;
//   }

//   void set(X L, X R, T t) {
//     enumerate_range(
//         L, R, [](int l, int r, T x) -> void {}, true);
//     dat[L] = t;
//     if (t != none_val) total_num++, total_len += R - L;
//     merge_at(L);
//     merge_at(R);
//   }

//   template <typename F>
//   void enumerate_all(F f) {
//     enumerate_range(LLIM, RLIM, f, false);
//   }

//   void merge_at(X p) {
//     if (p == LLIM || RLIM == p) return;
//     auto itp = dat.lower_bound(p);
//     assert((*itp).fi == p);
//     auto itq = prev(itp);
//     if ((*itp).se == (*itq).se) {
//       if ((*itp).se != none_val) --total_num;
//       dat.erase(itp);
//     }
//   }
// };

// TODO map
const INF = 2e15

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `SortedList`实现.
 */
class ODTMap<S> {
  private _len = 0
  private _count = 0
  private readonly _leftLimit = -INF
  private readonly _rightLimit = INF
  private readonly _data: SortedList<[start: number, end: number, value: S]> = new SortedList()
  private readonly _noneValue: S

  /**
   * 指定哨兵值建立一个ODTMap.
   * @param noneValue 表示空值的哨兵值.
   */
  constructor(noneValue: S) {
    this._noneValue = noneValue
  }

  /**
   * 返回包含`x`的区间的信息.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] {}

  set(start: number, end: number, value: S): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
  }

  enumerateAll(f: (start: number, end: number, value: S) => void): void {
    this.enumerateRange(this._leftLimit, this._rightLimit, f, false)
  }

  /**
   * 遍历范围`[start, end)`内的所有区间.
   */
  enumerateRange(
    start: number,
    end: number,
    f: (start: number, end: number, value: S) => void,
    erase = false
  ): void {
    if (start < this._leftLimit) start = this._leftLimit
    if (end > this._rightLimit) end = this._rightLimit
    if (start >= end) return
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateAll((start, end, value) => {
      const v = value === this._noneValue ? 'null' : value
      sb.push(`[${start},${end}):${v}`)
    })
    return `ODTMap{${sb.join(', ')}}`
  }

  /**
   * 区间个数.
   */
  get length(): number {
    return this._len
  }

  /**
   * 区间内元素个数之和.
   */
  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p <= 0 || this._rightLimit <= p) return
  }
}

export { ODTMap }
