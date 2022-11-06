import asyncio
import uuid
import inspect

mailboxes = {}

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


def send(pid, msg, *args):
    mailboxes[pid].put_nowait((msg, args))


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
    send(A, "hello", 5)
    send(A, "world")
    send(B, "hello", 4)
    send(B, "hello", 8)


loop = asyncio.get_event_loop()


loop.create_task(go())
loop.run_forever()
