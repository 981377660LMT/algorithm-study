aigc.sankuai.com/assistant/share/eebb7f397e93c1a7758d2253b5a7a2bb

---

- go
  **https://github.com/noneback/go-taskflow**
  https://github.com/fieldryand/goflow
  https://github.com/AkihiroSuda/go-dag

- ts
  https://github.com/walleXD/ts-dag
  https://github.com/markschad/ms-dag-ts
  https://github.com/PauloMigAlmeida/directed-acyclic-graph-builder-js

---

https://developer.aliyun.com/article/625843
流程调度系统需求

1. 建图：

   - 用户指定任务节点的依赖关系，算法自动构建 DAG，组装流程节点；

   ```ts
   const dagScheduler = new DAGScheduler<string>()
   const runner1 = dagScheduler.add('id1', {
     deps: ['id2', 'id3'],
     onTrigger: () => {
       console.log('task1')
     },
     onReset: () => {
       console.log('task1 reset')
     },
     onError: error => {
       console.error('task1 error', error)
     }
   })

   dagScheduler.run()
   ```

2. 触发/重新执行任务：

   - 用户可以手动触发任务(重新执行)；当任务重新执行后，由于依赖发生变化，所有子任务需要被重置；
   - 手动触发的任务如果依赖的前置任务未完成，不会执行当前节点任务；

3. 自动运行：
   - 当任务依赖的前置任务完成后，自动执行当前节点任务；

帮我分析下这个需求，背景是表单场景联动，我的需求是否合理，能否给出评价和实现方案、优化方案
