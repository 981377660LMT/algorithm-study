from typing import List, Tuple


def tree_dp_recursively(n: int, tree: List[List[int]]) -> Tuple[List[int], List[int]]:
    order = []
    st = [0]
    parent = [-1] * n
    # depth = [0] * n
    while st:
        v = st.pop()
        order.append(v)
        for c in tree[v]:
            if c != parent[v]:
                st.append(c)
                parent[c] = v
                # depth[c] = depth[v] + 1
    return parent, order[::-1]
