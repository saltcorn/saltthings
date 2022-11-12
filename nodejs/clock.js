const { createNode, spawn, send } = require("./index");

createNode({ createHttpServer: { port: 3135, host: "localhost" } });

spawn({
  processName: "clock",
  get_time(p, msg) {
    console.log("get_time");
    send(p, msg, new Date());
  },
});
