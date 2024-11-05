def Process(name):
    for _ in range(5):
        sys_write(name)

def main():
    sys_spawn(Process, 'A')
    sys_spawn(Process, 'B')
