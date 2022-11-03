const { createNode, spawn, send } = require("./index")


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
