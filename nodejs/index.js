let nextPid = 0
const mailboxes = {}
const http = require("http");

//testing: curl -X POST http://localhost:3125 -H 'Content-Type: application/json'  -d '[0, "hello", "world from CURL"]'
const myNode = { nodeID: null, nodeLocators: {} }

const createNode = (options = {}) => {
    myNode.nodeID = Math.floor(Math.random() * 16777215).toString(16);

    if (options.createHttpServer) {
        const requestListener = function (req, res) {
            let body = [];
            req.on('data', (chunk) => {
                body.push(chunk);
            }).on('end', () => {
                body = Buffer.concat(body).toString();
                console.log(body);
                [pid, ...msg] = JSON.parse(body)
                send({ pid, myNode }, ...msg)
                res.writeHead(200);
                res.end();
            });

        };
        const host = options.createHttpServer.host;
        const port = options.createHttpServer.port;
        const server = http.createServer(requestListener);
        server.listen(port, host, () => {
            console.log(`Server is running on http://${host}:${port}`);
        });
        myNode.nodeLocators.http = `http://${host}:${port}`
    }
    return ({ pid }, ...msg) => send({ pid, node: myNode }, ...msg)
}

const send = ({ pid, node }, ...msg) => {
    if (node && node.nodeID !== myNode.nodeID) {

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

createNode({ createHttpServer: { port: 3125, host: "localhost" } })

const A = spawn(async function () {
    console.log("A", this.pid);
    const m = await this.receive()
    console.log("A got", m);
    await this.loop({
        hello(s) {
            console.log("A got hello msg in loop with: " + s);
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
