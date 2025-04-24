import { useEffect, useRef } from 'react'
import type { useLayoutEffect } from 'react'

type EffectHookType = typeof useEffect | typeof useLayoutEffect

const createUpdateEffect: (hook: EffectHookType) => EffectHookType = hook => (effect, deps) => {
  const isMounted = useRef(false)

  // for react-refresh
  hook(() => {
    return () => {
      isMounted.current = false
    }
  }, [])

  hook(() => {
    if (!isMounted.current) {
      isMounted.current = true
    } else {
      return effect()
    }
  }, deps)
}

/**
 * useUpdateEffect 用法等同于 useEffect，但是会忽略首次执行，只在依赖更新时执行。
 * https://ahooks.js.org/zh-CN/hooks/use-update-effect
 */
export const useUpdateEffect = createUpdateEffect(useEffect)
