import uasyncio
import random
import time
from primitives import Queue
import urequests
import ujson
from uWeb.uWeb_uasyncio import uWeb_uasyncio as uWeb
from uWeb.uWeb_uasyncio import loadJSON

random.seed(time.ticks_ms())

def randChar():
    return 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_'[random.getrandbits(6)]

def randStr(N=16):
	return ''.join([randChar() for i in range(N)])

async def a_coroutine():
    pass

class aClass:
    pass


type_coroutine = type(a_coroutine)
type_class = type(aClass)

mailboxes = {}
myNode = {
    "nodeID": randStr(),
    "nodeLocators": {}
}

async def createNode(loop, options={}):
    global myNode
    
    if options.get("createHttpServer"):
        host = options['createHttpServer']['host']
        port = options['createHttpServer']['port']
        myNode['nodeLocators']['http'] = f"http://{host}:{port}"
        server = uWeb("0.0.0.0", port)
        def post(): #print JSON body from client
            print('Payload: ', loadJSON(server.request_body))
            await uasyncio.sleep(0)
        server.routes(({
            (uWeb.POST, "/"): post,
        }))
        loop.create_task(uasyncio.start_server(server.router, server.address, server.port)) 
      

def spawn(f):
    pid =randStr()
    q = Queue()
    mailboxes[pid] = q

    async def receive():
        value = await q.get()
        return value
    if type(f)==type_coroutine:
        uasyncio.create_task(f(receive))
    if type(f)==type_class:
        x = f()

        async def runner():
            while True:
                (msg, args) = await receive()
                op = getattr(x, msg, None)
                if callable(op):
                    op(*args)
        uasyncio.create_task(runner())

    return {"pid": pid, "node": myNode}

def send(proc, msg, *args):
    if type(proc) == str:
        mailboxes[proc].put_nowait((msg, args))
    if type(proc) == dict:
        if proc.get("node", {}).get("nodeID", "") == myNode["nodeID"]:
            mailboxes[proc['pid']].put_nowait((msg, args))
            return
        httpUrl = proc.get("node", {}).get(
            "nodeLocators", {}).get("http", False)
        if httpUrl:       
            post_data = ujson.dumps([proc, msg]+list(args))
            print(post_data)
            urequests.post(httpUrl, headers = {'content-type': 'application/json'}, data = post_data)            
            return

        print("unknown destination")

def run(f, createNodeOptions={}):
    loop = uasyncio.get_event_loop()
    async def g():
        await createNode(loop, createNodeOptions)
        await f()
    loop.create_task(g())
    loop.run_forever()

async def Af(receive):
    while True:
        print("Af waiting")
        (msg, args) = await receive()
        print("A got msg", msg, args)

class Foo:
    def hello(self, x):
        print("Foo.hello got", x)
class TickTock:
    def tock(self, tm):
        print("tock got", tm)

async def go():
    A = spawn(Af)
    
    B = spawn(Foo)
    send(A, "hello", 5)
    send(B, "hello", 4)
    send({
            "processName": "console",
            "node": {
                "nodeID": "",
                "nodeLocators": {"http": "http://localhost:3155"}}
        },
        "print",
        "tock")
    """
    await uasyncio.sleep(2)
    tickTock = spawn(TickTock)
    #await send(A, "hello", 5)
    #await send(A, "world")
    #await send(B, "hello", 4)
    #await send(B, "hello", 8)
    send(
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

    print("done")

run(go, {
    "createHttpServer": {"port": 3001, "host": "0.0.0.0"}
})