# Zed: 文本编辑中的 Rope 与 SumTree

https://zhuanlan.zhihu.com/p/780572678

1. 比字符串更好的是?

   一些编辑器确实将文本表示为行数组, 每行都用一个字符串表示
   vscode 的 Monaco 编辑器以这种方式工作了很长一段时间: [previous-text-buffer-data-structure](https://code.visualstudio.com/blogs/2018/03/23/text-buffer-reimplementation#_previous-text-buffer-data-structure)
   但字符串数组仍然会受到与单个字符串相同问题的困扰, 过多的内存消耗与性能问题迫使 vscode 团队寻找更好的办法

   比 String 更好的东西是：
   gap buffer, piece table, 与 rope

   - [gap buffer](https://en.wikipedia.org/wiki/Gap_buffer)
   - [piece table](https://en.wikipedia.org/wiki/Piece_table)
   - [rope](<https://en.wikipedia.org/wiki/Rope_(data_structure)>)

   不同编辑器在权衡利弊时做出了不同决定(优化以更适合自己):

   - [emacs/gap-buffer](https://www.gnu.org/software/emacs/manual/html_node/elisp/Buffer-Gap.html)
   - [vscode/piece-table](https://code.visualstudio.com/blogs/2018/03/23/text-buffer-reimplementation)
   - [vim/its-own-tree](https://github.com/vim/vim/blob/master/src/memline.c#L15)
   - [helix/rope](https://github.com/helix-editor/helix/blob/master/docs/architecture.md)

2. 为何 Zed 使用 rope?
   https://github.com/emacs-ng/emacs-ng/issues/378#issuecomment-907680382
   稳定 log
