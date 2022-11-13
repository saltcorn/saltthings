import uasyncio
import random
import time
from primitives import Queue
#import inspect

random.seed(time.ticks_ms())


def randChar():
  return 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_'[random.getrandbits(6)]

def randStr(N=16):
	return ''.join([randChar() for i in range(N)])

mailboxes = {}
myNode= {}

def spawn(f):
    pid =randStr()
    q = Queue()
    mailboxes[pid] = q

    async def receive():
        value = await q.get()
        return value

    #if inspect.iscoroutinefunction(f):
    uasyncio.create_task(f(receive))
    """if inspect.isclass(f):
        x = f()

        async def runner():
            while True:
                (msg, args) = await receive()
                op = getattr(x, msg, None)
                if callable(op):
                    op(*args)
        asyncio.create_task(runner())"""

    return {"pid": pid, "node": myNode}

def run(f, createNodeOptions={}):
    loop = uasyncio.get_event_loop()
    """async def g():
        await createNode(loop, createNodeOptions)
        await f()"""
    loop.create_task(f())
    loop.run_forever()

async def Af(receive):
    while True:
        (msg, args) = await receive()
        print("A got msg", msg, args)

async def go():
    A = spawn(Af)
    #B = spawn(Foo)
    #send(A, "hello", 5)
   

run(go, {
    "createHttpServer": {"port": 3001, "host": "0.0.0.0"}
})