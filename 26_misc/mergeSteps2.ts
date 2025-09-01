type StepType<T> = CreateStep<T> | UpdateStep<T> | DeleteStep<T>

class CreateStep<T> {
  readonly id: string
  readonly data: T

  constructor(id: string, data: T) {
    this.id = id
    this.data = data
  }

  merge(other: StepType<T>): StepType<T> | undefined {
    if (other.id !== this.id) {
      throw new Error('Cannot merge steps with different IDs')
    }
    if (other instanceof CreateStep) {
      throw new Error('Cannot merge two Create steps for the same ID')
    } else if (other instanceof UpdateStep) {
      return new CreateStep(this.id, { ...this.data, ...other.data })
    } else if (other instanceof DeleteStep) {
      return undefined
    }
  }
}

class UpdateStep<T> {
  readonly id: string
  readonly data: Partial<T>

  constructor(id: string, data: Partial<T>) {
    this.id = id
    this.data = data
  }

  merge(other: StepType<T>): StepType<T> | undefined {
    if (other.id !== this.id) {
      throw new Error('Cannot merge steps with different IDs')
    }
    if (other instanceof CreateStep) {
      throw new Error('Cannot merge two Create steps for the same ID')
    } else if (other instanceof UpdateStep) {
      return new UpdateStep(this.id, { ...this.data, ...other.data })
    } else if (other instanceof DeleteStep) {
      return other
    }
  }
}

class DeleteStep<T> {
  readonly id: string

  constructor(id: string) {
    this.id = id
  }

  merge(other: StepType<T>): StepType<T> | undefined {
    if (other.id !== this.id) {
      throw new Error('Cannot merge steps with different IDs')
    }
    if (other instanceof CreateStep) {
      return new UpdateStep(this.id, other.data)
    } else if (other instanceof UpdateStep) {
      throw new Error('Cannot update a deleted entity')
    } else if (other instanceof DeleteStep) {
      throw new Error('Cannot merge two Delete steps for the same ID')
    }
  }
}

/**
 * 合并操作序列.
 *
 * - `Create` + `Create` -> 报错
 * - `Create` + `Update` -> `Create`
 * - `Create` + `Delete` -> `undefined` (抵消)
 * - `Update` + `Create` -> 报错
 * - `Update` + `Update` -> `Update`
 * - `Update` + `Delete` -> `Delete`
 * - `Delete` + `Create` -> `Update` (或者报错)
 * - `Delete` + `Update` -> 报错
 * - `Delete` + `Delete` -> 报错
 */
function mergeSteps<T>(steps: StepType<T>[]): StepType<T>[] {
  const stepMap = new Map<string, StepType<T>>()
  for (const curStep of steps) {
    const { id } = curStep
    const preStep = stepMap.get(id)
    if (!preStep) {
      stepMap.set(id, curStep)
      continue
    }
    const newStep = preStep.merge(curStep)
    if (newStep === undefined) {
      stepMap.delete(id)
    } else {
      stepMap.set(id, newStep)
    }
  }
  return [...stepMap.values()]
}

export {}

