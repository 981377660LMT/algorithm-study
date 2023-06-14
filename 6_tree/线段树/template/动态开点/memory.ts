/* eslint-disable @typescript-eslint/no-explicit-any */

import v8 from 'v8'
import vm from 'vm'

function printMemoryUsage(): void {
  const used = process.memoryUsage() // 字节为单位
  Object.entries(used).forEach(([key, value]) => {
    console.log(`${key}:${format(value)}`)
  })

  function format(bytes: number): string {
    return `${(bytes / 1024 / 1024).toFixed(2)}MB`
  }
}

function getSizeOf(createObject: (...args: any[]) => any): string {
  gc()
  const preMemory = process.memoryUsage().heapUsed
  createObject()
  const afterMemory = process.memoryUsage().heapUsed
  return `${((afterMemory - preMemory) / 1024 / 1024).toFixed(2)}MB`
}

function gc(): void {
  v8.setFlagsFromString('--expose-gc')
  vm.runInNewContext('gc')
}

export { printMemoryUsage, getSizeOf }

if (require.main === module) {
  console.log(
    getSizeOf(() => {
      const arrayNode = new Array(1e6)
      for (let i = 0; i < arrayNode.length; i++) {
        arrayNode[i] = [1, 0, 2, 1]
      }
    })
  )

  console.log(
    getSizeOf(() => {
      const objectNode = new Array(1e6)
      for (let i = 0; i < objectNode.length; i++) {
        objectNode[i] = { l: 1, r: 0, v: 2, d: 1 }
      }
    })
  )
}
