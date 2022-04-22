[kd 树算法之思路篇](https://www.joinquant.com/view/community/detail/dd60bd4e89761b916fe36dc4d14bb272)
每一个节点都是按照上下或者左右进行平分的，因此如果两个点在树中的距离较近，那么它们的实际距离就是比较近的。
kd 树是一个二叉树结构，它的每一个节点记载了【特征坐标，切分轴，指向左枝的指针，指向右枝的指针】

- Python 的 scikit-learn 机器学习包提供了蛮算、kd 树和 ball 树三种 kNN 算法
- 给定一堆已有的样本数据，和一个被询问的数据点（红色五角星），我们如何找到离五角星最近的 15 个点？

https://leetcode-cn.com/problems/vFjcfV/solution/kd-treeban-zi-ti-by-mo-yan-24-63mv/
