from typing import List, Tuple


def getSubtreeInfo(
    tree: List[List[int]], root: int
) -> Tuple[List[int], List[int], List[int], List[int]]:
    """获取子树信息.height[i] 表示以 i 为根的子树的高度(距离最远的叶子节点的距离)."""
    n = len(tree)
    parent = [0] * n
    depth = [0] * n
    subsize = [0] * n
    height = [0] * n
    topological = [0] * n
    topological[0] = root
    parent[root] = root
    depth[root] = 0
    left, right = 0, 1
    while left < right:
        cur = topological[left]
        for next in tree[cur]:
            if next != parent[cur]:
                topological[right] = next
                right += 1
                parent[next] = cur
                depth[next] = depth[cur] + 1
        left += 1
    right -= 1
    while right >= 0:
        cur = topological[right]
        subsize[cur] = 1
        height[cur] = 0
        for next in tree[cur]:
            if next != parent[cur]:
                subsize[cur] += subsize[next]
                tmp = height[next] + 1
                if tmp > height[cur]:
                    height[cur] = tmp
        right -= 1
    parent[root] = -1
    return parent, depth, subsize, height


if __name__ == "__main__":
    # tree := [][]int32{{1, 2}, {0}, {0, 3, 4}, {2}, {2, 5}, {4}}
    tree = [[1, 2], [0], [0, 3, 4], [2], [2, 5], [4]]
    root = 0
    parent, depth, subsize, height = getSubtreeInfo(tree, root)
    print(parent, depth, subsize, height)
