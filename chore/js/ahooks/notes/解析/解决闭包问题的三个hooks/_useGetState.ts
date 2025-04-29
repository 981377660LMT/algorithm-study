import { Dispatch, SetStateAction, useState, useCallback } from 'react'
import { useLatest } from './useLatest'
import { useMemoizedFn } from './useMemoizedFn'

type GetStateAction<S> = () => S

/**
 * @deprecated 不要这样写，https://github.com/981377660LMT/ts/issues/811
 */
function useGetState<S>(
  initialState: S | (() => S)
): [S, Dispatch<SetStateAction<S>>, GetStateAction<S>] {
  const [state, _setState] = useState(initialState)
  const stateRef = useLatest(state)

  const setState = useMemoizedFn((action: S | ((prevState: S) => S)) => {
    const newValue = typeof action === 'function' ? (action as (prevState: S) => S)(state) : action
    stateRef.current = newValue
    _setState(newValue)
  })

  const getState = useCallback(() => stateRef.current, [])

  return [state, setState, getState]
}
