caomeinaixi":
    while True:
        try:
            main()
        except (EOFError, ValueError):
            break
else: