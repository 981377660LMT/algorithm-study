// # from titan_pylib.data_structures.array.persistent_array import PersistentArray
// from typing import List, Iterable, TypeVar, Generic, Optional

// T = TypeVar("T")

// class PersistentArray(Generic[T]):
//     class _Node:
//         __slots__ = "key", "left", "right"

//         def __init__(self, key: T):
//             self.key: T = key
//             self.left: Optional[PersistentArray._Node] = None
//             self.right: Optional[PersistentArray._Node] = None

//         def copy(self) -> "PersistentArray._Node":
//             node = PersistentArray._Node(self.key)
//             node.left = self.left
//             node.right = self.right
//             return node

//     def __init__(
//         self, a: Optional[Iterable[T]] = None, _root: Optional["PersistentArray._Node"] = None
//     ):
//         a = a or []
//         self.root = self._build(a) if _root is None else _root

//     def _build(self, a: Iterable[T]) -> Optional["PersistentArray._Node"]:
//         pool = [PersistentArray._Node(e) for e in a]
//         self.n = len(pool)
//         if not pool:
//             return None
//         n = len(pool)
//         for i in range(1, n + 1):
//             if 2 * i - 1 < n:
//                 pool[i - 1].left = pool[2 * i - 1]
//             if 2 * i < n:
//                 pool[i - 1].right = pool[2 * i]
//         return pool[0]

//     def _new(self, root: Optional["PersistentArray._Node"]) -> "PersistentArray[T]":
//         res = PersistentArray(_root=root)
//         res.n = self.n
//         return res

//     def set(self, k: int, v: T) -> "PersistentArray[T]":
//         assert 0 <= k < len(self), f"IndexError: {self.__class__.__name__}.set({k})"
//         node = self.root
//         if node is None:
//             return self._new(None)
//         new_node = node.copy()
//         res = self._new(new_node)
//         k += 1
//         b = k.bit_length()
//         for i in range(b - 2, -1, -1):
//             if k >> i & 1:
//                 node = node.right
//                 new_node.right = node.copy()
//                 new_node = new_node.right
//             else:
//                 node = node.left
//                 new_node.left = node.copy()
//                 new_node = new_node.left
//         new_node.key = v
//         return res

//     def get(self, k: int) -> T:
//         assert 0 <= k < len(self), f"IndexError: {self.__class__.__name__}.get({k})"
//         node = self.root
//         k += 1
//         b = k.bit_length()
//         for i in range(b - 2, -1, -1):
//             if k >> i & 1:
//                 node = node.right
//             else:
//                 node = node.left
//         return node.key

//     __getitem__ = get

//     def copy(self) -> "PersistentArray[T]":
//         return self._new(None if self.root is None else self.root.copy())

//     def tolist(self) -> List[T]:
//         node = self.root
//         a: List[T] = []
//         if not node:
//             return a
//         q = [node]
//         for node in q:
//             a.append(node.key)
//             if node.left:
//                 q.append(node.left)
//             if node.right:
//                 q.append(node.right)
//         return a

//     def __len__(self):
//         return self.n

//     def __str__(self):
//         return str(self.tolist())

//     def __repr__(self):
//         return f"{self.__class__.__name__}({self})"

package main

func main() {

}
