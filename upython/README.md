# Installation

## on PCs

homebrew micropython ships with an old version (v2) of uasyncio even though the version of micropython is latest. Do not use brew.

Install Ubuntu 22.10

Install micropython from apt (not snap)

```
$ sudo apt install micropython
```

Do this: https://forum.micropython.org/viewtopic.php?f=15&t=10146#p56429 i.e.

```
$ git clone https://github.com/micropython/micropython.git
Cloning into 'micropython'...
...
Resolving deltas: 100% (76088/76088), done.
$ rm -rf ~/.micropython/lib/uasyncio
$ cd micropython/
$ cp -r extmod/uasyncio ~/.micropython/lib
```

Now we're good:

```
$ micropython
MicroPython v1.19.1+ds-1 on 2022-07-23; linux [GCC 12.1.0] version
Use Ctrl-D to exit, Ctrl-E for paste mode
>>> import uasyncio
>>> uasyncio.__
__class__       __getattr__     __name__        __file__
__path__        __version__
>>> uasyncio.__version__
(3, 0, 0)
>>> uasyncio.Event
<class 'Event'>
```
