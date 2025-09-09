export interface Message {
  sender: string;
  recipient: string;
  content: string;
  timestamp: string;
}

export function createWS(username: string): WebSocket {
  return new WebSocket(`ws://localhost:8080/ws?username=${username}`);
}
