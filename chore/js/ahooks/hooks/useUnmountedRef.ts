/**
 * copy 自 ahook @link https://ahooks.js.org/zh-CN/hooks/use-unmounted-ref
 */
import { useEffect, useRef } from 'react'

const useUnmountedRef = () => {
  const unmountedRef = useRef(false)
  useEffect(() => {
    unmountedRef.current = false
    return () => {
      unmountedRef.current = true
    }
  }, [])
  return unmountedRef
}

export default useUnmountedRef
