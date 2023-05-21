
def factorial(n):
    if n==0:
        return 1
    else:
        return n*factorial(n-1)

print(factorial(20))

def deep(n):
    print(n)
    return deep(n+1)

deep(0)
