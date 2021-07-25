class BST:
    def __init__(self) -> None:
        pass

    def inorder(self):
        print(1)


class AVL(BST):
    def __init__(self) -> None:
        super().__init__()

    def inorder(self, other):
        return other


avl = AVL()
print(avl.inorder(1))

