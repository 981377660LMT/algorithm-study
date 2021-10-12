type ProduceFunc = <T>(base: T, recipe: (draft: T) => void) => any

/**
 *
 * @param base
 * @param recipe
 * 只需要实现简单对象和数组的基本使用方式
 * 需要保证未改变的数据部分不被拷贝
 * @summary
 * 对属性递归proxify
 */
const produce: ProduceFunc = (base, recipe) => {
  // your code here
}

// 太难了

function proxify(params: type) {}

const state = [
  {
    name: 'BFE',
  },
  {
    name: '.',
  },
]
const newState = produce(state, draft => {
  draft.push({ name: 'dev' })
  draft[0].name = 'bigfrontend'
  draft[1].name = '.' // set为相同值。
})

// 注意，未变化的部分并没有拷贝。
expect(newState).not.toBe(state)
expect(newState).toEqual([
  {
    name: 'bigfrontend',
  },
  {
    name: '.',
  },
  {
    name: 'dev',
  },
])
expect(newState[0]).not.toBe(state[0])
expect(newState[1]).toBe(state[1])
expect(newState[2]).not.toBe(state[2])
