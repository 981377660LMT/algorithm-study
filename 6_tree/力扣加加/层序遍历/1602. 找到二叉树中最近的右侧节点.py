from collections import deque


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def findNearestRightNode(self, root: TreeNode, u: TreeNode) -> TreeNode:
        ###### bfs波纹法，找u所在层的右兄弟
        queue = deque([root])
        while queue:
            cur_len = len(queue)
            for i in range(cur_len):
                cur = queue.popleft()
                if cur == u:
                    ## u是所在层最右一个，没有右兄弟
                    if i == cur_len - 1:
                        return None
                    ## 有右兄弟
                    return queue[0]
                if cur.left:
                    queue.append(cur.left)
                if cur.right:
                    queue.append(cur.right)

