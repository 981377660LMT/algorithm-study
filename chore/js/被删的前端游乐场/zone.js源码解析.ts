// Zone、ZoneType、ZonePrivate
//
// There can be many different zoneinstances in a system,
// but only one zone is active at any given time which can be retrieved using [Zone#current]
//
// ZoneJS 通过 Monkey patch （猴补丁）的方式，暴力地将浏览器或 Node 中的所有异步 API 进行了封装替换。
// ZoneJS 在加载时，对所有异步接口进行了封装，因此所有在 Zone 中执行的异步方法都会被当做为一个 Task 被其统一监管，
// 并且提供了相应的钩子函数（hooks），用来在异步任务执行前后或某个阶段做一些额外的操作，
// 因此可以实现：记录日志、监控性能、附加数据到异步执行上下文中等。
//
// Zone.current.fork(zoneSpec) // zoneSpec 的类型是 ZoneSpec
//
// // 只有 name 是必选项，其他可选
// interface ZoneSpec {
//   name: string; // zone 的名称，一般用于调试 Zones 时使用
//   properties?: { [key: string]: any; } ; // zone 可以附加的一些数据，通过 Zone.get('key') 可以获取
//   onFork: Function; // 当 zone 被 forked，触发该函数
//   onIntercept?: Function; // 对所有回调进行拦截
//   onInvoke?: Function; // 当回调被调用时，触发该函数
//   onHandleError?: Function; // 对异常进行统一处理
//   onScheduleTask?: Function; // 当任务进行调度时，触发该函数
//   onInvokeTask?: Function; // 当触发任务执行时，触发该函数
//   onCancelTask?: Function; // 当任务被取消时，触发该函数
//   onHasTask?: Function; // 通知任务队列的状态改变
// }
// !用于性能分析：监听异步方法的执行时间
// !在框架中的应用：实现自动重新渲染(用来检查异步任务是否执行完毕，然后触发对应的回调方法)
// Zone 还提供了许多方法来运行、计划和取消任务

import 'zone.js'

{
  const rootZone = Zone.current
  assert(rootZone === Zone.root)

  setTimeout(() => {
    assert(Zone.current === rootZone)
  }, 0)

  // 执行 run 方法，将切换 Zone.current 所保存的 Zone
  const zoneA = rootZone.fork({
    name: 'zoneA',
    onHandleError(parentZoneDelegate, currentZone, targetZone, error) {
      console.log('zoneA.onHandleError', currentZone.name, parentZoneDelegate.zone)
      return false
    }
  })

  zoneA.run(() => {
    assert(Zone.current === zoneA)
    setTimeout(() => {
      assert(Zone.current === zoneA)
    }, 0)
  })

  zoneA.runGuarded(() => {
    throw new Error('zoneA.runGuarded')
  })

  assert(Zone.current === rootZone)
}

function assert(condition: boolean, message?: string): void {
  if (!condition) {
    throw new Error(message)
  }
}
