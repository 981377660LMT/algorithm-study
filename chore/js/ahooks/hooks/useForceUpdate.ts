import { useCallback, useState } from 'react'

export function useForceUpdate(): () => void {
  const [_, setState] = useState<boolean>(false)
  const forceUpdate = useCallback(() => {
    setState(s => !s)
  }, [])
  return forceUpdate
}
