// @ts-nocheck

import {
  Modal,
  Timeline,
  Spin,
  Typography,
  message,
  Divider,
  Space,
  Skeleton,
} from "antd";
import { useDispatch, useSelector } from "react-redux";
import { closeStatusModal, useGetPipelineRunStatusQuery } from "../../store";
import { RootState } from "@/app/store";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  SyncOutlined,
} from "@ant-design/icons";
import { useEffect } from "react";
import { usePipelineRunSocket } from "../../hooks/usePipelineRunSocket";
import { PipelineRunStatus } from "../../types";

const { Text, Title } = Typography;

const PipelineStatusModal = () => {
  const dispatch = useDispatch();
  const isOpen = useSelector(
    (state: RootState) => state.gitlab.isStatusModalOpen
  );
  const pipelineRunId = useSelector(
    (state: RootState) => state.gitlab.statusModalPipelineRunId
  );
  const {
    data: pipelineRunStatus,
    isLoading,
    error,
  } = useGetPipelineRunStatusQuery(pipelineRunId || "", {
    skip: !pipelineRunId,
  });
  const [messageApi, contextHolder] = message.useMessage();
  const statusUpdates = usePipelineRunSocket(pipelineRunId || "");

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: "error",
        content: "Failed to load pipeline status. Please try again.",
      });
    }
  }, [error, messageApi]);

  useEffect(() => {
    console.log("These are the status updates", statusUpdates);
  }, [statusUpdates]);

  const getStatusItem = (update: PipelineRunStatus, isInitial = false) => {
    const status = update.status || "pending";
    const statusProps = {
      success: {
        icon: <CheckCircleOutlined />,
        color: "#52c41a",
        label: "Success",
      },
      completed: {
        icon: <CheckCircleOutlined />,
        color: "#52c41a",
        label: "Completed",
      },
      failed: {
        icon: <CloseCircleOutlined />,
        color: "#ff4d4f",
        label: "Failed",
      },
      running: {
        icon: <SyncOutlined spin />,
        color: "#fa8c16",
        label: "Running",
      },
      pending: {
        icon: <ClockCircleOutlined />,
        color: "#d9d9d9",
        label: "Pending",
      },
      accepted: {
        icon: <CheckCircleOutlined />,
        color: "#1890ff",
        label: "Accepted",
      },
    }[status] || {
      icon: <ClockCircleOutlined />,
      color: "#d9d9d9",
      label: "Pending",
    };

    return (
      <Timeline.Item dot={statusProps.icon} color={statusProps.color}>
        <div style={{ display: "flex", flexDirection: "column" }}>
          <Text strong style={{ fontSize: 15 }}>
            {isInitial
              ? `${update.macro_service_name} Pipeline`
              : update.current_service_id ||
                update.macro_service_name ||
                "Pipeline Service"}
          </Text>
          <Text type="secondary" style={{ fontSize: 13 }}>
            {statusProps.label}
          </Text>
          {update.message && (
            <Text type="secondary" style={{ fontStyle: "italic" }}>
              {update.message}
            </Text>
          )}
          <Text type="secondary" style={{ fontSize: 12 }}>
            {isInitial
              ? `Created: ${new Date(
                  update.created_at || update.timestamp
                ).toLocaleString()}`
              : new Date(update.timestamp).toLocaleString()}
          </Text>
          {update.requester_name && (
            <Text type="secondary" style={{ fontSize: 12 }}>
              Requester: {update.requester_name}
            </Text>
          )}
          {update.approver_name &&
            update.approver_name !== update.requester_name && (
              <Text type="secondary" style={{ fontSize: 12 }}>
                Approver: {update.approver_name}
              </Text>
            )}
          {update.comment && (
            <Text
              type="secondary"
              style={{ fontStyle: "italic", fontSize: 12 }}
            >
              Comment: {update.comment}
            </Text>
          )}
          {update.micro_service_names &&
            update.micro_service_names.length > 0 &&
            isInitial && (
              <Text type="secondary" style={{ fontSize: 12, marginTop: 4 }}>
                Services: {update.micro_service_names.join(", ")}
              </Text>
            )}
        </div>
      </Timeline.Item>
    );
  };

  const renderSkeletonTimeline = () => (
    <Timeline style={{ paddingLeft: 8, marginBottom: 24 }}>
      {[1, 2, 3].map((item) => (
        <Timeline.Item key={item} dot={<ClockCircleOutlined />} color="#d9d9d9">
          <div style={{ display: "flex", flexDirection: "column" }}>
            <Skeleton.Input
              style={{ width: 200, height: 20, marginBottom: 4 }}
              active
              size="small"
            />
            <Skeleton.Input
              style={{ width: 100, height: 16, marginBottom: 4 }}
              active
              size="small"
            />
            <Skeleton.Input
              style={{ width: 150, height: 14 }}
              active
              size="small"
            />
          </div>
        </Timeline.Item>
      ))}
    </Timeline>
  );

  // Merge and deduplicate updates
  const mergeUpdates = () => {
    const updateMap = new Map();
    
    // Add initial pipeline status if available
    if (pipelineRunStatus) {
      const initialUpdate = {
        ...pipelineRunStatus,
        timestamp: pipelineRunStatus.created_at,
        message: `Pipeline created for ${pipelineRunStatus.macro_service_name}`,
        isInitial: true,
      };
      updateMap.set('initial', initialUpdate);
    }

    // Add WebSocket updates, parsing JSON and merging with existing data
    statusUpdates.forEach((update) => {
      try {
        const parsedUpdate = JSON.parse(update);
        const key = parsedUpdate.current_service_id || parsedUpdate.macro_service_name || 'pipeline';
        
        // If we already have an update for this service, merge the data
        if (updateMap.has(key)) {
          const existing = updateMap.get(key);
          updateMap.set(key, {
            ...existing,
            ...parsedUpdate,
            // Keep the more recent timestamp
            timestamp: parsedUpdate.timestamp || existing.timestamp,
          });
        } else {
          updateMap.set(key, parsedUpdate);
        }
      } catch (e) {
        console.error('Failed to parse WebSocket update:', e);
      }
    });

    // Convert back to array and sort by timestamp
    return Array.from(updateMap.values()).sort(
      (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
  };

  const allUpdates = mergeUpdates();

  return (
    <>
      {contextHolder}
      <Modal
        title={
          <Title level={4} style={{ margin: 0 }}>
            Pipeline Run Status
          </Title>
        }
        open={isOpen}
        onCancel={() => dispatch(closeStatusModal())}
        footer={null}
        width={650}
        centered
        styles={{ padding: "20px 24px" }}
      >
        {isLoading && renderSkeletonTimeline()}

        {!isLoading && pipelineRunStatus && (
          <div>
            {allUpdates.length === 1 && 
             allUpdates[0].status !== 'completed' && 
             allUpdates[0].status !== 'success' && 
             allUpdates[0].status !== 'failed' && (
              <Text
                type="secondary"
                style={{ display: "block", marginBottom: 16 }}
              >
                Waiting for pipeline updates...
              </Text>
            )}
            <Timeline style={{ paddingLeft: 8, marginBottom: 24 }}>
              {allUpdates.map((update, index) =>
                getStatusItem(update, update.isInitial || false)
              )}
            </Timeline>

            <Divider style={{ margin: "16px 0" }} />
          </div>
        )}
      </Modal>
    </>
  );
};

export default PipelineStatusModal;