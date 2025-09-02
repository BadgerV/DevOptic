import { Table, Tag, Tooltip, Typography } from "antd";
import { ExecutionHistory } from "../../types";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  SyncOutlined,
  ExclamationCircleOutlined,
  DownOutlined,
  RightOutlined,
  PlayCircleOutlined,
} from "@ant-design/icons";
import { useDispatch } from "react-redux";
import { openStatusModal } from "../../store";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import duration from "dayjs/plugin/duration";

dayjs.extend(relativeTime);
dayjs.extend(duration);

const { Text } = Typography;

interface HistoryTableProps {
  history: ExecutionHistory[];
  filter: "all" | "running" | "success" | "failed" | "pending";
  searchTerm: string;
}

const HistoryTable = ({ history, filter, searchTerm }: HistoryTableProps) => {
  const dispatch = useDispatch();

  const filteredHistory = history
    .filter((record) => filter === "all" || record.status === filter)
    .filter((record) =>
      searchTerm
        ? record.requester_name
            .toLowerCase()
            .includes(searchTerm.toLowerCase()) ||
          record.pipeline_run_id
            .toLowerCase()
            .includes(searchTerm.toLowerCase()) ||
          record.macro_service_name
            ?.toLowerCase()
            .includes(searchTerm.toLowerCase())
        : true
    );

  const getStatusConfig = (status: string) => {
    const configs = {
      success: {
        icon: <CheckCircleOutlined />,
        color: "success",
        label: "Completed",
      },
      completed: {
        icon: <CheckCircleOutlined />,
        color: "success",
        label: "Completed",
      },
      failed: {
        icon: <CloseCircleOutlined />,
        color: "error",
        label: "Failed",
      },
      running: {
        icon: <SyncOutlined spin />,
        color: "processing",
        label: "Running",
      },
      pending: {
        icon: <ClockCircleOutlined />,
        color: "default",
        label: "Pending",
      },
    };
    return (
      configs[status as keyof typeof configs] || {
        icon: <ExclamationCircleOutlined />,
        color: "default",
        label: status,
      }
    );
  };

  const formatDuration = (startTime: string, endTime: string) => {
    if (!startTime || !endTime) return "-";
    const start = dayjs(startTime);
    const end = dayjs(endTime);
    const diff = end.diff(start);

    if (diff < 60000) {
      // Less than 1 minute
      return `${Math.round(diff / 1000)}s`;
    } else if (diff < 3600000) {
      // Less than 1 hour
      return `${Math.round(diff / 60000)}m`;
    } else if (diff < 86400000) {
      // Less than 1 day
      const hours = Math.floor(diff / 3600000);
      const minutes = Math.round((diff % 3600000) / 60000);
      return `${hours}h ${minutes}m`;
    } else {
      const days = Math.floor(diff / 86400000);
      const hours = Math.round((diff % 86400000) / 3600000);
      return `${days}d ${hours}h`;
    }
  };

  const renderMicroServices = (services: string[]) => {
    if (!services || services.length === 0)
      return <Text type="secondary">No dependencies</Text>;

    const displayLimit = 2;
    const visibleServices = services.slice(0, displayLimit);
    const hiddenCount = services.length - displayLimit;

    return (
      <div>
        {visibleServices.map((service, index) => (
          <Tag key={index}  style={{ marginBottom: 2 }}>
            {service}
          </Tag>
        ))}
        {hiddenCount > 0 && (
          <Tooltip title={services.slice(displayLimit).join(", ")}>
            <Tag  color="blue">
              +{hiddenCount} more
            </Tag>
          </Tooltip>
        )}
      </div>
    );
  };

  const columns = [
    {
      title: "Execution Summary",
      key: "summary",
      width: 280,
      render: (_: any, record: any) => (
        <div>
          <div style={{ fontWeight: 500, marginBottom: 4 }}>
            <PlayCircleOutlined style={{ marginRight: 6, color: "#1890ff" }} />
            {record.macro_service_name || "Unknown Service"}
          </div>
          <Text type="secondary" style={{ fontSize: 12 }}>
            by {record.requester_name} • Started{" "}
            {dayjs(record.started_at).fromNow()}
          </Text>
          {record.error_message && (
            <div style={{ marginTop: 4 }}>
              <Text style={{ fontSize: 12, color: "#ff4d4f" }}>
                <ExclamationCircleOutlined style={{ marginRight: 4 }} />
                {record.error_message.length > 40
                  ? `${record.error_message.substring(0, 40)}...`
                  : record.error_message}
              </Text>
            </div>
          )}
        </div>
      ),
    },
    {
      title: "Dependencies",
      key: "dependencies",
      width: 180,
      render: (_: any, record: any) =>
        renderMicroServices(record.micro_service_names),
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: 120,
      sorter: (a: ExecutionHistory, b: ExecutionHistory) =>
        a.status.localeCompare(b.status),
      render: (status: string) => {
        const statusConfig = getStatusConfig(status);
        return (
          <Tag icon={statusConfig.icon} color={statusConfig.color}>
            {statusConfig.label}
          </Tag>
        );
      },
    },
    {
      title: "Duration",
      key: "duration",
      width: 100,
      render: (_: any, record: any) => (
        <Text style={{ fontSize: 12, fontFamily: "monospace" }}>
          {formatDuration(record.started_at, record.completed_at)}
        </Text>
      ),
    },
    {
      title: "Approval Chain",
      key: "approval",
      width: 140,
      render: (_: any, record: any) => (
        <div style={{ fontSize: 12 }}>
          <div>
            <Text type="secondary">Approved by:</Text>
          </div>
          <Text strong>
            {record.approver_name || (
              <Text type="secondary" italic>
                System
              </Text>
            )}
          </Text>
        </div>
      ),
    },
    {
      title: "Completed",
      key: "completed",
      width: 120,
      sorter: (a: ExecutionHistory, b: ExecutionHistory) =>
        (a.completed_at || "").localeCompare(b.completed_at || ""),
      render: (_: any, record: any) => (
        <div>
          {record.completed_at ? (
            <Tooltip
              title={dayjs(record.completed_at).format("MMM D, YYYY h:mm A")}
            >
              <Text style={{ fontSize: 12 }}>
                {dayjs(record.completed_at).fromNow()}
              </Text>
            </Tooltip>
          ) : (
            <Text type="secondary" style={{ fontSize: 12 }}>
              In Progress
            </Text>
          )}
        </div>
      ),
    },
  ];

  const expandedRowRender = (record: any) => (
    <div
      style={{
        margin: "16px 0",
        padding: "16px",
        backgroundColor: "#fafafa",
        borderRadius: "6px",
      }}
    >
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))",
          gap: "16px",
        }}
      >
        <div>
          <Text strong>Execution Details:</Text>
          <div style={{ marginTop: "8px", fontSize: "12px" }}>
            <div>
              <Text type="secondary">Pipeline Run ID:</Text>
            </div>
            <div style={{ fontFamily: "monospace", marginBottom: "8px" }}>
              {record.pipeline_run_id}
            </div>
            <div>
              <Text type="secondary">Pipeline Unit ID:</Text>
            </div>
            <div style={{ fontFamily: "monospace" }}>
              {record.pipeline_unit_id}
            </div>
          </div>
        </div>

        <div>
          <Text strong>Timeline:</Text>
          <div style={{ marginTop: "8px" }}>
            <div style={{ marginBottom: "4px" }}>
              <Text type="secondary">Started:</Text>{" "}
              {dayjs(record.started_at).format("MMM D, YYYY h:mm A")}
            </div>
            {record.completed_at && (
              <div style={{ marginBottom: "4px" }}>
                <Text type="secondary">Completed:</Text>{" "}
                {dayjs(record.completed_at).format("MMM D, YYYY h:mm A")}
              </div>
            )}
            <div>
              <Text type="secondary">Total Duration:</Text>{" "}
              <Text strong>
                {formatDuration(record.started_at, record.completed_at)}
              </Text>
            </div>
          </div>
        </div>

        <div>
          <Text strong>Participants:</Text>
          <div style={{ marginTop: "8px" }}>
            <div style={{ marginBottom: "4px" }}>
              <Text type="secondary">Requester:</Text>{" "}
              <Tag>{record.requester_name}</Tag>
            </div>
            <div>
              <Text type="secondary">Approver:</Text>{" "}
              <Tag
              
                color={record.approver_name ? "blue" : "default"}
              >
                {record.approver_name || "System"}
              </Tag>
            </div>
          </div>
        </div>
      </div>

      {record.micro_service_names && record.micro_service_names.length > 0 && (
        <div style={{ marginTop: "16px" }}>
          <Text strong>All Dependencies:</Text>
          <div style={{ marginTop: "8px" }}>
            {record.micro_service_names.map(
              (service: string, index: number) => (
                <Tag key={index} style={{ marginBottom: "4px" }}>
                  {service}
                </Tag>
              )
            )}
          </div>
        </div>
      )}

      {record.error_message && (
        <div style={{ marginTop: "16px" }}>
          <Text strong style={{ color: "#ff4d4f" }}>
            Error Details:
          </Text>
          <div
            style={{
              marginTop: "8px",
              padding: "12px",
              backgroundColor: "#fff2f0",
              border: "1px solid #ffccc7",
              borderRadius: "4px",
              fontFamily: "monospace",
              fontSize: "12px",
              color: "#a8071a",
            }}
          >
            {record.error_message}
          </div>
        </div>
      )}
    </div>
  );

  return (
    <Table
      dataSource={filteredHistory}
      columns={columns}
      rowKey="id"
      pagination={{ pageSize: 10 }}
      onRow={(record: any) => ({
        onClick: () => {
          console.log(record);
          dispatch(openStatusModal(record.pipeline_run_id));
        },
        style: { cursor: "pointer" },
      })}
      expandable={{
        expandedRowRender,
        expandIcon: ({ expanded, onExpand, record }) =>
          expanded ? (
            <DownOutlined
              onClick={(e) => {
                e.stopPropagation();
                onExpand(record, e);
              }}
            />
          ) : (
            <RightOutlined
              onClick={(e) => {
                e.stopPropagation();
                onExpand(record, e);
              }}
            />
          ),
      }}
      size="middle"
    />
  );
};

export default HistoryTable;
