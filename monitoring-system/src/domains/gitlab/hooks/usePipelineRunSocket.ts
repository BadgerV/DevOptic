import { useEffect, useState } from "react";
import { Subscription } from "rxjs";
import {
  createPipelineSocket,
  PipelineSocketService,
} from "@/shared/services/real-time/pipelineSocket";

interface StatusUpdate {
  pipeline_run_id: string;
  current_service_id?: string;
  status: "pending" | "accepted" | "running" | "success" | "rejected";
  message: string;
  timestamp: string;
  approver_id?: string;
  comment?: string;
}

export const usePipelineRunSocket = (pipelineRunId: string) => {
  const [statusUpdates, setStatusUpdates] = useState<StatusUpdate[]>([]);

  useEffect(() => {
    if (!pipelineRunId) return;

    const socketService: PipelineSocketService =
      createPipelineSocket(pipelineRunId);
    const subscription: Subscription = socketService.connect().subscribe({
      next: (message) => {
        setStatusUpdates((prev) => [...prev, message.payload]);
      },
      error: (error) => {
        console.error("WebSocket subscription error:", error);
      },
      complete: () => {
        console.log("WebSocket subscription completed");
      },
    });

    return () => {
      socketService.disconnect();
      subscription.unsubscribe();
    };
  }, [pipelineRunId]);

  return statusUpdates;
};
