import { WebSocketMessage, ProgressMessage, ActionMessage, ActionType } from "@/types/websocket";

let notificationsSocket: WebSocket | null = null;
let interactiveSocket: WebSocket | null = null;
const notificationsCallbacks: Set<(message: ProgressMessage) => void> = new Set();
const interactiveCallbacks: Set<(message: ActionMessage) => void> = new Set();

export const connectNotifications = () => {
  if (notificationsSocket?.readyState === WebSocket.OPEN) return;

  notificationsSocket = new WebSocket(`ws://${window.location.host}/ws/notifications`);
  
  notificationsSocket.onmessage = (event) => {
    const message = JSON.parse(event.data) as ProgressMessage;
    if (message.type === "progress") {
      notificationsCallbacks.forEach(callback => callback(message));
    }
  };

  notificationsSocket.onerror = (error) => {
    console.error('Notifications WebSocket error:', error);
  };

  notificationsSocket.onclose = () => {
    console.log('Notifications WebSocket closed');
    // Try to reconnect after 5 seconds
    setTimeout(connectNotifications, 5000);
  };
};

export const connectInteractive = () => {
  if (interactiveSocket?.readyState === WebSocket.OPEN) return;

  interactiveSocket = new WebSocket(`ws://${window.location.host}/ws/interactive`);
  
  interactiveSocket.onmessage = (event) => {
    const message = JSON.parse(event.data) as ActionMessage;
    interactiveCallbacks.forEach(callback => callback(message));
  };

  interactiveSocket.onerror = (error) => {
    console.error('Interactive WebSocket error:', error);
  };

  interactiveSocket.onclose = () => {
    console.log('Interactive WebSocket closed');
    // Try to reconnect after 5 seconds
    setTimeout(connectInteractive, 5000);
  };
};

export const subscribeToNotifications = (callback: (message: ProgressMessage) => void) => {
  notificationsCallbacks.add(callback);
  return () => {
    notificationsCallbacks.delete(callback);
  };
};

export const subscribeToInteractive = (callback: (message: ActionMessage) => void) => {
  interactiveCallbacks.add(callback);
  return () => {
    interactiveCallbacks.delete(callback);
  };
};

export const sendInteractiveMessage = (type: ActionType, action: string, data?: any) => {
  if (interactiveSocket?.readyState === WebSocket.OPEN) {
    interactiveSocket.send(JSON.stringify({ type, action, data }));
  } else {
    console.error('Interactive WebSocket is not connected');
  }
};

export const disconnect = () => {
  notificationsSocket?.close();
  interactiveSocket?.close();
  notificationsSocket = null;
  interactiveSocket = null;
  notificationsCallbacks.clear();
  interactiveCallbacks.clear();
};
