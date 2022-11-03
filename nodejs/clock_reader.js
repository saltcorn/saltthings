const { createNode, spawn, send } = require("./index")


createNode({ createHttpServer: { port: 3145, host: "localhost" } })

const clockPID = {
    pid: 0,
    node: {
        nodeLocators: { http: "http://localhost:3135" }
    }
}

spawn({
    __init() {
        setInterval(() => send(this, "tick"), 1000)
    },
    tick() {
        console.log("tick");
        send(clockPID, "get_time", this, "tock")
    },
    tock(time) {
        console.log("tock", time);

    }
})


