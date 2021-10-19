const cb = () => console.log('Animation frame fired')
const recorder = recordAnimationFrames(cb)
// logs 'Animation frame fired' on each animation frame
recorder.stop() // stops logging
recorder.start() // starts again
const recorder2 = recordAnimationFrames(cb, false)
// `start` needs to be explicitly called to begin recording frames

function recordAnimationFrames(cb: () => void, autoStart = true) {
  return new AnimationFramesRecorder(cb, autoStart)
}

interface IAnimationFramesRecorder {
  stop(): void
  start(): void
}

class AnimationFramesRecorder implements IAnimationFramesRecorder {
  constructor(
    private callback: () => void,
    private autoStart = true,
    private running = false,
    private handle = Infinity
  ) {
    autoStart && this.start()
  }

  stop(): void {
    if (!this.running) return
    this.running = false
    cancelAnimationFrame(this.handle)
  }

  start(): void {
    if (this.running) return
    this.running = true
    this.run()
  }

  private run(): void {
    this.handle = requestAnimationFrame(() => {
      this.callback()
      this.running && this.run()
    })
  }
}
