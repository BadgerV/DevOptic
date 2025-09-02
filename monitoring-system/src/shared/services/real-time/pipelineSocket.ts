import { Observable, Subject } from "rxjs";

interface WebSocketMessage {
  type: string;
  id: string;
  payload: {
    pipeline_run_id: string;
    current_service_id?: string;
    status: "pending" | "accepted" | "running" | "success" | "rejected";
    message: string;
    timestamp: string;
    approver_id?: string;
    comment?: string;
  };
  timestamp: string;
}

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL_WEB_SOCKET || "http://localhost:8080";

export class PipelineSocketService {
  private ws: WebSocket | null = null;
  private messageSubject = new Subject<WebSocketMessage>();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectInterval = 5000; // 5 seconds
  private pipelineRunId: string;

  constructor(pipelineRunId: string) {
    this.pipelineRunId = pipelineRunId;
  }

  connect(): Observable<WebSocketMessage> {
    const token = localStorage.getItem("token");
    this.ws = new WebSocket(
      `ws://${API_BASE_URL}/api/v1/gitlab/ws/pipeline-runs/${this.pipelineRunId}?token=${token}`
    );

    this.ws.onopen = () => {
      console.log(`WebSocket connected for pipeline run ${this.pipelineRunId}`);
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        if (
          message.type === "pipeline_status_change" &&
          message.id === this.pipelineRunId
        ) {
          this.messageSubject.next(message);
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    this.ws.onclose = () => {
      console.log(`WebSocket closed for pipeline run ${this.pipelineRunId}`);
      this.attemptReconnect();
    };

    return this.messageSubject.asObservable();
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error("Max reconnect attempts reached");
      this.messageSubject.complete();
      return;
    }

    this.reconnectAttempts++;
    setTimeout(() => {
      console.log(`Reconnecting WebSocket (attempt ${this.reconnectAttempts})`);
      this.connect();
    }, this.reconnectInterval);
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
      this.messageSubject.complete();
    }
  }
}

export const createPipelineSocket = (pipelineRunId: string) =>
  new PipelineSocketService(pipelineRunId);
