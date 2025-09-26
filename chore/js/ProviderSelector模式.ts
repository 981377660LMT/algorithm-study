import { useMemo } from 'react'

// ===================================================================
// 1. 定义抽象类型
// ===================================================================

/**
 * 原始的、未经处理的数据类型。
 * 例如：一个巨大的对象数组。
 */
type TRawData = any

/**
 * 经过昂贵计算后得到的数据结构类型。
 * 例如：一个用于快速查找的 Map 或一个分组后的对象。
 */
type TProcessedData = any

/**
 * 用于从已处理数据中选择特定部分的参数类型。
 * 例如：一个 ID 字符串或一个索引号。
 */
type TSelectorParam = any

/**
 * 从已处理数据中选择出的最终数据片段的类型。
 */
type TSelectedData = any

// ===================================================================
// 2. 定义核心函数签名
// ===================================================================

/**
 * 占位符：一个昂贵的计算函数。
 * @param rawData - 原始数据。
 * @returns 经过处理的数据。
 */
declare function expensiveComputation(rawData: TRawData): TProcessedData

/**
 * 占位符：一个轻量级的选择函数。
 * @param processedData - 已经过计算和缓存的数据。
 * @param selectorParam - 用于查找的参数。
 * @returns 所需的数据片段。
 */
declare function selectData(
  processedData: TProcessedData,
  selectorParam: TSelectorParam
): TSelectedData

/**
 * 占位符：一个用于获取原始数据的 Hook。
 * 在真实应用中，这可能来自 Zustand、Redux、React Context 或 SWR。
 */
declare function useRawDataSource(): TRawData

// ===================================================================
// 3. 实现抽象模式
// ===================================================================

/**
 * 模式第一部分：数据提供者 Hook (Provider Hook)
 *
 * 职责：
 * 1. 获取原始数据。
 * 2. 执行一次昂贵的计算。
 * 3. 使用 `useMemo` 缓存计算结果，直到原始数据发生变化。
 *
 * 这个 Hook 是整个模式的核心，它确保了昂贵的计算被集中处理和缓存。
 */
function useProcessedSharedData(): TProcessedData {
  const rawData = useRawDataSource()

  const processedData = useMemo(() => {
    // 昂贵的计算只在这里执行，并且其结果会被缓存。
    return expensiveComputation(rawData)
  }, [rawData]) // 仅当原始数据变化时才重新计算

  return processedData
}

/**
 * 模式第二部分：数据消费者 Hook (Selector Hook)
 *
 * 职责：
 * 1. 调用 Provider Hook (`useProcessedSharedData`) 来获取已缓存的、处理好的数据。
 * 2. 执行一个轻量级的选择操作，以从共享数据中提取特定组件所需的部分。
 * 3. 返回最终所需的数据片段。
 *
 * 这个 Hook 可以被任意数量的组件安全、高效地调用，而不会触发重复的昂贵计算。
 */
function useSelectedData(selectorParam: TSelectorParam): TSelectedData {
  // 获取已缓存的、全局处理好的数据。此步骤非常快。
  const processedData = useProcessedSharedData()

  // 从缓存数据中选择所需的部分。此步骤也非常快。
  const selectedData = useMemo(() => {
    return selectData(processedData, selectorParam)
  }, [processedData, selectorParam]) // 仅当共享数据或选择参数变化时才重新选择

  return selectedData
}

export {}
