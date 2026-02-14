import { useState, useEffect, useRef } from 'react';
import './App.css';

function App() {
  const params = new URLSearchParams(window.location.search);
  const roomId = params.get("room") || "SUI_333";
  const playerId = parseInt(params.get("player")) || 1;

  const [board, setBoard] = useState(Array(9).fill("-"));
  const [status, setStatus] = useState("connecting");
  const [winnerId, setWinnerId] = useState(null);
  const socketRef = useRef(null);

  useEffect(() => {
    // Connect to Go Backend
    const socket = new WebSocket(`ws://localhost:8080/ws?room=${roomId}`);
    socketRef.current = socket;

    socket.onopen = () => setStatus("active");
socket.onmessage = (event) => {
  try {
    const data = JSON.parse(event.data);
    console.log("Server Data:", data);

    const boardStr = data.Board || data.board;
    const gameState = data.GameState || data.game_state;
    const winner = data.Winner_id !== undefined ? data.Winner_id : data.winner_id;

    if (boardStr) {
      setBoard(boardStr.split(""));
    }
    if (gameState) {
      setStatus(gameState);
    }
    if (winner !== undefined) {
      setWinnerId(winner);
    }
  } catch (err) {
    console.error("Failed to parse server message:", err);
  }
};
socket.onclose = () => setStatus("disconnected");

    return () => socket.close();
  }, [roomId]);

  const makeMove = (index) => {
    if (status === "finished" || board[index] !== "-") return;

    const move = {
      room_id: roomId,
      player_id: playerId,
      cell_index: index
    };
    
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify(move));
    }
  };

  return (
    <div className="container">
      <div className="header">
        <h1>Room: {roomId}</h1>
        <h3>You are: Player {playerId} ({playerId === 1 ? 'X' : 'O'})</h3>
        <p className={`status ${status}`}>Status: {status}</p>
      </div>

      {status === "finished" && (
        <div className="overlay">
          <h2>Game Over!</h2>
          <p>{winnerId === 0 ? "It's a Draw!" : `Player ${winnerId} Wins!`}</p>
          <button onClick={() => window.location.reload()}>Play Again</button>
        </div>
      )}

      <div className="grid">
        {board.map((cell, i) => (
          <button 
            key={i} 
            className={`cell ${cell !== "-" ? "occupied" : ""}`} 
            onClick={() => makeMove(i)}
          >
            {cell === "-" ? "" : cell}
          </button>
        ))}
      </div>
    </div>
  );
}

export default App;
