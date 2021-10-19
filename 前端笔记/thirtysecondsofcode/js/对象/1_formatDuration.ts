const formatDuration = (ms: number) => {
  if (ms < 0) ms = -ms

  const time = {
    day: Math.floor(ms / 86400000),
    hour: Math.floor(ms / 3600000) % 24,
    minute: Math.floor(ms / 60000) % 60,
    second: Math.floor(ms / 1000) % 60,
    millisecond: Math.floor(ms) % 1000,
  }

  return Object.entries(time)
    .filter(([_, val]) => val !== 0)
    .map(([unit, val]) => `${val} ${unit}${val !== 1 ? 's' : ''}`)
    .join(', ')
}
formatDuration(1001) // '1 second, 1 millisecond'
formatDuration(34325055574)
// '397 days, 6 hours, 44 minutes, 15 seconds, 574 milliseconds'
