优化措施：

1. 防抖：research 上报防抖 100ms，100ms 内如果 search 触发则取消上一次 research
2. 查找策略：searchFast、searchSlow 两种查找方式
3. 缓存：字符串，但没必要
4. 增量计算：

   - 定义 diff 或者 action 的数据结构，要求原子性
     例如增删改/增删

     ```ts
     // diff 风格
     type Diff = {
       removed: string[]
       added: string[]
       colSortChanged: boolean
       rowSortChanged: boolean
       // meta: any // 未知的 diff，放弃增量计算(meta字段表示其他信息，例如放弃增量计算)
     }

     // action 风格
     type Action = {
       type: 'add' | 'remove' | 'update'
       payload: string
     }
     ```

   - changeDetector 监听数据变化, buffer(stash)保存变化，summarize 生成 diff，flush onChange 输出
   - handleChange 中 patchDiff ，patch 可以启发式
   - 注意 diff 过大时放弃增量计算，直接全量计算

5. 分片
6. worker
7. data 计算用 offset，协同计算用 id

---

单元格匹配的数据结构:Map<RowOffset,ColOffset[]>
