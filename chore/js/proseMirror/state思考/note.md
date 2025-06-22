# pm 范式

1. transform

   - **能力**：Transform 提供了一种方式来描述和应用对文档的变更，如插入文本、删除节点等。这些变更通过事务（Transaction）来表示，每个事务可以包含多个步骤（Step），每个步骤描述了一个具体的变更操作。
   - **依赖**：Transform 依赖于 Model，因为它需要知道文档的结构来正确地应用变更。同时，Transform 的结果（事务）通常会被应用到 State 上，导致 State 的更新。

2. history
   ```ts
   import { history, undo, redo } from 'prosemirror-history'
   ```
3. collab
