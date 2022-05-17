# ODroML

Welcome to the repository of the core ODroML runtime. **Contrary to the name, it is not an ML dialect.**

## Install
```
go install github.com/odroml/odroml@latest
```

## How does FizzBuzz look?
```odroml
operator /?(a, b: int): bool {
    return a % b == 0;
}

for var i = 0; i < 100; i = i + 1 {
    var a, b: string;
    if i /? 3 {
        a = "Fizz";
    }
    if i /? 5 {
        b = "Buzz";
    }
    let v = a + b;
    match v {
        case "" {
            print(i);
        }
        default {
            print(v);
        }
    }
}
```
