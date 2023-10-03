// run in browser console to send message
const socket = new WebSocket('ws://localhost:8000/ws');
socket.addEventListener('message', (event) => {
  console.log(event.data);
});
socket.send('');
