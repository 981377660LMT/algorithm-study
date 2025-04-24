import { useCallback, useState } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import useUnmountedRef from './useUnmountedRef'

function useSafeState<S = undefined>(initialState?: S | (() => S)): [S, Dispatch<SetStateAction<S>>] {
  const unmountedRef = useUnmountedRef()
  const [state, setState] = useState(initialState)
  const setCurrentState = useCallback(currentState => {
    /** if component is unmounted, stop update */
    if (unmountedRef.current) return
    setState(currentState)
  }, [])

  return [state, setCurrentState]
}

export default useSafeState
