let nextPid = 0;
const mailboxes = {};
const http = require("http");

//testing: curl -X POST http://localhost:3125 -H 'Content-Type: application/json'  -d '[0, "hello", "world from CURL"]'
const myNode = { nodeID: null, nodeLocators: {}, registeredProcesses: {} };

const createNode = (options = {}) => {
  myNode.nodeID = Math.floor(Math.random() * 16777215).toString(16);

  if (options.createHttpServer) {
    const requestListener = function (req, res) {
      let body = [];
      req
        .on("data", (chunk) => {
          body.push(chunk);
        })
        .on("end", () => {
          body = Buffer.concat(body).toString();
          //console.log(body);
          [proc, ...msg] = JSON.parse(body);
          if (proc.processName) {
            const pid = myNode.registeredProcesses[proc.processName];
            send({ pid, myNode }, ...msg);
          } else send({ pid: proc, myNode }, ...msg);
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
    myNode.nodeLocators.http = `http://${host}:${port}`;
  }
  return ({ pid }, ...msg) => send({ pid, node: myNode }, ...msg);
};

const sendHTTP = (httpLoc, pid, ...msg) => {
  const url = new URL(httpLoc);
  var options = {
    host: url.hostname,
    path: url.pathname,
    port: url.port,
    method: "POST",
  };

  var req = http.request(options);
  req.on("error", (e) => {
    console.error(e);
  });
  req.write(JSON.stringify([pid, ...msg]));
  req.end();
};

const send = ({ pid, node }, ...msg) => {
  if (node && node.nodeID !== myNode.nodeID) {
    if (node.nodeLocators.http) sendHTTP(node.nodeLocators.http, pid, ...msg);
    return;
  }
  if (mailboxes[pid].resolver) {
    const r = mailboxes[pid].resolver;
    mailboxes[pid].resolver = null;
    r(msg);
  } else mailboxes[pid].queue.push(msg);
};

const spawn = (f, options = {}) => {
  const pid = nextPid;
  nextPid++;
  mailboxes[pid] = { queue: [], resolver: null };
  const receive = () => {
    if (mailboxes[pid].queue.length > 0) return mailboxes[pid].queue.shift();
    else {
      if (mailboxes[pid].resolver)
        throw new Error("Double resolver in PID " + pid);
      return new Promise((resolve) => {
        mailboxes[pid].resolver = resolve;
      });
    }
  };
  const node = myNode;
  const that = { receive, pid, node, mbox: { pid, node } };

  //console.log({ that });
  if (f.call) f.call(that);
  else {
    if (f.processName) myNode.registeredProcesses[f.processName] = pid;
    if (f.__init) f.__init.call(that);
    loop(that, f);
  }
  return { pid, node, send: (msg, arg) => send({ pid, node }, msg, arg) };
};

const registerProcess = (name, { pid }) =>
  (myNode.registeredProcesses[name] = pid);

async function loop(that, dispatch) {
  while (true) {
    const [nm, ...args] = await that.receive();
    const f = dispatch[nm];
    if (f) {
      await f.call(that, ...args);
    } else throw new Error("Unknown message in loop: " + nm);
  }
}

module.exports = { spawn, createNode, send, registerProcess, loop };
