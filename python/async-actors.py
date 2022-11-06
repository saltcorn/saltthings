import asyncio
import uuid

mailboxes = {}

"""
todo:
- msg arguments
- spawn class
"""


def spawn(f):
    pid = str(uuid.uuid4())
    q = asyncio.Queue()
    mailboxes[pid] = q

    async def receive():
        value = await q.get()
        return value
    asyncio.create_task(f(receive))
    return pid


def send(pid, msg, *args):
    mailboxes[pid].put_nowait((msg, args))


async def Af(receive):
    while True:
        (msg, args) = await receive()
        print("A got msg", msg, args)


async def go():
    A = spawn(Af)
    send(A, "hello", 5)
    send(A, "world")

loop = asyncio.get_event_loop()


loop.create_task(go())
loop.run_forever()
