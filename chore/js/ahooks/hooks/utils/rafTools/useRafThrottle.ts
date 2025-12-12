import { useRef, useCallback, useEffect } from 'react';

export function useRafThrottle<T extends (...args: any[]) => void>(callback: T) {
  const rafRef = useRef<number>();
  const callbackRef = useRef(callback);
  useEffect(() => {
    callbackRef.current = callback;
  });

  const throttled = useCallback((...args: Parameters<T>) => {
    if (rafRef.current) {
      return;
    }
    rafRef.current = requestAnimationFrame(() => {
      callbackRef.current(...args);
      rafRef.current = undefined;
    });
  }, []) as T;

  useEffect(
    () => () => {
      if (rafRef.current) {
        cancelAnimationFrame(rafRef.current);
      }
    },
    []
  );

  return throttled;
}
