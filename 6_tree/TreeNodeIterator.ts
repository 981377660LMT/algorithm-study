// /**
//  * 获取两个行之间的所有行的id(包括这两个行).
//  *
//  * @param rootGroup 分组模型.
//  * @param path1 第一个行的路径.
//  * @param rowIndex1 第一个行的索引.
//  * @param path2 第二个行的路径.
//  * @param rowIndex2 第二个行的索引.
//  */
// function getRowIdsBetween(rootGroup: Group, path1: number[], rowIndex1: number, path2: number[], rowIndex2: number): ID[] {
//   const compareRes = compareArray([...path1, rowIndex1], [...path2, rowIndex2])
//   if (compareRes === 1) {
//     ;[path1, path2] = [path2, path1]
//     ;[rowIndex1, rowIndex2] = [rowIndex2, rowIndex1]
//   }

//   // // // 四个状态
//   // const curLeafGroupPath = path1;
//   // let parentGroup: NonLeafGroup | undefined;
//   // const parentGroupGenerator = getGroupByPath(rootGroup, path1.slice(0, -1)!)!.getChildren!();
//   // const curGroup = getLeafGroupByPath(rootGroup, path1);
//   // if (!curGroup) return [];
//   // const curGroupGenerator = curGroup.getRows();

//   // const getRowIdsInGroup = (leafGroup: LeafGroup, fromIndex: number, toIndex: number): ID[] => {
//   //   const rowGenerator = leafGroup.getRows();
//   //   const res: ID[] = [];
//   //   for (let i = fromIndex; i <= toIndex; i++) {
//   //     rowGenerator.moveTo(i);
//   //     res.push(rowGenerator.id);
//   //   }
//   //   return res;
//   // };

//   // const nextLeafGroup = (): void => {
//   //   while (parentGroup) {
//   //     const len = parentGroupGenerator.length;
//   //     const pos = curGroup.path[curGroup.path.length - 1];
//   //     if (pos + 1 < len) {
//   //       parentGroupGenerator.moveTo(pos + 1);
//   //       const nextGroup = parentGroupGenerator.value;
//   //     }
//   //   }

//   //   if (parentGroup) {
//   //   }

//   //   const nextLeafGroupPath = curLeafGroupPath;
//   //   const pos = curGroup.path[curGroup.path.length - 1];
//   //   // if (pos===parentGroupRowGenerator.)
//   // };

//   // let ok = false;
//   // while (ok) {
//   //   if (compareArray([...path1, rowIndex1], [...path2, rowIndex2]) === 0) {
//   //     ok = true;
//   //   }
//   // }

//   // const res: ID[] = [];
//   // function dfs (cur: Group, curPath: number[]): void {}
//   // dfs(rootGroup, []);
//   // return res;
// }

interface ITreeNodeIterator<NoneLeaf, Leaf> {
  readonly value: NoneLeaf | Leaf
  prev(): void
  next(): void
  prevLeaf(): void
  nextLeaf(): void
}

class TreeNodeIterator<NoneLeaf, Leaf> implements ITreeNodeIterator<NoneLeaf, Leaf> {
  prev(): void {
    throw new Error('Method not implemented.')
  }

  next(): void {
    throw new Error('Method not implemented.')
  }

  prevLeaf(): void {
    throw new Error('Method not implemented.')
  }

  nextLeaf(): void {
    throw new Error('Method not implemented.')
  }

  get value(): NoneLeaf | Leaf {}
}

function createNodeIterator<NoneLeaf, Leaf>(root: NoneLeaf, path: number[]): TreeNodeIterator<NoneLeaf, Leaf> {}

export { TreeNodeIterator }
