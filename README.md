# grape 🍇
grape: a multi-threaded grep implementation written in Go

## clone grape

```
$ git clone https://github.com/carlomunguia/grape
$ cd grape
```

## run grape

```
$ cd grape (the grape repo can be placed anywhere other files sit)
```
use the go run command within grape, followed by a search term. example:
```
$ go run ./grape hello .
```
the above will pull all instances of hello within the file structure (like classic grep! :) )
