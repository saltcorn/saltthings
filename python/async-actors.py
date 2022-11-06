import asyncio
import uuid

mailboxes = {}


def spawn(f):
    pid = str(uuid.uuid4())
    q = asyncio.Queue()
    mailboxes[pid] = q

    async def receive():
        value = await q.get()
        return value
    asyncio.create_task(f(receive))
    return pid


def send(pid, msg):
    mailboxes[pid].put_nowait(msg)


async def loop(pid, f):
    pass


async def Af(receive):
    msg = await receive()
    print("A got msg", msg)
    msg = await receive()
    print("A got msg", msg)


async def go():
    A = spawn(Af)
    send(A, "hello")
    send(A, "world")

loop = asyncio.get_event_loop()


loop.create_task(go())
loop.run_forever()
