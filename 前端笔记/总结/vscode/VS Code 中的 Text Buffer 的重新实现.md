它是基于传统的 piece table 数据结构的一次优化改造。

1. 之前的 text buffer 数据结构
   编辑器的核心模型是基于文本行的，例如，开发者是一行一行地读写代码，编译器提供的运行时诊断/堆栈追踪提供了行数和列数，tokenization 引擎是一行一行地运行
   开始的时候，我们使用的是以行为单位的数组的方法，并且运行良好，因为常规的文本输入都相对较小。
   当用户正在输入时，我们在数组中定位到用户正在输入的行，并且修改这一行的字符串进行替换；当用户插入新的一行时，我们在行数组( line array) 中插入一个新的行对象，由 JS 引擎帮我们完成繁重的底层内存操作。

   但是，我们在不断地收到 VS Code 的崩溃反馈：当打开某些文件的时候会导致内存不足。例如，用户打开一个 35 MB 的文件失败。问题原因在于，这个文件里有太多的行，1370 万行。`我们之前为每一行创建了一个 ModelLine 对象，每个对象大概占用了 40 - 60 字节，所以 line array 使用了大约 600 MB 的内存来存放文档。这个内存大小是文件原始大小的 20 倍。`

2. 新的 text buffer 实现
   我们开始寻找一种数据结构，占用尽可能少的 metadata。在评估完一些数据结构以后，我发现 piece table 可能是一个很好的备选项。

```TS
class PieceTable {

    original: string; // original contents

    added: string; // user added contents

    nodes: Node[];

}

class Node {

    type: NodeType;

    start: number;

    length: number;

}

enum NodeType {

    Original,

    Added

}
```

piece table 的初始内存大小，十分接近原始文件的大小，编辑动作需要的内容正比于编辑动作的数量和新增文本的大小。所以，一般来说，piece table 在内存方面具有巨大的优势。 但是，低内存的代价就是，访问一个行（真实的行）的速度慢。例如，`如果你想要得到第 1000 行的内容，唯一的办法是从文档的开始的每一个字符开始遍历，从第 999 个换行符开始，读取字符，直到下一个换行符号。`

3. 使用 caching 进行更快的行查找
   传统的 piece table 的 nodes 中只包含了偏移量，但是我们可以在里面添加换行信息，从而可以更快的行查找。

```JS
class Node {

type: NodeType;

start: number;

length: number;

lineStarts: number[]; // 存储换行符的 offset

}
```

例如，如果你想从一个给定的 Node 中找到第二行，你可以从 `node.lineStarts[0]` and `node.lineStarts[1]` 的相对偏移量中读取到`这一行的文本的位置。` (差分数组)

4. 使用平衡二叉树(红黑树)来加速行查找
   传统的 piece table 的 Node 按顺序存在 list 中，他的顺序就代表 Node 在文档中的顺序
   通过红黑树将 list 结构转为平衡二叉让查找行的时间复杂度降到了 logn

   1. 通过数据结构和算法来优化性能，永远都是最优解，而不遇到问题就抱怨 javascript 语言本身。

   2. 任何数据过长的时候，应该意识到潜在可能的性能问题。

5. Piece Tree
