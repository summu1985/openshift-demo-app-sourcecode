// load-generator.js
const WebSocket = require('ws');

const WEBSOCKET_URL = process.argv[2] || 'ws://localhost:8080';
const TOTAL_CONNECTIONS = parseInt(process.argv[3], 10) || 100;
const connections = [];

console.log(`Opening ${TOTAL_CONNECTIONS} WebSocket connections to ${WEBSOCKET_URL}...`);

for (let i = 0; i < TOTAL_CONNECTIONS; i++) {
    const ws = new WebSocket(WEBSOCKET_URL);

    ws.on('open', () => {
        console.log(`Connection ${i + 1} opened`);
        // Optionally send ping to keep alive
        setInterval(() => {
            if (ws.readyState === WebSocket.OPEN) {
                ws.send('ping');
            }
        }, 30000); // every 30s
    });

    ws.on('error', (err) => {
        console.error(`Connection ${i + 1} error:`, err.message);
    });

    ws.on('close', () => {
        console.log(`Connection ${i + 1} closed`);
    });

    connections.push(ws);
}

// Graceful shutdown
process.on('SIGINT', () => {
    console.log('\nClosing all WebSocket connections...');
    connections.forEach(ws => ws.close());
    process.exit();
});
