# https://nyaannyaan.github.io/library/tree/cartesian-tree.hpp
# // return value : pair<graph, root>
# template <typename T>
# pair<vector<vector<int>>, int> CartesianTree(vector<T> &a) {
#   int N = (int)a.size();
#   vector<vector<int>> g(N);
#   vector<int> p(N, -1), st;
#   st.reserve(N);
#   for (int i = 0; i < N; i++) {
#     int prv = -1;
#     while (!st.empty() && a[i] < a[st.back()]) {
#       prv = st.back();
#       st.pop_back();
#     }
#     if (prv != -1) p[prv] = i;
#     if (!st.empty()) p[i] = st.back();
#     st.push_back(i);
#   }
#   int root = -1;
#   for (int i = 0; i < N; i++) {
#     if (p[i] != -1)
#       g[p[i]].push_back(i);
#     else
#       root = i;
#   }
#   return make_pair(g, root);
# }
from typing import List, Tuple


def cartesianTree(nums: List[int]) -> Tuple[int, List[List[int]]]:
    """给定插入序列构建笛卡尔树,复杂度O(n)

    Args:
        nums (List[int]): 插入序列

    Returns:
        Tuple[int, List[List[int]]]: 根节点, 有向的树

    Notes:
        构建[0,n)的过程:
            1.当要构建 [left,right) 区间时，找到[left,right) 里值最小的元素的小标 minIndex.
            2.将minIndex作为根, [left,minIndex) 作为左子树, [minIndex+1,right) 作为右子树，继续构建.
    """
    n = len(nums)
    tree = [[] for _ in range(n)]
    parent = [-1] * n
    stack = []
    for i in range(n):
        last = -1
        while stack and nums[stack[-1]] > nums[i]:
            last = stack.pop()
        if stack:
            parent[i] = stack[-1]
        if last != -1:
            parent[last] = i
        stack.append(i)

    root = -1
    for i in range(n):
        if parent[i] != -1:
            tree[parent[i]].append(i)
        else:
            root = i
    return root, tree


print(cartesianTree([9, 3, 7, 1, 8, 12, 10, 20, 15, 18, 5]))
