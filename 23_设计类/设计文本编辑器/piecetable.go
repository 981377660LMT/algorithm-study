// https://code.visualstudio.com/blogs/2018/03/23/text-buffer-reimplementation#_previous-text-buffer-data-structure
//
// Piece Table 是一种数据结构，用于表示和管理文本编辑器中的文本。
// 它通过维护原始文本和编辑后新增文本的片段（pieces）来高效地支持插入和删除操作，而无需频繁地移动或修改原始文本数据。
