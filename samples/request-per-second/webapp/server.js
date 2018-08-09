const http = require('http');
const appInsights = require('applicationinsights');

const hostname = '0.0.0.0';
const port = 8080;

const instrumentation_key = process.env.INSTRUMENTATION_KEY
appInsights.setup(instrumentation_key).start();

const server = http.createServer((req, res) => {
    res.statusCode = 200;
    res.setHeader('Content-Type', 'text/plain');
    res.end('how many requests per second');
});

server.listen(port, hostname, () => {
  console.log(`Server running at http://${hostname}:${port}/`);
});
