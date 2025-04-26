// https://github.com/981377660LMT/ts/issues/806
//
// ts里的F-Bounded多态
//
// !Java 里可以这样写，ts 可以再简化一点
// interface ITreeNode<T extends ITreeNode<T>> {
//   left?: T
//   right?: T
// }
//
// !typescript 里最好这样写
interface ITreeNode {
  left?: this
  right?: this
}

export {}
