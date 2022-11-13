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
myNode = {
    "nodeID": randStr(),
    "nodeLocators": {}
}

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

async def send(proc, msg, *args):
    print("send", type(proc))

    if type(proc) == str:
        mailboxes[proc].put_nowait((msg, args))
    if type(proc) == dict:
        if proc.get("node", {}).get("nodeID", "") == myNode["nodeID"]:
            print("delivering", msg, args)
            print("to", mailboxes[proc['pid']])
            mailboxes[proc['pid']].put_nowait((msg, args))
            return

def run(f, createNodeOptions={}):
    loop = uasyncio.get_event_loop()
    """async def g():
        await createNode(loop, createNodeOptions)
        await f()"""
    loop.create_task(f())
    loop.run_forever()

async def Af(receive):
    while True:
        print("Af waiting")
        (msg, args) = await receive()
        print("A got msg", msg, args)

async def go():
    A = spawn(Af)
    
    #B = spawn(Foo)
    await send(A, "hello", 5)
    #await uasyncio.sleep(3)

    print("done")

run(go, {
    "createHttpServer": {"port": 3001, "host": "0.0.0.0"}
})