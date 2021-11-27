class Node(object):
    def __init__(self, val=None, children=None):
        self.val = val
        self.children = children


class Codec:
    def serialize(self, root: 'Node') -> str:
        """Encodes a tree to a single string.
        
        :type root: Node
        :rtype: str
        """
        res = []

        def dfs(root):
            if not root:
                return
            res.append(str(root.val))
            res.append(str(len(root.children)))
            for ch in root.children:
                dfs(ch)

        dfs(root)
        # print(res)
        return ",".join(res)

    def deserialize(self, data: str) -> 'Node':
        """Decodes your encoded data to tree.
        
        :type data: str
        :rtype: Node
        """
        if not data:
            return
        data = iter(data.split(","))

        def dfs():
            root = Node(int(next(data)), [])
            num = int(next(data))
            for _ in range(num):
                root.children.append(dfs())
            return root

        return dfs()
