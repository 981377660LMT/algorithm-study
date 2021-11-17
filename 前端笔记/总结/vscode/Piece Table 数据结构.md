1. 当你在文本编辑器中打开一个文件，内容从硬盘中读取到内存中的数据结构中。如果你想自己开发一个文本编辑器，该如何在内存中存储一个打开的文件？

   1. 我们的第一直觉可能是使用字符串数组：数组中的每个字符串代表文件中的一行， 例如：
      lines = [

   "the quick brown fox", # line 1 of the file

   "jumped over the lazy dog", # line 2 of the file

   ]

   它的确也是一种事实可行的方法，可能你会觉得这种直观性的优势相对于它的潜在缺点是利大于弊。 事实上， VS code 在 2018 年之前，一直是使用这种数据结构。

   然而，当我们遇到大文件时，这种方法就非常消耗性能。至于为什么，考虑一下，如果在文本的中间`插入新的一行`： "went to the park and"
   为了给新的一行创建空间，数组中该行后面每一行都需要在内存中平移。对于大文件来说，这一步操作非常消耗性能。文件越大，意味着更多的数据操作。

   2. 一种 Append-Only 的表示方法
      Piece Table 就是文本编辑器（不管是老的，还是新的）中的强大的数据结构。它的一个关键特征就是，它是`以 append 的方式记录了所有对文件的修改。`
      在文本添加一行文字后，Piece Table 里的内容大概长这样：

      ```JS
        {

        "original": "the quick brown fox\njumped over the lazy dog",

        "add": "went to the park and\n",

        ...

        }
      ```

      当编辑器显示一个打开的文本文件时，它会组合两个 buffer 区的片段，最终形成在屏幕上所见的文本。有些片段会被忽视，例如当用户删除了某些文本。
      到目前为止，编辑器只知道在文本中插入了 `"went to the park and\n"`, 但是是在哪里插入的并不知道。 我们还缺少 Piece Table 的足够信息来正确地展示文本。缺失的部分就是`位置信息`。
      Piece Table 需要跟踪 哪些片段是来自 original buffer，哪些片段是来自 add buffer。Piece Table 是通过遍历一个 Piece descriptor 的列表来完成这个目的。 每一个 piece descriptor 包含有个三个关键信息：

      - source: 说明这个片段，来自哪个 buffer（original/add）.

      - start: 这个片段在该 buffer 中开始的位置。

      - length: 这个片段的长度。

        ```JS
        {

        "original": "the quick brown fox\njumped over the lazy dog",

        "add": "",

        "pieces": [Piece(start=0, length=44, source="original")],

        }
        ```

        我们在第二行添加一段文本后，这个时候，Piece Table 从一个 piece 变成了 3 个 piece， 大概长这个样子：

        ```JS
          {

          "original": "the quick brown fox\njumped over the lazy dog",

          "add": "went to the park and\n",

          "pieces": [

          Piece(start=0, length=20, source="original"),

          Piece(start=0, length=21, source="add"),

          Piece(start=20, length=24, source="original"),

          ],

          }
        ```

        通过在 Piece Table 中连续地读取 piece，文本编辑器可以将一个个 piece 转换成屏幕上所见的内容，并且可以再存回文件。

        当我们从文件中删除部分文本时，我们会把 一个原来的 piece 给拆分成 2 个 piece：

        1. piece 1 指向被删除文本的左侧。

        2. piece 2 指向被删除文本的右侧。

        删除文本仍然存在于 buffer 中，只有没有任何 piece 对它进行引用。

这篇文章把 Piece Table 的数据结构的原理讲得很清楚，但是在实际应用中，这种简单的结构并不能完全满足实际需求，比如文档中文本样式问题并没有考虑，查找速度优化问题等等，我们还需要以此为基础和参考，通过进一步的优化，衍生出更适合实际需求的数据结构。
