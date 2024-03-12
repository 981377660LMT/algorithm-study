// ScapeGoatTree (SGT) 替罪羊树
// https://zhuanlan.zhihu.com/p/180545164
// 为了防止二叉搜索树左右不平衡，我们引入平衡树，而其中思路最简单的是替罪羊树（Scapegoat tree）。
// 替罪羊树 是一种依靠重构操作维持平衡的重量平衡树。
// 替罪羊树会在插入、删除操作时，检测途经的节点，若发现失衡，则将以该节点为根的子树重构。
// 当发现某个子树很不平衡时，暴力重构该子树使之平衡。
// 若左子树或右子树占当前树的比例大于alpha(一般是0.7-0.8) ，则进行重构(很多教程还会判断已删除节点的个数占总大小的比例来决定重不重构)
// 重构分为两步操作：先进行一遍中序遍历，把该子树“拉平”，把其中所有数存入一个数组里（BST的性质决定这个数组一定是有序的）；
// 然后，再用这些数据重新建一个平衡的二叉树，放回原位置。
// 一个节点导致树的不平衡，就要导致整棵子树被拍扁，估计这也是“替罪羊”这个名字的由来吧
//
// !优点：可以维护子树内信息(其他平衡树做不到).
// 缺点：无法持久化、无法维护区间信息.
//
// https://riteme.site/blog/2016-4-6/scapegoat.html
// https://www.nowcoder.com/discuss/353148839920082944
// https://juejin.cn/post/6844904128150241294 使用替罪羊树实现KD-Tree的增删改查
//
// TODO: stack 维护一个内存池(重构操作会频繁的收回一个节点编号，并且再分配节点编号).

package main

// https://github.com/EndlessCheng/codeforces-go/blob/49f6570d86c17f5064a5b079360f4128acc520c4/copypasta/scapegoat_tree.go#L24
func main() {

}

// alpha 的值越小，那么替罪羊树就越容易重构，那么树也就越平衡，查询的效率也就越高，自然修改（加点和删点）的效率也就低了；
// 反之，alpha 的值越大，那么替罪羊树就越不容易重构，那么树也就越不平衡，查询的效率也就越低，自然修改（加点和删点）的效率也就高了。
// 所以，查询多，alpha 就应该小一些；修改多，alpha 就应该大一些。
// alpha = 4/5
const ALPHA_NUM int32 = 4
const ALPHA_DENO int32 = 5

// https://www.luogu.com.cn/problem/P3369
// https://www.luogu.com.cn/problem/P6136
