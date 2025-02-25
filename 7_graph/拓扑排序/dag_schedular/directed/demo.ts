/* eslint-disable @typescript-eslint/no-namespace */

// !调度流程/使用案例
// 原理：按照拓扑序运行任务.
// ```ts
// export async function run<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, context: T) {
//   for (let i = 0; i < schedule.dag.sorted.length; i++) {
//     const runnable = schedule.dag.sorted[i]
//     const result = runnable(context)
//     if (result instanceof Promise) {
//       // eslint-disable-next-line no-await-in-loop
//       await result
//     }
//   }
// }
// ```

import { Runnable, Schedule } from '.'

// 1. 自定义上下文（也可以放在项目的某个全局声明文件中）
interface IContext {
  userId?: string
}

async function main() {
  // 2. 定义任务
  const fetchUser: Runnable<IContext> = async context => {
    console.log('获取用户信息...')
    context.userId = '123'
  }

  const renderUserId: Runnable<IContext> = context => {
    console.log('当前用户 ID: ', context.userId)
    // 业务逻辑...
  }

  // 3. 构造调度对象
  const schedule = new Schedule<IContext>()
  schedule.add(renderUserId)
  schedule.add(fetchUser, { before: renderUserId })

  // 4. 运行调度
  schedule.build()
  await schedule.run({})

  console.log('任务执行完毕')
}

main()
