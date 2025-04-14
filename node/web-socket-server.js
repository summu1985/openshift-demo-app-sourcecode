// ws-server.js
const WebSocket = require('ws');
const http = require('http');
const express = require('express');
const app = express();

let connectionCount = 0;

const server = http.createServer(app);
const wss = new WebSocket.Server({ server });

wss.on('connection', ws => {
    connectionCount++;
    ws.on('close', () => {
        connectionCount--;
    });
});

// /metrics endpoint
app.get('/metrics', (req, res) => {
    res.set('Content-Type', 'text/plain');
    res.send(`websocket_connection_count ${connectionCount}\n`);
});

server.listen(8080, () => {
    console.log('WebSocket server running on port 8080');
});
