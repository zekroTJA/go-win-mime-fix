# go-win-mime-fix

## Introduction

If you are using go on Windows, the following function might return an unexpected result:

```go
res := mime.TypeByExtension(".js")
```

In this case, if the error occures, `res` would be something `text/plain; ...` instead of `application/javascript` like expected.

This is caused on some windows systems because, as [@mattn](https://github.com/mattn) explains in a [comment](https://github.com/golang/go/issues/32350#issuecomment-525138689) on [issue #32350](https://github.com/golang/go/issues/32350) of golang, some editors set the value of `Content Type` of the registry key `HKEY_CLASSES_ROOT\.js` on installation.

> As a result that I swimed in internet for 10 minutes, I figure out this value of the registory key seems to be possibly changed by some text editor or casual confguration.

So, this tool will detect if the error will be casued on your system, if it is becasue the unexpected value of the registry key mentioned above and will fix this by setting the keys value to the right value.

## Usage

Either you run the tool by cloning the repository and comile it yourself:
```
> git clone https://github.com/zekroTJA/go-win-mime-fix .
> go run cmd/main.go -fix
```

...or by downloading the precompiled exe from [releases](https://github.com/zekroTJA/go-win-mime-fix/releases):
```
> go-win-mime-fix.exe -fix
```

---

© 2020 Ringo Hoffmann (zekro Development)  
Covered by the MIT Licence.