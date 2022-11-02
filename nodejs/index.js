let nextPid = 0
const mailboxes = {}


const myNode = { nodeID: null, nodeLocators: {} }

const createNode = (options = {}) => {
    myNode.nodeID = Math.floor(Math.random() * 16777215).toString(16);

    return ({ pid }, ...msg) => send({ pid, node: myNode }, ...msg)
}

const send = ({ pid, node }, ...msg) => {
    if (node.nodeID !== myNode.nodeID) {

        return
    }
    if (mailboxes[pid].resolver) {
        const r = mailboxes[pid].resolver
        mailboxes[pid].resolver = null
        r(msg);
    } else
        mailboxes[pid].queue.push(msg);
}

const spawn = (f, options = {}) => {
    const pid = nextPid
    nextPid++
    mailboxes[pid] = { queue: [], resolver: null }
    const receive = () => {
        if (mailboxes[pid].queue.length > 0)
            return mailboxes[pid].queue.shift();
        else {
            if (mailboxes[pid].resolver) throw new Error("Double resolver in PID " + pid)
            return new Promise((resolve) => {
                mailboxes[pid].resolver = resolve
            });

        }

    }
    const node = myNode
    const that = { receive, pid, node }
    const loop = async (dispatch) => {
        while (true) {
            const [nm, arg] = await receive()
            const f = dispatch[nm]
            if (f) {
                await f.call(that, arg)
            } else
                throw new Error("Unknown message in loop: " + nm)
        }
    }
    that.loop = loop
    console.log({ that });
    if (f.call)
        f.call(that)
    else {
        if (f.__init)
            f.__init.call(that)
        loop(f)
    }
    return { pid, node, send: (msg, arg) => send({ pid, node }, msg, arg) }
}

createNode()

const A = spawn(async function () {
    console.log("A", this.pid);
    const m = await this.receive()
    console.log("A got", m);
    await this.loop({
        hello() {
            console.log("A got hello msg in loop");
        },
        foo({ x, y }) {
            console.log("A got foo", x, y);
        }

    })
})
send(A, "hello")
send(A, "hello")
A.send("foo", { x: 4, y: 2 })

const B = spawn({
    __init() {
        this.w = 9
    },
    bar({ z }) {
        console.log("B.bar got", z, this.w);
        A.send("foo", { x: z, y: this.w })
    }
})

send(B, "bar", { z: 7 })
