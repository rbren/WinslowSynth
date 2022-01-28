This setup is working for me:
```
export LIBRARY_PATH_TO_USE="/usr/local/lib/:/Users/rbren/git/go/src/github.com/rbren/midi/portaudio/lib/.libs"
export LIBRARY_PATH="$LIBRARY_PATH_TO_USE"
export CGO_CFLAGS="-I`pwd`/portaudio/include"
export PKG_CONFIG_PATH=`pwd`/portaudio
export CGO_ENABLED=1
export CC=gcc

go run -exec "env DYLD_LIBRARY_PATH=$LIBRARY_PATH_TO_USE" main.go
```
