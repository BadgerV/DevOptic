import { Table, Button, Modal, Input, message, Tag, Tooltip, Typography } from "antd";
import { useState } from "react";
import {
  useApproveAuthorizationRequestMutation,
  useRejectAuthorizationRequestMutation,
  openStatusModal,
} from "../../store";
import { AuthorizationRequest } from "../../types";
import { 
  CheckCircleOutlined, 
  CloseCircleOutlined, 
  ClockCircleOutlined,
  InfoCircleOutlined,
  DownOutlined,
  RightOutlined
} from "@ant-design/icons";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "@/app/store";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);


const { Text } = Typography;

const ApprovalsTable = ({ requests, filter }: any) => {
  const dispatch = useDispatch();
  const userId =
    useSelector((state: RootState) => state.auth.user?.userId) || "";
  const [rejectModalVisible, setRejectModalVisible] = useState(false);
  const [approveModalVisible, setApproveModalVisible] = useState(false);
  const [selectedRequestId, setSelectedRequestId] = useState<string | null>(
    null
  );
  const [approveComment, setApproveComment] = useState("");
  const [rejectComment, setRejectComment] = useState("");
  const [approveRequest, { isLoading: isApproving }] =
    useApproveAuthorizationRequestMutation();
  const [rejectRequest, { isLoading: isRejecting }] =
    useRejectAuthorizationRequestMutation();
  const [messageApi, contextHolder] = message.useMessage();

  const filteredRequests =
    filter === "all"
      ? requests
      : requests.filter((r: any) => r.status === filter);

  const handleApprove = (id: string) => {
    setSelectedRequestId(id);
    setApproveModalVisible(true);
  };

  const handleApproveConfirm = async (pipelineRunId: string) => {
    if (!selectedRequestId) return;
    try {
      await approveRequest({
        id: selectedRequestId,
        data: { approver_id: userId, comment: approveComment || " " },
      }).unwrap();
      messageApi.open({
        type: "success",
        content: "Request approved successfully!",
      });
      dispatch(openStatusModal(pipelineRunId));
      setApproveModalVisible(false);
      setApproveComment("");
      setSelectedRequestId(null);
    } catch {
      messageApi.open({
        type: "error",
        content: "Failed to approve request. Please try again.",
      });
    }
  };

  const handleReject = (id: string) => {
    setSelectedRequestId(id);
    setRejectModalVisible(true);
  };

  const handleRejectConfirm = async () => {
    if (!selectedRequestId) return;
    try {
      await rejectRequest({
        id: selectedRequestId,
        data: { comment: rejectComment },
      }).unwrap();
      messageApi.open({
        type: "success",
        content: "Request rejected successfully!",
      });
      setRejectModalVisible(false);
      setRejectComment("");
      setSelectedRequestId(null);
    } catch {
      messageApi.open({
        type: "error",
        content: "Failed to reject request. Please try again.",
      });
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "approved":
      case "accepted":
        return <CheckCircleOutlined style={{ color: "#52c41a" }} />;
      case "rejected":
        return <CloseCircleOutlined style={{ color: "#ff4d4f" }} />;
      case "pending":
        return <ClockCircleOutlined style={{ color: "#fa8c16" }} />;
      default:
        return null;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "approved":
      case "accepted":
        return "success";
      case "rejected":
        return "error";
      case "pending":
        return "warning";
      default:
        return "default";
    }
  };

  const renderMicroServices = (services: string[]) => {
    if (!services || services.length === 0) return <Text type="secondary">None</Text>;
    
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
            <Tag color="blue">
              +{hiddenCount} more
            </Tag>
          </Tooltip>
        )}
      </div>
    );
  };

  const columns = [
    {
      title: "Request Summary",
      key: "summary",
      width: 300,
      render: (_: any, record: any) => (
        <div>
          <div style={{ fontWeight: 500, marginBottom: 4 }}>
            {record.macro_service_name}
          </div>
          <Text type="secondary" style={{ fontSize: 12 }}>
            by {record.requester_name} • {dayjs(record.created_at).fromNow()}
          </Text>
          {record.comment && (
            <div style={{ marginTop: 4 }}>
              <Text style={{ fontSize: 12 }} type="secondary">
                <InfoCircleOutlined style={{ marginRight: 4 }} />
                {record.comment.length > 50 
                  ? `${record.comment.substring(0, 50)}...` 
                  : record.comment}
              </Text>
            </div>
          )}
        </div>
      ),
    },
    {
      title: "Services",
      key: "services",
      width: 200,
      render: (_: any, record: any) => renderMicroServices(record.micro_service_names),
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: 120,
      render: (status: string) => (
        <Tag 
          icon={getStatusIcon(status)} 
          color={getStatusColor(status)}
          style={{ textTransform: 'capitalize' }}
        >
          {status}
        </Tag>
      ),
    },
    {
      title: "Approver",
      key: "approver",
      width: 120,
      render: (_: any, record: any) => (
        record.approver_name ? (
          <Text>{record.approver_name}</Text>
        ) : (
          <Text type="secondary" italic>Unassigned</Text>
        )
      ),
    },
    {
      title: "Last Updated",
      dataIndex: "updated_at",
      key: "updated_at",
      width: 140,
      render: (date: string) => (
        <Tooltip title={dayjs(date).format("MMM D, YYYY h:mm A")}>
          <Text style={{ fontSize: 12 }}>
            {dayjs(date).fromNow()}
          </Text>
        </Tooltip>
      ),
    },
    {
      title: "Actions",
      key: "actions",
      width: 150,
      render: (_: any, record: AuthorizationRequest) => (
        <>
          {record.status === "pending" && (
            <div style={{ display: 'flex', gap: 8 }}>
              <Button
                type="primary"
                size="small"
                onClick={() => handleApprove(record.id)}
                loading={isApproving}
              >
                Approve
              </Button>
              <Button
                type="default"
                size="small"
                onClick={() => handleReject(record.id)}
                loading={isRejecting}
              >
                Reject
              </Button>
            </div>
          )}
        </>
      ),
    },
  ];

  const expandedRowRender = (record: any) => (
    <div style={{ margin: '16px 0', padding: '16px', backgroundColor: '#fafafa', borderRadius: '6px' }}>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div>
          <Text strong>Request ID:</Text>
          <div style={{ fontFamily: 'monospace', fontSize: '12px', marginTop: '4px' }}>
            {record.id}
          </div>
        </div>
        <div>
          <Text strong>Pipeline Run ID:</Text>
          <div style={{ fontFamily: 'monospace', fontSize: '12px', marginTop: '4px' }}>
            {record.pipeline_run_id}
          </div>
        </div>
        <div>
          <Text strong>Created:</Text>
          <div style={{ marginTop: '4px' }}>
            {dayjs(record.created_at).format("MMM D, YYYY h:mm A")}
          </div>
        </div>
        <div>
          <Text strong>Last Updated:</Text>
          <div style={{ marginTop: '4px' }}>
            {dayjs(record.updated_at).format("MMM D, YYYY h:mm A")}
          </div>
        </div>
      </div>
      
      {record.micro_service_names && record.micro_service_names.length > 0 && (
        <div style={{ marginTop: '16px' }}>
          <Text strong>All Micro Services:</Text>
          <div style={{ marginTop: '8px' }}>
            {record.micro_service_names.map((service: string, index: number) => (
              <Tag key={index} style={{ marginBottom: '4px' }}>
                {service}
              </Tag>
            ))}
          </div>
        </div>
      )}
      
      {record.comment && (
        <div style={{ marginTop: '16px' }}>
          <Text strong>Full Comment:</Text>
          <div style={{ 
            marginTop: '8px', 
            padding: '8px', 
            backgroundColor: 'white', 
            border: '1px solid #d9d9d9',
            borderRadius: '4px'
          }}>
            {record.comment}
          </div>
        </div>
      )}
    </div>
  );

  return (
    <>
      {contextHolder}
      <Table
        dataSource={filteredRequests}
        columns={columns}
        rowKey="id"
        pagination={{ pageSize: 10 }}
        expandable={{
          expandedRowRender,
          expandIcon: ({ expanded, onExpand, record }) =>
            expanded ? (
              <DownOutlined onClick={e => onExpand(record, e)} />
            ) : (
              <RightOutlined onClick={e => onExpand(record, e)} />
            )
        }}
        size="middle"
      />
      <Modal
        title="Approve Request"
        open={approveModalVisible}
        onOk={() =>
          handleApproveConfirm(
            requests.find((r: any) => r.id === selectedRequestId)
              ?.pipeline_run_id || ""
          )
        }
        onCancel={() => {
          setApproveModalVisible(false);
          setApproveComment("");
          setSelectedRequestId(null);
        }}
        okText="Approve"
        okButtonProps={{ loading: isApproving }}
      >
        <Input.TextArea
          value={approveComment}
          onChange={(e) => setApproveComment(e.target.value)}
          placeholder="Optional comment"
          rows={4}
        />
      </Modal>
      <Modal
        title="Reject Request"
        open={rejectModalVisible}
        onOk={handleRejectConfirm}
        onCancel={() => {
          setRejectModalVisible(false);
          setRejectComment("");
          setSelectedRequestId(null);
        }}
        okText="Reject"
        okButtonProps={{ loading: isRejecting }}
      >
        <Input.TextArea
          value={rejectComment}
          onChange={(e) => setRejectComment(e.target.value)}
          placeholder="Optional comment"
          rows={4}
        />
      </Modal>
    </>
  );
};

export default ApprovalsTable;