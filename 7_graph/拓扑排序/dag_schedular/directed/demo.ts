/* eslint-disable @typescript-eslint/no-namespace */

// !调度流程/使用案例

import { Options, Runnable, Schedule, SingleOptionsFn } from '.'
import { DirectedGraph } from './directed-graph/directed-graph'

// 1. 自定义上下文（也可以放在项目的某个全局声明文件中）
declare global {
  namespace Scheduler {
    interface Context {
      userId: string
    }
  }
}

// 2. 定义任务
const task: Runnable = async context => {
  console.log('当前用户 ID: ', context.userId)
  // 业务逻辑...
}

// 3. 构造调度对象
const dag = new DirectedGraph<Parameters<typeof task>[0]>()
const schedule: Schedule = { dag, tags: new Map(), symbols: new Map() }

// 4. 配置调度选项
const configureSchedule: SingleOptionsFn = (options: Options) => {
  // 在这里根据 options.runnable 或 options.tag 添加额外逻辑
  options.dag.addNode(options.runnable!)
}
configureSchedule.__type = 'single'
