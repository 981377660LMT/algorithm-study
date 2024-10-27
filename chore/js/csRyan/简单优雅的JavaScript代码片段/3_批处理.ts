// https://segmentfault.com/a/1190000043866551
// 合并请求，成批发出
//
// 场景：
// !后端提供的接口具备批量查询能力（比如，同时查询10个资源的详情信息），但是调用者每次只需要请求一个资源的详情信息。
// !没有将接口本身的批量请求能力利用起来，发出过多请求，导致接口流控限制或效率低下。
//
// import React from "react";
//
// const List = ({ ids }: { ids: string[] }) => {
//   return (
//     <div>
//       {ids.map((id) => (
//         <ResourceDetail key={id} id={id} />
//       ))}
//     </div>
//   );
// };
//
// const ResourceDetail = ({ id }: { id: string }) => {
//   // 在这个组件请求资源详情并渲染...
//   // 注意在这个组件中，你只关心这一个资源的信息，你不再具备”整个列表“的数据视角。
// };
//
// // 后端接口支持同时请求多个资源的信息
// declare function api(ids: string[]): Promise<Record<string, { details: any }>>;
//
// 实现原理：
// !本质上是一个请求的缓冲队列（buffer），虽然它一个一个地接受请求，但是它不会立即将请求发出，而是等待一段时间(debounce)，将请求累积起来，然后将累积起来的请求成批发出。
// 通过callback的方式来返回数据，是因为要支持分批多次返回（拿到一批响应就立刻返回给调用者），而不是等待所有结果到达以后再全部一次性返回
// 这里的 callback 可以当作双向通信的通道.

import { debounce } from 'lodash'

/**
 * 批处理请求.
 *
 * 将「批量请求」接口封装成「逐个请求」的接口.
 * 使用者无需关心请求是如何被合并发出的.
 */
function withBatch<InputItem, OutputItem, Key extends PropertyKey>(
  fn: (inputs: InputItem[]) => Promise<OutputItem[]>,
  keyFn: (input: InputItem) => Key,
  options?: { debounceWait?: number; maxWait?: number }
): (input: InputItem) => Promise<OutputItem> {
  type BufferItem = {
    key: Key
    value: InputItem
    promise: Promise<OutputItem>
    resolve: (value: OutputItem) => void
    reject: (error: unknown) => void
  }

  const { debounceWait = 50, maxWait = 200 } = options || {}
  const buffer: Map<Key, BufferItem> = new Map()

  const excute = debounce(
    () => {
      if (!buffer.size) return
      const batch = new Map(buffer)
      const inputs = [...batch.values()]
      buffer.clear()
      fn(inputs.map(({ value }) => value))
        .then(results => {
          for (const [index, result] of results.entries()) {
            const { resolve } = inputs[index]
            resolve(result)
          }
        })
        .catch(error => {
          for (const { reject } of inputs) {
            reject(error)
          }
        })
    },
    debounceWait /** 等待若干毫秒作为一个请求收集窗口，然后将收集到的所有请求作为一批发出 */,
    {
      maxWait /** 避免不断有请求到来，导致debounce一直无法被调用 */,
      leading: false,
      trailing: true
    }
  )

  const schedule = (value: InputItem): Promise<OutputItem> => {
    const key = keyFn(value)
    const item = buffer.get(key)
    // 如果已经有相同的input在buffer中，则不重复调度它，而是与前一个input共享同一个结果
    if (item) return item.promise
    const bufferItem = { key, value, ...createControllablePromise<OutputItem>() }
    buffer.set(key, bufferItem)
    excute()
    return bufferItem.promise
  }

  return schedule

  function createControllablePromise<T>(): {
    promise: Promise<T>
    resolve: (value: T) => void
    reject: (reason: unknown) => void
  } {
    let res: any = {}
    res.promise = new Promise<T>((resolve, reject) => {
      res.resolve = resolve
      res.reject = reject
    })
    return res
  }
}

export default {
  withBatch
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const api = (ids: string[]): Promise<Record<string, { details: any }>[]> => {
    console.log('called', ids)
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(ids.map(id => ({ [id]: { details: id } })))
      }, 1000)
    })
  }
  const wrappedApi = withBatch(api, id => id)

  wrappedApi('1').then(console.log)
  wrappedApi('1').then(console.log)
  wrappedApi('2').then(console.log)
  wrappedApi('2').then(console.log)
  wrappedApi('1').then(console.log)
}
