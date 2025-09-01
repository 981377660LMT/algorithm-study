// ### 操作合并逻辑分析表
//
// | 上一步 (PreStep) | 当前步 (CurStep) | 最终结果 (Result) | 说明 |
// | :--- | :--- | :--- | :--- |
// | (无) | `create` | `create` | 首次操作，记录为创建。 |
// | (无) | `update` | `update` | 首次操作，记录为更新一个已存在的实体。 |
// | (无) | `delete` | `delete` | 首次操作，记录为删除一个已存在的实体。 |
// | `create` | `update` | `create` (数据合并) | 对新创建的实体进行更新，本质上还是一个包含了所有数据的创建操作。 |
// | `create` | `delete` | **(无/移除)** | 创建后立即删除，操作相互抵消，无最终效果。 |
// | `update` | `update` | `update` (数据合并) | 连续更新，将变更数据合并。 |
// | `update` | `delete` | `delete` | 更新后删除，最终操作为删除。 |
// | `delete` | `create` | `update` 或者报错 | 删除后又创建，相当于用新数据更新了原始实体。 |
// | `create` | `create` | `create` (忽略 `create`) | 无效序列，忽略后续的 `create`。 |
// | `update` | `create` | `update` (忽略 `create`) | 无效序列，忽略后续的 `create`。 |
// | `delete` | `update` | `delete` (忽略 `update`) | 无效序列，无法更新已删除的实体，忽略 `update`。 |
// | `delete` | `delete` | `delete` (忽略 `delete`) | 连续删除，最终结果仍为删除。 |

import { isEqual } from 'lodash-es'

type StepType<T = Record<string, unknown>> =
  | { type: 'create'; id: PropertyKey; data: T }
  | { type: 'update'; id: PropertyKey; data: Partial<T> }
  | { type: 'delete'; id: PropertyKey }

/**
 * 合并 create、update、delete 操作.
 *
 * @alias compactSteps
 */
function mergeSteps<T>(
  steps: StepType<T>[],
  mergeData: <U, V>(pre: U, cur: V) => U & V
): StepType<T>[] {
  const res = new Map<PropertyKey, StepType<T>>()

  for (const curStep of steps) {
    const { id } = curStep
    const preStep = res.get(id)
    if (!preStep) {
      res.set(id, curStep)
      continue
    }

    const preType = preStep.type
    const curType = curStep.type
    switch (preType) {
      case 'create':
        if (curType === 'update') {
          const newCreateStep = {
            ...preStep,
            data: mergeData(preStep.data, curStep.data)
          }
          res.set(id, newCreateStep)
        } else if (curType === 'delete') {
          res.delete(id)
        }
        break

      case 'update':
        if (curType === 'update') {
          const newUpdateStep = {
            ...preStep,
            data: mergeData(preStep.data, curStep.data)
          }
          res.set(id, newUpdateStep)
        } else if (curType === 'delete') {
          res.set(id, curStep)
        }
        break

      case 'delete':
        if (curType === 'create') {
          // !注意
          const newUpdateStep: StepType<T> = {
            id,
            type: 'update',
            data: curStep.data
          }
          res.set(id, newUpdateStep)
        }
        break

      default:
        throw new Error(`Unsupported step type: ${preType}`)
    }
  }

  return [...res.values()]
}

export {}

