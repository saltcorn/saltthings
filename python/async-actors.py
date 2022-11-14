import asyncio
import uuid
import inspect
import aiohttp
from aiohttp import web

mailboxes = {}
myNode = {
    "nodeID": str(uuid.uuid4()),
    "nodeLocators": {}
}

"""
todo:
-register process name
"""


async def createNode(loop, options={}):
    global myNode
    
    if options.get("createHttpServer"):
        host = options['createHttpServer']['host']
        port = options['createHttpServer']['port']
        myNode['nodeLocators']['http'] = f"http://{host}:{port}"
        

        async def handler(request):
            pid, msg, *args = await request.json()
            #print(pid,msg,args)
            await send(pid,msg,*args)
            return web.Response(text="OK")
        app = web.Application(loop=loop)
        app.add_routes([web.post('/', handler)])
        srv = await loop.create_server(app.make_handler(), '0.0.0.0', port)
        return srv # TODO remove?


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

    return {"pid": pid, "node": myNode}


async def send(proc, msg, *args):
    if type(proc) == str:
        mailboxes[proc].put_nowait((msg, args))
    if type(proc) == dict:
        if proc.get("node", {}).get("nodeID", "") == myNode["nodeID"]:
            mailboxes[proc['pid']].put_nowait((msg, args))
            return
        httpUrl = proc.get("node", {}).get(
            "nodeLocators", {}).get("http", False)
        if httpUrl:
            async with aiohttp.ClientSession() as session:
                await session.post(httpUrl, json=[proc, msg, *args])

            return

        print("unknown destination")


def run(f, createNodeOptions={}):
    loop = asyncio.get_event_loop()
    async def g():
        await createNode(loop, createNodeOptions)
        await f()
    loop.create_task(g())
    loop.run_forever()


async def Af(receive):
    while True:
        (msg, args) = await receive()
        print("A got msg", msg, args)


class Foo:
    def hello(self, x):
        print("Foo.hello got", x)
class TickTock:
    def tock(self, tm):
        print("tock got", tm)


async def go():
    #A = spawn(Af)
    #B = spawn(Foo)
    tickTock = spawn(TickTock)
    #await send(A, "hello", 5)
    #await send(A, "world")
    #await send(B, "hello", 4)
    #await send(B, "hello", 8)
    await send(
        {
            "processName": "clock",
            "node": {
                "nodeID": "",
                "nodeLocators": {"http": "http://localhost:3135"}}
        },
        "get_time",
        tickTock,
        "tock",
        )


run(go, {
    "createHttpServer": {"port": 3001, "host": "0.0.0.0"}
})
