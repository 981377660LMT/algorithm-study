type Draft<T> = { [K in keyof T]: T[K] }

const DRAFT_STATE = Symbol('draftState')

interface DraftState {
  original: any
  copy: any
  parent?: DraftState
  modified: boolean
  drafts: Map<PropertyKey, any>
}

function isPlainObject(value: any): boolean {
  if (value === null || typeof value !== 'object') return false
  const proto = Object.getPrototypeOf(value)
  return proto === Object.prototype || proto === null
}

/**
 * 对传入的对象或数组生成一个代理（Proxy），延迟克隆，拦截 get、set、delete 操作，并记录修改状态.
 * 嵌套对象在第一次访问时才会递归创建子草稿（draft).
 */
function createDraft<T extends object>(base: T, parent?: DraftState): Draft<T> {
  const state: DraftState = {
    original: base,
    copy: Array.isArray(base) ? [...(base as any)] : { ...base },
    parent,
    modified: false,
    drafts: new Map()
  }

  const proxy = new Proxy(base, {
    get(target, prop, receiver) {
      if (prop === DRAFT_STATE) return state
      const source = state.modified ? state.copy : state.original
      const value = Reflect.get(source, prop, receiver)
      if (isPlainObject(value)) {
        if (!state.drafts.has(prop)) state.drafts.set(prop, createDraft(value, state))
        return state.drafts.get(prop)
      }
      return value
    },
    set(target, prop, value, receiver) {
      if (!state.modified) {
        state.modified = true
        state.copy = Array.isArray(state.original) ? [...state.original] : { ...state.original }
      }
      return Reflect.set(state.copy, prop, value, receiver)
    },
    deleteProperty(target, prop) {
      if (!state.modified) {
        state.modified = true
        state.copy = Array.isArray(state.original) ? [...state.original] : { ...state.original }
      }
      return Reflect.deleteProperty(state.copy, prop)
    }
  })

  return proxy as Draft<T>
}

/**
 * 遍历所有草稿，如果没有修改直接返回原始对象；如果修改则返回新的对象，并对嵌套草稿递归调用.
 */
function finalize<T>(draft: T): T {
  const state: DraftState = (draft as any)[DRAFT_STATE]
  if (!state) return draft
  if (!state.modified) return state.original
  const result = state.copy
  Object.keys(result).forEach(key => {
    const value = (result as any)[key]
    ;(result as any)[key] = finalize(value)
  })
  return result
}

function produce<T extends object>(base: T, recipe: (draft: Draft<T>) => void): T {
  const draft = createDraft(base)
  recipe(draft)
  return finalize(draft)
}

export {}