if (typeof require !== 'undefined' && require.main === module) {
  // --- 测试辅助工具 ---

  /**
   * 简单的深比较函数，用于比较 data 对象。
   */
  function deepEqual(a: any, b: any): boolean {
    if (a === undefined && b === undefined) return true
    if (a === undefined || b === undefined) return false
    return JSON.stringify(a) === JSON.stringify(b)
  }

  /**
   * 断言函数，用于比较两个 Step 数组是否相等。
   */
  function assertEqual(actual: StepType<any>[], expected: StepType<any>[], testName: string) {
    let pass = true
    if (actual.length !== expected.length) {
      pass = false
    } else {
      // 对结果按 id 排序，以忽略顺序问题
      const sortFn = (a: StepType<any>, b: StepType<any>) =>
        String(a.id).localeCompare(String(b.id))
      const sortedActual = [...actual].sort(sortFn)
      const sortedExpected = [...expected].sort(sortFn)

      for (let i = 0; i < sortedActual.length; i++) {
        const act = sortedActual[i]
        const exp = sortedExpected[i]
        if (
          act.constructor.name !== exp.constructor.name ||
          act.id !== exp.id ||
          !deepEqual((act as any).data, (exp as any).data)
        ) {
          pass = false
          break
        }
      }
    }

    if (pass) {
      console.log(`✅ 通过: ${testName}`)
    } else {
      console.error(`❌ 失败: ${testName}`)
      console.error('   预期:', expected)
      console.error('   实际:', actual)
    }
  }

  /**
   * 断言函数，用于验证代码块是否按预期抛出错误。
   */
  function assertThrows(testFn: () => void, testName: string) {
    try {
      testFn()
      console.error(`❌ 失败: ${testName} - 未按预期抛出错误。`)
    } catch (e) {
      console.log(`✅ 通过: ${testName} - 已按预期抛出错误。`)
    }
  }

  // --- 测试用例 ---

  function runTests() {
    console.log('--- 开始测试 mergeSteps (面向对象版本) ---')

    // 1. 测试有效合并场景
    assertEqual(
      mergeSteps([new CreateStep('1', { name: 'A' }), new UpdateStep('1', { value: 10 })]),
      [new CreateStep('1', { name: 'A', value: 10 })],
      'Create + Update -> Create'
    )

    assertEqual(
      mergeSteps([new CreateStep('1', { name: 'A' }), new DeleteStep('1')]),
      [],
      'Create + Delete -> (抵消)'
    )

    assertEqual(
      mergeSteps([new UpdateStep('1', { name: 'A' }), new UpdateStep('1', { value: 10 })]),
      [new UpdateStep('1', { name: 'A', value: 10 })],
      'Update + Update -> Update'
    )

    assertEqual(
      mergeSteps([new UpdateStep('1', { name: 'A' }), new DeleteStep('1')]),
      [new DeleteStep('1')],
      'Update + Delete -> Delete'
    )

    assertEqual(
      mergeSteps([new DeleteStep('1'), new CreateStep('1', { name: 'B' })]),
      [new UpdateStep('1', { name: 'B' })],
      'Delete + Create -> Update'
    )

    // 2. 测试无效合并（应抛出错误）的场景
    assertThrows(() => {
      mergeSteps([new CreateStep('1', { name: 'A' }), new CreateStep('1', { name: 'B' })])
    }, 'Create + Create -> 报错')

    assertThrows(() => {
      mergeSteps([new UpdateStep('1', { name: 'A' }), new CreateStep('1', { name: 'B' })])
    }, 'Update + Create -> 报错')

    assertThrows(() => {
      mergeSteps([new DeleteStep('1'), new UpdateStep('1', { name: 'A' })])
    }, 'Delete + Update -> 报错')

    assertThrows(() => {
      mergeSteps([new DeleteStep('1'), new DeleteStep('1')])
    }, 'Delete + Delete -> 报错')

    // 3. 复杂场景和多 ID 测试
    assertEqual(
      mergeSteps([
        new CreateStep('1', { name: 'Item 1' }),
        new CreateStep('2', { name: 'Item 2' }),
        new UpdateStep('1', { status: 'active' }),
        new DeleteStep('2'),
        new CreateStep('3', { name: 'Item 3' })
      ]),
      [
        new CreateStep('1', { name: 'Item 1', status: 'active' }),
        new CreateStep('3', { name: 'Item 3' })
      ],
      '多个ID的复杂操作组合'
    )

    assertEqual(
      mergeSteps([
        new CreateStep('1', { name: 'A' }),
        new UpdateStep('1', { value: 10 }),
        new UpdateStep('1', { status: 'pending' }),
        new DeleteStep('1')
      ]),
      [],
      '长链式操作: Create -> Update -> Update -> Delete'
    )

    console.log('--- 测试结束 ---')
  }

  runTests()
}
