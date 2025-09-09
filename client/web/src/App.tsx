import React, { useEffect, useState } from "react";
import { createWS, Message } from "./utils/ws";

export default function App() {
  const [usernameInput, setUsernameInput] = useState(""); // input field
  const [username, setUsername] = useState("");           // confirmed username
  const [recipient, setRecipient] = useState("");
  const [msgInput, setMsgInput] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);

  // Open WebSocket only once username is confirmed
  useEffect(() => {
    if (username) {
      const socket = createWS(username);
      socket.onmessage = (event) => {
        const msg: Message = JSON.parse(event.data);
        setMessages((prev) => [...prev, msg]);
      };
      setWs(socket);

      return () => {
        socket.close();
      };
    }
  }, [username]);

  const sendMessage = () => {
    if (!recipient || !msgInput || !ws) return;
    const msg: Message = {
      sender: username,
      recipient,
      content: msgInput,
      timestamp: new Date().toISOString(),
    };
    ws.send(JSON.stringify(msg));
    setMessages((prev) => [...prev, msg]);
    setMsgInput("");
  };

  if (!username) {
    // Username input page
    return (
      <div style={{ padding: "20px" }}>
        <h2>Enter your username</h2>
        <input
          value={usernameInput}
          onChange={(e) => setUsernameInput(e.target.value)}
        />
        <button onClick={() => setUsername(usernameInput)}>Connect</button>
      </div>
    );
  }

  // Main chat UI
  return (
    <div style={{ padding: "20px" }}>
      <h2>Wispr Web Chat</h2>

      <div>
        Recipient:{" "}
        <input
          value={recipient}
          onChange={(e) => setRecipient(e.target.value)}
        />
      </div>

      <div style={{ marginTop: "10px" }}>
        <input
          value={msgInput}
          onChange={(e) => setMsgInput(e.target.value)}
          placeholder="Type a message"
        />
        <button onClick={sendMessage}>Send</button>
      </div>

      <div
        style={{
          marginTop: "20px",
          maxHeight: "300px",
          overflowY: "auto",
          border: "1px solid #ccc",
          padding: "10px",
        }}
      >
        {messages.map((m, idx) => (
          <div key={idx}>
            [{m.timestamp}] <b>{m.sender}</b> → {m.recipient}: {m.content}
          </div>
        ))}
      </div>
    </div>
  );
}
