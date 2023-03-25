// https://suisen-cp.github.io/cp-library-cpp/library/datastructure/segment_tree/segment_tree_beats.hpp
// https://rsm9.hatenablog.com/entry/2021/02/01/220408

// 吉老师线段树
// 0. 在propagete函数里加一句.

// void propagete(int k, Id lazy) {
// 	d[k] = mapping(f, d[k]);
// 	if (k < size) {
// 			lz[k] = composition(f, lz[k]);
// !		if (d[k].fail) pushDown(k), pushUp(k);
// 	}
// }

// 1.SegmentTreeBeats 的 E 需要带有fail属性.
// 2.mapping函数更新`非叶子结点的` e 时，如果 e 持有的情报不足以更新导致更新失败，
//   需要将mapping返回值的fail属性设置为true.
//   叶子结点:size==1
// 3.mapping函数之外的更新不允许失败.

package main

func main() {

}

// 子数组mex之和
// https://yukicoder.me/submissions/652366
