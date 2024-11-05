#!/usr/bin/env python3

import sys
import random
from pathlib import Path


class OS:
    """
    A minimal executable operating system model. Processes
    are state machines (Python generators) that can be paused
    or continued with local states being saved.
    """

    """
    We implement four system calls:

    - read: read a random bit value.
    - write: write a string to the buffer.
    - spawn: create a new state machine (process).
    """
    SYSCALLS = ["read", "write", "spawn"]

    class Process:
        """
        A "freezed" state machine. The state (local variables,
        program counters, etc.) are stored in the generator
        object.
        """

        def __init__(self, func, *args):
            # func should be a generator function. Calling
            # func(*args) returns a generator object.
            self._func = func(*args)

            # This return value is set by the OS's main loop.
            self.retval = None

        def step(self):
            """
            Resume the process with OS-written return value,
            until the next system call is issued.
            """
            syscall, args, *_ = self._func.send(self.retval)
            self.retval = None
            return syscall, args

    def __init__(self, src):
        # This is a hack: we directly execute the source
        # in the current Python runtime--and main is thus
        # available for calling.
        exec(src, globals())
        self.procs = [OS.Process(main)]
        self.buffer = ""

    def run(self):
        # Real operating systems waste all CPU cycles
        # (efficiently, by putting the CPU into sleep) when
        # there is no running process at the moment. Our model
        # terminates if there is nothing to run.
        while self.procs:

            # There is also a pointer to the "current" process
            # in today's operating systems.
            current = random.choice(self.procs)

            try:
                # Operating systems handle interrupt and system
                # calls, and "assign" CPU to a process.
                match current.step():
                    case "read", _:
                        current.retval = random.choice([0, 1])
                    case "write", s:
                        self.buffer += s
                    case "spawn", (fn, *args):
                        self.procs += [OS.Process(fn, *args)]
                    case _:
                        assert 0

            except StopIteration:
                # The generator object terminates.
                self.procs.remove(current)

        return self.buffer


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} file")
        exit(1)

    src = Path(sys.argv[1]).read_text()

    # Hack: patch sys_read(...) -> yield "sys_read", (...)
    for syscall in OS.SYSCALLS:
        src = src.replace(f"sys_{syscall}", f'yield "{syscall}", ')

    stdout = OS(src).run()
    print(stdout)
