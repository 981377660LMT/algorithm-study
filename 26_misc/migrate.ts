// !迁移

import { useState } from 'react'

/**
 * 定义了迁移过程所需的各种操作。
 * @template TSourceData 源数据的类型。
 * @template TDestinationData 目标数据的类型。
 */
interface IMigrationHandlers<TSourceData, TDestinationData> {
  /**
   * 从源位置读取数据。
   * @returns 源数据或 undefined。
   */
  sourceReader: () => TSourceData | undefined

  /**
   * 从目标位置读取数据。
   * @returns 目标数据或 undefined。
   */
  destinationReader: () => TDestinationData | undefined

  /**
   * 将数据写入目标位置。
   * @param data 要写入的数据。
   */
  destinationWriter: (data: TDestinationData) => void

  /**
   * 清理源位置的数据。这是一个异步操作，因为它可能涉及API调用。
   */
  sourceCleaner: () => Promise<void>

  /**
   * (可选) 将源数据转换为目标数据的格式。
   * 如果未提供，则假定 TSourceData 和 TDestinationData 类型兼容。
   * @param sourceData 源数据。
   * @returns 目标数据。
   */
  transformer?: (sourceData: TSourceData) => TDestinationData
}

export {}

function foo<T>() {
  function use(pageId: string): [T | undefined, (value: T | undefined) => void]
  function use(pageId: string, defaultValue: T): [T, (value: T) => void]
  function use(pageId: string, defaultValue?: T): any {}
}
