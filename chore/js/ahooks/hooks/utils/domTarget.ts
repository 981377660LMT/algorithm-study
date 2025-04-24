import type { MutableRefObject } from 'react'

type TargetValue<T> = T | undefined | null

type TargetType = HTMLElement | Element | Window | Document

export type BasicTarget<T extends TargetType = Element> = (() => TargetValue<T>) | TargetValue<T> | MutableRefObject<TargetValue<T>>

export function getTargetElement<T extends TargetType>(target: BasicTarget<T>, defaultElement?: T) {
  if (!target) {
    return defaultElement
  }

  let targetElement: TargetValue<T>

  if (typeof target === 'function') {
    targetElement = target()
  } else if ('current' in target) {
    targetElement = target.current
  } else {
    targetElement = target
  }

  return targetElement
}

/**
 * 比较两组dom是否相同
 */
export default function targetAreSame(oldTarget: BasicTarget | BasicTarget[], target: BasicTarget | BasicTarget[]): boolean {
  target = Array.isArray(target) ? target : [target]
  oldTarget = Array.isArray(oldTarget) ? oldTarget : [oldTarget]

  const els = target.map(item => getTargetElement(item))
  const oldEls = target.map(item => getTargetElement(item))

  for (let i = 0; i < oldEls.length; i++) {
    if (!Object.is(oldEls[i], els[i])) return false
  }
  return true
}
