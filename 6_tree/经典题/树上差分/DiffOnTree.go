// TODO 树上差分
// https://blog.csdn.net/justidle/article/details/104508212
// 边差分、点差分

// 树上差分，就是利用差分的性质，对路径上的重要节点进行修改（而不是暴力全改），
// 作为其差分数组的值，最后在求值时，利用 dfs 遍历求出差分数组的前缀和，
// 就可以达到降低复杂度的目的。树上差分时需要求 LCA

// update
// build
// get
