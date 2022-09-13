# bush

A `tree`-like directory walker written in Go

## Install + Run

Simply clone this repo, install dependencies and build the Go code:

```
git clone https://github.com/h5law/bush
cd bush
go mod tidy
go build -o bush
```

Then move the `bush` executable into the `$PATH` location

```
sudo mv bush /usr/local/bin
```

## TODO

- Figure out how to remove the excess lines in the indent of files in last
directories where connecting lines are not needed