if (require.main === module) {
  // 简单合并函数
  const simpleMerge = <U, V>(pre: U, cur: V): U & V => ({ ...pre, ...cur } as U & V)

  // 简单测试助手函数
  function assertEqual(actual: any, expected: any, testName: string) {
    if (isEqual(actual, expected)) {
      console.log(`✅ 通过: ${testName}`)
    } else {
      console.log(`❌ 失败: ${testName}`)
      console.log('   预期:', expected)
      console.log('   实际:', actual)
    }
  }

  // 测试用例
  function runTests() {
    console.log('开始测试 mergeSteps 函数...')

    // 1. 基本操作测试 - 单个 create/update/delete
    assertEqual(
      mergeSteps([{ type: 'create', id: 1, data: { name: 'Item 1' } }], simpleMerge),
      [{ type: 'create', id: 1, data: { name: 'Item 1' } }],
      '单个create操作'
    )

    assertEqual(
      mergeSteps([{ type: 'update', id: 1, data: { name: 'Updated' } }], simpleMerge),
      [{ type: 'update', id: 1, data: { name: 'Updated' } }],
      '单个update操作'
    )

    assertEqual(
      mergeSteps([{ type: 'delete', id: 1 }], simpleMerge),
      [{ type: 'delete', id: 1 }],
      '单个delete操作'
    )

    // 2. 组合操作测试
    // create + update = create(合并数据)
    assertEqual(
      mergeSteps(
        [
          { type: 'create', id: 1, data: { name: 'Item 1', count: 10 } },
          { type: 'update', id: 1, data: { count: 20 } }
        ],
        simpleMerge
      ),
      [{ type: 'create', id: 1, data: { name: 'Item 1', count: 20 } }],
      'create + update = create(合并数据)'
    )

    // create + delete = 无操作
    assertEqual(
      mergeSteps(
        [
          { type: 'create', id: 1, data: { name: 'Item 1' } },
          { type: 'delete', id: 1 }
        ],
        simpleMerge
      ),
      [],
      'create + delete = 无操作'
    )

    // update + update = update(合并数据)
    assertEqual(
      mergeSteps(
        [
          { type: 'update', id: 1, data: { name: 'Updated' } },
          { type: 'update', id: 1, data: { count: 30 } }
        ],
        simpleMerge
      ),
      [{ type: 'update', id: 1, data: { name: 'Updated', count: 30 } }],
      'update + update = update(合并数据)'
    )

    // update + delete = delete
    assertEqual(
      mergeSteps(
        [
          { type: 'update', id: 1, data: { name: 'Updated' } },
          { type: 'delete', id: 1 }
        ],
        simpleMerge
      ),
      [{ type: 'delete', id: 1 }],
      'update + delete = delete'
    )

    // delete + create = update
    assertEqual(
      mergeSteps(
        [
          { type: 'delete', id: 1 },
          { type: 'create', id: 1, data: { name: 'New Item' } }
        ],
        simpleMerge
      ),
      [{ type: 'update', id: 1, data: { name: 'New Item' } }],
      'delete + create = update'
    )

    // 3. 忽略无效序列测试
    // create + create = create(忽略后续create)
    assertEqual(
      mergeSteps(
        [
          { type: 'create', id: 1, data: { name: 'Item 1' } },
          { type: 'create', id: 1, data: { name: 'Item 1 New' } }
        ],
        simpleMerge
      ),
      [{ type: 'create', id: 1, data: { name: 'Item 1' } }],
      'create + create = create(忽略后续create)'
    )

    // update + create = update(忽略后续create)
    assertEqual(
      mergeSteps(
        [
          { type: 'update', id: 1, data: { name: 'Updated' } },
          { type: 'create', id: 1, data: { name: 'New Item' } }
        ],
        simpleMerge
      ),
      [{ type: 'update', id: 1, data: { name: 'Updated' } }],
      'update + create = update(忽略后续create)'
    )

    // delete + update = delete(忽略update)
    assertEqual(
      mergeSteps(
        [
          { type: 'delete', id: 1 },
          { type: 'update', id: 1, data: { name: 'Updated' } }
        ],
        simpleMerge
      ),
      [{ type: 'delete', id: 1 }],
      'delete + update = delete(忽略update)'
    )

    // delete + delete = delete(忽略后续delete)
    assertEqual(
      mergeSteps(
        [
          { type: 'delete', id: 1 },
          { type: 'delete', id: 1 }
        ],
        simpleMerge
      ),
      [{ type: 'delete', id: 1 }],
      'delete + delete = delete(忽略后续delete)'
    )

    // 4. 复杂场景测试
    // 多个ID的操作组合
    assertEqual(
      mergeSteps(
        [
          { type: 'create', id: 1, data: { name: 'Item 1' } },
          { type: 'create', id: 2, data: { name: 'Item 2' } },
          { type: 'update', id: 1, data: { status: 'active' } },
          { type: 'delete', id: 2 },
          { type: 'create', id: 3, data: { name: 'Item 3' } }
        ],
        simpleMerge
      ),
      [
        { type: 'create', id: 1, data: { name: 'Item 1', status: 'active' } },
        { type: 'create', id: 3, data: { name: 'Item 3' } }
      ],
      '多个ID的操作组合'
    )

    // 链式操作: create -> update -> delete -> create
    assertEqual(
      mergeSteps(
        [
          { type: 'create', id: 1, data: { name: 'Original' } },
          { type: 'update', id: 1, data: { status: 'active' } },
          { type: 'delete', id: 1 },
          { type: 'create', id: 1, data: { name: 'Recreated' } }
        ],
        simpleMerge
      ),
      [{ type: 'create', id: 1, data: { name: 'Recreated' } }],
      '链式操作: create -> update -> delete -> create'
    )

    // update -> delete -> create
    assertEqual(
      mergeSteps(
        [
          { type: 'update', id: 1, data: { name: 'Updated' } },
          { type: 'delete', id: 1 },
          { type: 'create', id: 1, data: { name: 'New Item' } }
        ],
        simpleMerge
      ),
      [{ type: 'update', id: 1, data: { name: 'New Item' } }],
      '链式操作: update -> delete -> create'
    )

    console.log('测试完成')
  }

  // 运行测试
  runTests()
}
