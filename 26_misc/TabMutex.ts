import { v4 as uuid } from 'uuid'

/**
 * Helper function to calculate the endTime, lock acquisition time, and then
 * resolve the promise with all the lock stats
 */
const resolveWithStats = (resolve, stats) => {
  const currentTime = new Date().getTime()
  stats.acquireEnd = currentTime
  stats.acquireDuration = stats.acquireEnd - stats.acquireStart
  stats.lockStart = currentTime
  resolve(stats)
}

class TabMutex {
  clientId: string
  xPrefix: string
  yPrefix: string
  timeout: number
  localStorage: any
  lockStats: any
  constructor({
    clientId = uuid(),
    xPrefix = '_MUTEX_LOCK_X_',
    yPrefix = '_MUTEX_LOCK_Y_',
    timeout = 3000,
    localStorage = window.localStorage
  } = {}) {
    this.clientId = clientId
    this.xPrefix = xPrefix
    this.yPrefix = yPrefix
    this.timeout = timeout

    this.localStorage = localStorage
    this.resetStats()
  }

  lock(key) {
    const x = this.xPrefix + key
    const y = this.yPrefix + key
    this.resetStats()

    if (!this.lockStats.acquireStart) {
      this.lockStats.acquireStart = new Date().getTime()
    }

    return new Promise((resolve, reject) => {
      const acquireLock = _key => {
        const elapsedTime = new Date().getTime() - this.lockStats.acquireStart
        if (elapsedTime >= this.timeout) {
          return reject(new Error(`Lock could not be acquired within ${this.timeout}ms`))
        }

        this.setItem(x, this.clientId)

        // if y exists, another client is getting a lock, so retry in a bit
        let lsY = this.getItem(y)
        if (lsY) {
          this.lockStats.restartCount++
          setTimeout(() => acquireLock(_key))
          return
        }

        // ask for inner lock
        this.setItem(y, this.clientId)

        // if x was changed, another client is contending for an inner lock
        const lsX = this.getItem(x)
        if (lsX !== this.clientId) {
          this.lockStats.contentionCount++

          // Give enough time for critical section:
          setTimeout(() => {
            lsY = this.getItem(y)
            if (lsY === this.clientId) {
              resolveWithStats(resolve, this.lockStats)
            } else {
              // we lost the lock, restart the process again
              this.lockStats.restartCount++
              this.lockStats.locksLost++
              setTimeout(() => acquireLock(_key))
            }
          }, 50)
          return
        }

        // no contention:
        resolveWithStats(resolve, this.lockStats)
      }

      acquireLock(key)
    })
  }

  release(key) {
    const x = this.xPrefix + key
    const y = this.yPrefix + key
    return new Promise(resolve => {
      this.localStorage.removeItem(y)
      this.localStorage.removeItem(x)
      this.lockStats.lockEnd = new Date().getTime()
      this.lockStats.lockDuration = this.lockStats.lockEnd - this.lockStats.lockStart
      resolve(this.lockStats)
      this.resetStats()
    })
  }

  /**
   * Helper function to wrap all values in an object that includes the time (so
   * that we can expire it in the future) and json.stringify's it
   */
  setItem(key, value) {
    return this.localStorage.setItem(
      key,
      JSON.stringify({
        expiresAt: new Date().getTime() + this.timeout,
        value
      })
    )
  }

  /**
   * Helper function to parse JSON encoded values set in localStorage
   */
  getItem(key) {
    const item = this.localStorage.getItem(key)
    if (!item) return null
    let parsed
    try {
      parsed = JSON.parse(item)
      if (new Date().getTime() - parsed.expiresAt >= this.timeout) {
        this.localStorage.removeItem(key)
        return null
      }
    } catch (error) {
      //
    }

    return parsed?.value
  }

  /**
   * Helper function to reset statistics. A single FastMutex client can be used
   * to perform multiple successive lock()s so we need to reset stats each time
   */
  resetStats() {
    this.lockStats = {
      restartCount: 0,
      locksLost: 0,
      contentionCount: 0,
      acquireDuration: 0,
      acquireStart: null
    }
  }
}

export default TabMutex
