import React, { useCallback, useEffect, useRef } from 'react'

function useInterval(fn: Function, time: number): () => void {
  const ref = useRef(fn)
  ref.current = fn

  let cleanUpFnRef = useRef<() => void>()
  const clean = useCallback(() => {
    cleanUpFnRef.current?.()
  }, [])

  useEffect(() => {
    const timer = setInterval(() => ref.current(), time)
    cleanUpFnRef.current = () => {
      clearInterval(timer)
    }
    return clean
  }, [])

  return clean
}

// 为什么要用 useCallback 包裹返回的函数呢？
// !因为这个返回的函数可能作为参数传入别的组件，
// 这样用 useCallback 包裹就可以避免该参数的变化，
// 配合 memo 可以起到减少没必要的渲染的效果。
