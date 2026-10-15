const socket = new WebSocket('ws://localhost:8080/ws');

const statusElement = document.getElementById('status');
const tempElement = document.getElementById('temperatura');
const umidElement = document.getElementById('umidade');
const lumElement = document.getElementById('luminosidade');

socket.addEventListener('open', () => {
    console.log('Conectado ao WebSocket!');
    statusElement.textContent = 'Conectado e recebendo dados...';
    statusElement.style.color = 'green';
});

socket.addEventListener('message', (event) => {
    console.log('Mensagem recebida do servidor:', event.data);

    try {
        const data = JSON.parse(event.data);

        if (data.temperatura !== undefined) {
            tempElement.textContent = `${data.temperatura.toFixed(1)} °C`;
        }
        if (data.umidade !== undefined) {
            umidElement.textContent = `${data.umidade.toFixed(1)} %`;
        }
        if (data.luminosidade !== undefined) {
            lumElement.textContent = `${data.luminosidade}`;
        }
    } catch (e) {
        console.error('Erro ao processar a mensagem JSON:', e);
    }
});

socket.addEventListener('error', (err) => {
    console.error('Erro no WebSocket:', err);
    statusElement.textContent = 'Erro de conexão com o servidor.';
    statusElement.style.color = 'red';
});

socket.addEventListener('close', () => {
    console.log('Conexão WebSocket fechada.');
    statusElement.textContent = 'Conexão fechada.';
    statusElement.style.color = 'gray';
});
