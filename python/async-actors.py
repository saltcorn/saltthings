import asyncio
import uuid
import inspect
import aiohttp

mailboxes = {}
myNode = {"nodeID": str(uuid.uuid4())}

"""
todo:
-http client and server
-register process name


"""


def spawn(f):
    pid = str(uuid.uuid4())
    q = asyncio.Queue()
    mailboxes[pid] = q

    async def receive():
        value = await q.get()
        return value

    if inspect.iscoroutinefunction(f):
        asyncio.create_task(f(receive))
    if inspect.isclass(f):
        x = f()

        async def runner():
            while True:
                (msg, args) = await receive()
                op = getattr(x, msg, None)
                if callable(op):
                    op(*args)
        asyncio.create_task(runner())

    return pid


async def send(proc, msg, *args):
    print(type(proc))
    if type(proc) == str:
        mailboxes[proc].put_nowait((msg, args))
    if type(proc) == dict:
        print("isDixt")
        if proc.get("node", {}).get("nodeID","") == myNode["nodeID"]:
            print("still local")

            mailboxes[proc].put_nowait((msg, args))
            return
        httpUrl =proc.get("node", {}).get("nodeLocators",{}).get("http",False)
        if httpUrl:
            print("go http!")
            async with aiohttp.ClientSession() as session:
                await session.post(httpUrl, json=[proc, msg, *args]) 

            return

        print("unknown destination")


def run(f):
    loop = asyncio.get_event_loop()
    loop.create_task(f())
    loop.run_forever()


async def Af(receive):
    while True:
        (msg, args) = await receive()
        print("A got msg", msg, args)


class Foo:
    def hello(self, x):
        print("Foo.hello got", x)


async def go():
    A = spawn(Af)
    B = spawn(Foo)
    #send(A, "hello", 5)
    #send(A, "world")
    #send(B, "hello", 4)
    #send(B, "hello", 8)
    await send({"processName": "console", "node": {
        "nodeID": "", 
        "nodeLocators": {"http": "http://localhost:3155"}}}, 
        "print", "foobar")


run(go)
