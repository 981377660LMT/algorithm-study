import {
  applyPatches,
  createDraft,
  current,
  enablePatches,
  finishDraft,
  immerable,
  nothing,
  original,
  Patch,
  produce,
  produceWithPatches
} from 'immer'

// export {
//   Draft,
//   Immer,
//   Immutable,
//   Objectish,
//   Patch,
//   PatchListener,
//   Producer,
//   StrictMode,
//   WritableDraft,
//   applyPatches,
//   castDraft,
//   castImmutable,
//   createDraft,
//   current,
//   enableMapSet,
//   enablePatches,
//   finishDraft,
//   freeze,
//   DRAFTABLE as immerable,
//   isDraft,
//   isDraftable,
//   NOTHING as nothing,
//   original,
//   produce,
//   produceWithPatches,
//   setAutoFreeze,
//   setUseStrictShallowCopy
// }

// 启用 patch listener 插件
enablePatches()

const state0 = {
  users: [
    { id: 1, name: '张三', active: false },
    { id: 2, name: '李四', active: true }
  ],
  settings: {
    darkMode: false,
    notifications: {
      email: true,
      push: false,
      sms: true
    }
  }
}

// !1.produce 是 Immer 的核心 API，接受一个初始状态和一个函数（称为 recipe），允许您对草稿进行修改。
const state1 = produce(
  state0,
  draft => {
    draft.users.push({ id: 3, name: '王五', active: true })
  },
  (patches, inversePatches) => {
    // interface Patch {
    //   op: "replace" | "remove" | "add";
    //   path: (string | number)[];
    //   value?: any;
    // }
    // type PatchListener = (patches: Patch[], inversePatches: Patch[]) => void;

    console.log('Patches:', patches)
    console.log('Inverse Patches:', inversePatches)
  }
)

// !2.与 produce 类似，但返回一个包含三个元素的数组：[下一个状态, 补丁数组, 逆补丁数组]。
const [state2, patches, inversePatches] = produceWithPatches(state0, draft => {
  draft.users.push({ id: 3, name: '王五', active: true })
})
console.log('State2:', state2)
console.log('Patches:', patches)
console.log('Inverse Patches:', inversePatches)

// !3.createDraft 和 finishDraft
// 允许手动创建和完成草稿，适用于不能立即完成状态修改的场景。
{
  const baseState = { count: 0, name: 'default' }
  // 创建草稿
  const draft = createDraft(baseState)

  // 在任何时候修改草稿
  draft.count++
  draft.name = 'test'

  // 完成草稿并生成最终状态
  const nextState = finishDraft(draft)

  console.log(nextState) // { count: 1, name: 'test' }
}
{
  const baseState = { count: 0, name: 'default' }
  let patches: Patch[] = []
  let inversePatches: Patch[] = []

  const draft = createDraft(baseState)
  draft.count += 10
  const nextState = finishDraft(draft, (p, ip) => {
    patches = p
    inversePatches = ip
  })
  console.log('Next State:', nextState) // { count: 10, name: 'default' }
  console.log('Patches:', patches) // 补丁数组
  console.log('Inverse Patches:', inversePatches) // 逆补丁数组
}

// !4.applyPatches
// 应用由 produceWithPatches 或带补丁监听器的 produce 生成的补丁。
{
  const baseState = { users: [{ name: 'John' }] }
  let patches: Patch[] = []

  const nextState = produce(
    baseState,
    draft => {
      draft.users.push({ name: 'Mike' })
    },
    p => {
      patches = p
    }
  )

  // 将原始状态回退到 nextState
  const recreatedState = applyPatches(baseState, patches)
  console.log(recreatedState === nextState) // 结构相同但不是同一对象

  // 在未来应用更多补丁
  const futureState = applyPatches(nextState, [
    { op: 'replace', path: ['users', 0, 'name'], value: 'John Doe' }
  ])
}

// !5.工具函数 current、original、isDraft、isDraftable、freeze
// 创建草稿的当前状态的快照（不冻结），用于调试或将草稿值安全地传递给外部
// 获取草稿对应的原始对象
{
  const nextState = produce({ count: 0 }, draft => {
    draft.count++
    // 查看当前修改后的草稿状态
    console.log(current(draft)) // { count: 1 }

    // 很适合在调试中使用
    // console.log(draft); // 输出 Proxy 对象，不便于查看
  })
}

// !6.特殊常量
// nothing: 特殊的哨兵值，从 recipe 返回它可以将状态设置为 undefined。
{
  const state = { a: 1, b: 2 }
  const nextState1 = produce(state, draft => {
    // 移除 b 属性
    delete draft.b
  })
  console.log(nextState1) // { a: 1 }

  const nextState2 = produce(state, draft => {
    // 使用 nothing 返回 undefined
    return nothing
  })
  console.log(nextState2) // undefined
}
// immerable: 符号常量，用于在类上标记，使 Immer 能够处理类实例。
{
  class Person {
    [immerable] = true // 标记此类为可草稿

    constructor(public name: string, public age: number) {}
  }

  const john = new Person('John', 30)
  const olderJohn = produce(john, draft => {
    draft.age += 1
  })

  console.log(john.age) // 30
  console.log(olderJohn.age) // 31
  console.log(olderJohn instanceof Person) // true
}
