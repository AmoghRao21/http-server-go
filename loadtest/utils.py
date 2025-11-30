import random

def rand_text(n=8):
    letters = "abcdefghijklmnopqrstuvwxyz"
    return "".join(random.choice(letters) for _ in range(n))

def rand_int(a=1, b=1000):
    return random.randint(a, b)
