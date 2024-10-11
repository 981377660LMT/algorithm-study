import { SetStateAction, useCallback, useEffect, useRef, useState } from 'react'

/**
 * 屏蔽受控组件和非受控组件的差异.
 *
 * @see
 * - {@link https://github.com/alibaba/hooks/blob/master/packages/hooks/src/useControllableValue/index.ts}
 * - {@link https://juejin.cn/post/7207810396420669477}
 *
 * @summary
 * - 非受控组件：如果 props 中没有 value，则组件内部自己管理 state
 * - 受控组件：如果 props 有 value 字段，则由父级接管控制 state
 * - 无 value，有 onChange 的组件：只要 props 中有 onChange 字段，则在 state 变化时，就会触发 onChange 函数
 *
 * @alias useControllableValue
 */
export function useMergeState<T>(
  defaultStateValue: T,
  props?: {
    defaultValue?: T
    value?: T
    onChange?: (value: T) => void
  }
): [T, React.Dispatch<React.SetStateAction<T>>] {
  const { defaultValue, value: propsValue, onChange } = props || {}

  const isFirstRender = useRef(true)

  const [stateValue, setStateValue] = useState<T>(() => {
    if (propsValue !== undefined) {
      return propsValue!
    } else if (defaultValue !== undefined) {
      return defaultValue!
    } else {
      return defaultStateValue
    }
  })

  useEffect(() => {
    if (propsValue === undefined && !isFirstRender.current) {
      setStateValue(propsValue!)
    }
    isFirstRender.current = false
  }, [propsValue])

  const mergedValue = propsValue === undefined ? stateValue : propsValue

  const setState = useCallback(
    (value: SetStateAction<T>) => {
      let res = isFunction(value) ? value(stateValue) : value
      if (propsValue === undefined) {
        setStateValue(res)
      }
      onChange?.(res)
    },
    [stateValue]
  )

  return [mergedValue, setState]
}

function isFunction(value: unknown): value is Function {
  return typeof value === 'function'
}
