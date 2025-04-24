/* eslint-disable @typescript-eslint/ban-types */

/**
 * https://github.com/alibaba/hooks/tree/master/packages/hooks/src/createUseStorageState
 */

import { useState, useEffect, useRef } from 'react'
import { useMemoizedFn } from '../useMemoizedFn'

const isUndefined = (val: unknown): val is undefined => typeof val === 'undefined'
const isFunction = (val: unknown): val is Function => typeof val === 'function'

export interface IFuncUpdater<T> {
  (previousState?: T): T
}

export interface Options<T> {
  serializer?: (value: T) => string
  deserializer?: (value: string) => T
  defaultValue?: T | IFuncUpdater<T>
}

export function createUseStorageState(getStorage: () => Storage | undefined) {
  function useStorageState<T>(key: string, options?: Options<T>) {
    let storage: Storage | undefined

    // https://github.com/alibaba/hooks/issues/800
    try {
      storage = getStorage()
    } catch (err) {
      // eslint-disable-next-line no-console
      console.error(err)
    }

    const serializer = (value: T) => {
      if (options?.serializer) {
        return options?.serializer(value)
      }
      return JSON.stringify(value)
    }

    const deserializer = (value: string) => {
      if (options?.deserializer) {
        return options?.deserializer(value)
      }
      return JSON.parse(value)
    }

    function getStoredValue() {
      try {
        const raw = storage?.getItem(key)
        if (raw) {
          return deserializer(raw)
        }
      } catch (e) {
        // eslint-disable-next-line no-console
        console.error(e)
      }
      if (isFunction(options?.defaultValue)) {
        return options?.defaultValue && options.defaultValue()
      }
      return options?.defaultValue
    }

    const [state, setState] = useState<T>(() => getStoredValue())

    const isMountedRef = useRef(false)
    useEffect(() => {
      if (isMountedRef.current) {
        setState(getStoredValue())
      } else {
        isMountedRef.current = true
      }
    }, [key])

    const updateState = (value: T | IFuncUpdater<T>) => {
      const currentState = isFunction(value) ? value(state) : value
      setState(currentState)

      if (isUndefined(currentState)) {
        storage?.removeItem(key)
      } else {
        try {
          storage?.setItem(key, serializer(currentState))
        } catch (e) {
          // eslint-disable-next-line no-console
          console.error(e)
        }
      }
    }

    return [state, useMemoizedFn(updateState)] as const
  }

  return useStorageState
}
