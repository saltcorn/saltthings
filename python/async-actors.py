import asyncio
import uuid

mailboxes = {}


async def spawn(f):
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


async def Af(receive):
    msg = await receive()
    print("A got msg", msg)


async def go():
    A = await spawn(Af)
    send(A, "hello")

loop = asyncio.get_event_loop()


loop.create_task(go())
loop.run_forever()
