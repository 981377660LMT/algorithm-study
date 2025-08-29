/**
 * 分组并发。
 * 对一组项目进行分组，并以“组间并发，组内串行”的方式处理它们。
 *
 * @template T 项目的类型。
 * @template K 分组键的类型 (string | number)。
 * @template R 每个组处理后返回结果的类型。
 *
 * @param items 要处理的项目数组。
 * @param getKey 一个函数，用于从每个项目中提取分组键。
 * @param processGroup 一个异步函数，负责按顺序处理单个组内的所有项目。
 * @returns 一个 Promise，当所有组都处理完毕后，它会解析为包含每个组结果的数组。
 */
async function groupByAndProcessConcurrently<T, K extends PropertyKey, R>(
  items: T[],
  getKey: (index: number) => K,
  processGroup: (group: T[], key: K) => Promise<R>
): Promise<R[]> {
  // 1. 根据 getKey 函数提供的逻辑对项目进行分组。
  const groups = new Map<K, T[]>()
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    const key = getKey(i)
    if (!groups.has(key)) groups.set(key, [])
    groups.get(key)!.push(item)
  }

  // 2. 为每个组创建一个处理任务（Promise）。
  const tasks = Array.from(groups.entries()).map(([key, group]) => {
    // 调用传入的 processGroup 函数，它定义了组内的串行处理逻辑。
    return processGroup(group, key)
  })

  // 3. 并发执行所有组的处理任务，并等待它们全部完成。
  return Promise.all(tasks)
}

export {}
