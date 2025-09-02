import { io, Socket } from 'socket.io-client';

export type EventCallback = (data: any) => void;

class WebSocketClient {
  private socket: Socket | null = null;
  private url: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private eventHandlers = new Map<string, Set<EventCallback>>();

  constructor() {
    this.url = import.meta.env.VITE_WEBSOCKET_URL || 'ws://localhost:8000';
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.socket = io(this.url, {
          transports: ['websocket'],
          autoConnect: true,
        });

        this.socket.on('connect', () => {
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          resolve();
        });

        this.socket.on('disconnect', () => {
          console.log('WebSocket disconnected');
          this.handleReconnect();
        });

        this.socket.on('connect_error', (error) => {
          console.error('WebSocket connection error:', error);
          reject(error);
        });

        // Setup event forwarding
        this.setupEventForwarding();
      } catch (error) {
        reject(error);
      }
    });
  }

  private setupEventForwarding() {
    if (!this.socket) return;

    // Forward all events to registered handlers
    this.socket.onAny((eventName: string, data: any) => {
      const handlers = this.eventHandlers.get(eventName);
      if (handlers) {
        handlers.forEach(handler => handler(data));
      }
    });
  }

  private handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      setTimeout(() => {
        console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        this.connect().catch(console.error);
      }, 1000 * this.reconnectAttempts);
    }
  }

  on(event: string, callback: EventCallback): () => void {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, new Set());
    }
    this.eventHandlers.get(event)!.add(callback);

    // Return cleanup function
    return () => this.off(event, callback);
  }

  off(event: string, callback?: EventCallback) {
    if (callback) {
      this.eventHandlers.get(event)?.delete(callback);
    } else {
      this.eventHandlers.delete(event);
    }
  }

  emit(event: string, data?: any) {
    this.socket?.emit(event, data);
  }

  disconnect() {
    this.socket?.disconnect();
    this.socket = null;
    this.eventHandlers.clear();
  }

  isConnected(): boolean {
    return this.socket?.connected ?? false;
  }
}

export const websocketClient = new WebSocketClient();
