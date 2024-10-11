import React from 'react'
import { useMergeState } from './useMergedState'

interface CalendarProps {
  value?: Date
  defaultValue?: Date
  onChange?: (date: Date) => void
}

function Calendar(props: CalendarProps) {
  const { value: propsValue, defaultValue, onChange } = props

  const [mergedValue, setValue] = useMergeState(new Date(), {
    value: propsValue,
    defaultValue,
    onChange
  })

  return (
    <div>
      {mergedValue?.toLocaleDateString()}
      <div
        onClick={() => {
          setValue(new Date('2024-5-1'))
        }}
      >
        2023-5-1
      </div>
      <div
        onClick={() => {
          setValue(new Date('2024-5-2'))
        }}
      >
        2023-5-2
      </div>
      <div
        onClick={() => {
          setValue(new Date('2024-5-3'))
        }}
      >
        2023-5-3
      </div>
    </div>
  )
}
