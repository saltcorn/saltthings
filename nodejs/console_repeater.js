const { createNode, spawn, send } = require("./index");

createNode({ createHttpServer: { port: 3155, host: "localhost" } });

spawn({
  processName: "console",
  print(s) {
    console.log("console>", s);
  },
});
