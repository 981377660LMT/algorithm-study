from typing import DefaultDict, List, Optional, Set, Tuple
from collections import defaultdict


INF = int(1e20)


# Definition for a binary tree node.
class TreeNode:
    def __init__(self, x):
        self.val = x
        self.left = None
        self.right = None


# 2 <= 树中节点个数 <= 10^5
# -10000 <= 树中节点的值 <= 10000
# 树上启发式合并
# 为了保证总复杂度是 O(nlogn)，需要把小子树合并到大子树里，也就是启发式合并。
class Solution:
    def getMaxLayerSum(self, root: Optional[TreeNode]) -> int:
        ...


##################################
import sys
from typing import Optional

sys.setrecursionlimit(9999999)

# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, x):
#         self.val = x
#         self.left = None
#         self.right = None


# 暴力过了
# class Num:
#     def __init__(self):
#         self.val = 0

#     def setV(self, n):
#         self.n = n

#     def raiseV(self, n):
#         if n>self.n:
#             self.n = n

#     def getV(self):
#         return self.n

# class Solution:
#     def getMaxLayerSum(self, root: Optional[TreeNode]) -> int:
#         def getSum(node):
#             if node is None:
#                 return [],[]
#             tl,cl = getSum(node.left)
#             tr,cr = getSum(node.right)
#             r = [node.val]
#             c = [1]
#             r.extend([0]*max(len(tl),len(tr)))
#             c.extend([0]*max(len(tl),len(tr)))
#             for i, n in enumerate(tl):
#                 r[i+1] += n
#                 c[i+1] += cl[i]
#             for i, n in enumerate(tr):
#                 r[i+1] += n
#                 c[i+1] += cr[i]
#             node.array = r
#             node.cts = c
#             return r, c
#         getSum(root)

#         ans = Num()
#         ans.setV(max(root.array))

#         def getMax(ra, rc, na, nc, df):
#             ba = []
#             ba.extend(ra)
#             for i, v in enumerate(na):
#                 ba[i+df] -= na[i]
#                 if i>0:
#                     ba[i+df-1] += na[i]
#             # lowest value
#             if len(ra)==len(na)+df and rc[-1]==nc[-1]:
#                 if len(ba)>0:
#                     ans.raiseV(max(ba[:-1]))
#             else:
#                 ans.raiseV(max(ba))

#         def getCount(node, depth):
#             if node is None:
#                 return
#             if node.left is not None and node.right is not None:
#                 getCount(node.left, depth+1)
#                 getCount(node.right, depth+1)
#             else:
#                 getMax(root.array, root.cts, node.array, node.cts, depth)
#                 getCount(node.left, depth+1)
#                 getCount(node.right, depth+1)

#         getCount(root, 0)
#         return ans.getV()

