import {
  Modal,
  Descriptions,
  Tag,
  Button,
  Space,
  Typography,
  message,
  Badge,
  Row,
  Col,
  Card,
  Spin,
  Alert,
} from "antd";
import {
  ApiOutlined,
  ClockCircleOutlined,
  GlobalOutlined,
  ReloadOutlined,
  ContainerOutlined,
  TagsOutlined,
  UserOutlined,
  EditOutlined,
} from "@ant-design/icons";
import { useEffect, useState } from "react";
import { useGetEndpointByIDQuery, useUpdateEndpointMutation } from "../../store/endpointMonitoringApi";
import { EndpointsEssentials } from "../../types";
import { checkEndpointLive } from "../../services/endpointMonitoringApiCalls";
import ServiceCreationModal from "../serviceCreationModal/ServiceCreationModal";

const { Title, Text, Paragraph } = Typography;

interface ServiceData {
  id?: string;
  service_name: string;
  url?: string;
  server_name: string;
  api_method?: string;
  expected_code?: number;
  endpoint_id?: string;
  total_checks: number;
  avg_latency: number;
  successful_checks: number;
  uptime_percentage?: number;
  last_run?: "success" | "failure" | string;
  failure_count?: number;
  gitlab_url?: string;
  docker_container_name?: string;
  kubernetes_pod_name?: string;
  tags?: string[];
  description?: string;
  has_been_modified?: boolean;
  created_at?: string;
  updated_at?: string;
  last_changed_by?: string;
  callType?: "API" | "Service";
  severity?: number;
  downtime_count?: number;
  uptime?: number;
}

interface ServiceDetailsModalProps {
  visible: boolean;
  onClose: () => void;
  id: string | undefined;
}

interface LiveCheckResult {
  status: "success" | "failure";
  response_time?: number;
  status_code?: number;
  timestamp: string;
  error?: string;
}

const ServiceDetailsModal: React.FC<ServiceDetailsModalProps> = ({
  visible,
  onClose,
  id,
}) => {
  const [isChecking, setIsChecking] = useState(false);
  const [lastCheckResult, setLastCheckResult] = useState<LiveCheckResult | null>(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [formData, setFormData] = useState<any>({});

  const { data, isLoading, isError } = useGetEndpointByIDQuery(id as any, {
    skip: !id,
  });

  const [updateEndpoint, { isLoading: isUpdating }] = useUpdateEndpointMutation();

  const normalize = (raw: any): ServiceData => {
    const uptime = Number(raw?.UptimePercentage ?? 0);
    const failureCount = Number(raw?.FailureCount ?? 0);

    const severity =
      uptime >= 99
        ? 0
        : uptime >= 95
        ? 1
        : uptime >= 90
        ? 2
        : failureCount >= 4
        ? 3
        : 2;

    return {
      id: raw?.ID,
      service_name: raw?.ServiceName ?? "-",
      url: raw?.URL,
      server_name: raw?.ServerName ?? "-",
      api_method: raw?.APIMethod,
      expected_code: raw?.ExpectedCode,
      endpoint_id: raw?.EndpointID,
      total_checks: raw?.TotalChecks ?? 0,
      avg_latency: raw?.AvgLatency ?? 0,
      successful_checks: raw?.SuccessfulChecks ?? 0,
      uptime_percentage: uptime,
      last_run: raw?.LastRunSucceeded,
      failure_count: failureCount,
      gitlab_url: raw?.GitlabURL,
      docker_container_name: raw?.DockerContainerName,
      kubernetes_pod_name: raw?.KubernetesPodName,
      tags: raw?.Tags ?? [],
      description: raw?.Description,
      has_been_modified: raw?.HasBeenModified,
      created_at: raw?.CreatedAt,
      updated_at: raw?.UpdatedAt,
      last_changed_by: raw?.LastChangedBy,
      callType: raw?.APIMethod ? "API" : "Service",
      severity,
      downtime_count: failureCount,
      uptime,
    };
  };

  const serviceData: ServiceData | any = data?.data;

  useEffect(() => {
    if (serviceData) {
      setFormData({
        service_name: serviceData.service_name,
        server_name: serviceData.server_name,
        url: serviceData.url,
        api_method: serviceData.api_method,
        expected_status: serviceData.expected_code?.toString(),
        gitlab_url: serviceData.gitlab_url,
        docker_container_name: serviceData.docker_container_name,
        kubernetes_pod_name: serviceData.kubernetes_pod_name,
        tags: serviceData.tags,
        description: serviceData.description,
      });
    }
  }, [serviceData]);

  const formatValue = (value: any): string => {
    if (value === null || value === undefined || value === "") {
      return "-";
    }
    return String(value);
  };

  const formatDate = (dateString?: string): string => {
    if (!dateString || dateString === "-") return "-";
    return new Date(dateString).toLocaleString();
  };

  const getStatusColor = (status?: string): "success" | "error" | "default" => {
    switch (status) {
      case "success":
        return "success";
      case "failure":
        return "error";
      default:
        return "default";
    }
  };

  const getUptimeColor = (percentage: number): string => {
    if (percentage >= 99) return "#52c41a";
    if (percentage >= 95) return "#faad14";
    return "#ff4d4f";
  };

  const getSeverityColor = (severity?: number): string => {
    if (!severity || severity === 0) return "#52c41a";
    if (severity === 1) return "#faad14";
    if (severity === 2) return "#fa8c16";
    if (severity >= 3) return "#ff4d4f";
    return "#d9d9d9";
  };

  const getSeverityText = (severity?: number): string => {
    if (!severity || severity === 0) return "Normal";
    if (severity === 1) return "Low";
    if (severity === 2) return "Medium";
    if (severity === 3) return "High";
    if (severity >= 4) return "Critical";
    return "Unknown";
  };

  const handleLiveCheck = async () => {
    setIsChecking(true);
    try {
      if (!serviceData?.id) {
        throw new Error("No endpoint selected");
      }
      const result = await checkEndpointLive(serviceData.id);
      setLastCheckResult(result);
      if (result.status === "success") {
        message.success(
          `Endpoint is live! Response time: ${result.response_time}ms, Status code: ${result.status_code}`
        );
      } else {
        message.error(
          `Endpoint check failed with status code: ${result.status_code || "N/A"}. Error: ${result.error || "Unknown"}`
        );
      }
    } catch (error) {
      message.error(
        `Failed to check endpoint status: ${
          error instanceof Error ? error.message : "Unknown error"
        }`
      );
      setLastCheckResult({
        status: "failure",
        error: error instanceof Error ? error.message : "Unknown error",
        timestamp: new Date().toISOString(),
      });
    } finally {
      setIsChecking(false);
    }
  };

  const handleUpdateSubmit = async (updatedData: any) => {
    try {
      await updateEndpoint({ id: serviceData.id, ...updatedData }).unwrap();
      message.success("Service updated successfully");
      setIsEditModalOpen(false);
    } catch (error) {
      message.error("Failed to update service");
    }
  };

  const handleModalClose = () => {
    setLastCheckResult(null);
    onClose();
  };

  if (isLoading) {
    return (
      <Modal open={visible} onCancel={handleModalClose} footer={null}>
        <Spin tip="Loading endpoint details..." />
      </Modal>
    );
  }

  if (isError || !serviceData) {
    return (
      <Modal open={visible} onCancel={handleModalClose} footer={null}>
        <Alert type="error" message="Failed to load endpoint details" />
      </Modal>
    );
  }

  const uptimeValue = serviceData.uptime_percentage ?? serviceData.uptime ?? 0;
  const failureCount = serviceData.failure_count ?? serviceData.downtime_count ?? 0;

  return (
    <>
      <Modal
        title={
          <div style={{ display: "flex", alignItems: "center" }}>
            <ApiOutlined style={{ marginRight: 8, color: "#1890ff" }} />
            Service Details: {serviceData.service_name}
          </div>
        }
        open={visible}
        onCancel={handleModalClose}
        width={800}
        footer={[
          <Button
            key="edit"
            type="default"
            icon={<EditOutlined />}
            onClick={() => setIsEditModalOpen(true)}
          >
            Edit
          </Button>,
          <Button key="close" onClick={handleModalClose}>
            Close
          </Button>,
        ]}
        styles={{
          body: { padding: "24px" },
        }}
      >
        <div style={{ maxHeight: "70vh", overflowY: "auto" }}>
          <Card
            size="small"
            style={{
              marginBottom: 24,
              backgroundColor: "#f8f9fa",
              border: "1px solid #e9ecef",
            }}
          >
            <Row align="middle" justify="space-between">
              <Col>
                <Space>
                  <Text strong>Endpoint Status Check</Text>
                  {lastCheckResult && (
                    <Badge
                      status={lastCheckResult.status === "success" ? "success" : "error"}
                      text={
                        lastCheckResult.status === "success"
                          ? `Live (${lastCheckResult.response_time}ms)`
                          : "Failed"
                      }
                    />
                  )}
                </Space>
              </Col>
              <Col>
                <Button
                  type="primary"
                  icon={<ReloadOutlined />}
                  loading={isChecking}
                  onClick={handleLiveCheck}
                >
                  {isChecking ? "Checking..." : "Check Now"}
                </Button>
              </Col>
            </Row>
          </Card>

          <Card
            title={
              <span>
                <GlobalOutlined style={{ marginRight: 8 }} />
                Service Information
              </span>
            }
            size="small"
            style={{ marginBottom: 24 }}
          >
            <Descriptions column={2} size="small">
              <Descriptions.Item label="ID">{formatValue(serviceData.id)}</Descriptions.Item>
              <Descriptions.Item label="Service Name">
                {formatValue(serviceData.service_name)}
              </Descriptions.Item>
              <Descriptions.Item label="URL" span={2}>
                {formatValue(serviceData.url) !== "-" ? (
                  <Text copyable>{serviceData.url}</Text>
                ) : (
                  "-"
                )}
              </Descriptions.Item>
              <Descriptions.Item label="Server Name">
                {formatValue(serviceData.server_name)}
              </Descriptions.Item>
              <Descriptions.Item label="API Method">
                {serviceData.api_method ? (
                  <Tag color="blue">{serviceData.api_method}</Tag>
                ) : (
                  "-"
                )}
              </Descriptions.Item>
              <Descriptions.Item label="Call Type">
                {serviceData.callType ? (
                  <Tag color={serviceData.callType === "API" ? "green" : "orange"}>
                    {serviceData.callType}
                  </Tag>
                ) : (
                  "-"
                )}
              </Descriptions.Item>
              <Descriptions.Item label="Expected Status">
                {serviceData.expected_code ? (
                  <Tag color="green">{serviceData.expected_code}</Tag>
                ) : (
                  "-"
                )}
              </Descriptions.Item>
            </Descriptions>
          </Card>

          <Card
            title={
              <span>
                <ClockCircleOutlined style={{ marginRight: 8 }} />
                Monitoring Statistics
              </span>
            }
            size="small"
            style={{ marginBottom: 24 }}
          >
            <Row gutter={16}>
              <Col span={12}>
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="Endpoint ID">
                    {formatValue(serviceData.endpoint_id)}
                  </Descriptions.Item>
                  <Descriptions.Item label="Total Checks">
                    {formatValue(serviceData.total_checks)}
                  </Descriptions.Item>
                  <Descriptions.Item label="Successful Checks">
                    <Text style={{ color: "#52c41a" }}>
                      {formatValue(serviceData.successful_checks)}
                    </Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="Failure Count">
                    <Text style={{ color: "#ff4d4f" }}>{formatValue(failureCount)}</Text>
                  </Descriptions.Item>
                </Descriptions>
              </Col>
              <Col span={12}>
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="Average Latency">
                    {formatValue(serviceData.avg_latency) !== "-"
                      ? `${serviceData.avg_latency}ms`
                      : "-"}
                  </Descriptions.Item>
                  <Descriptions.Item label="Uptime Percentage">
                    <Text style={{ color: getUptimeColor(uptimeValue)}}>
                      {uptimeValue}%
                    </Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="Severity Level">
                    <Badge
                      color={getSeverityColor(serviceData.severity)}
                      text={getSeverityText(serviceData.severity)}
                    />
                  </Descriptions.Item>
                  <Descriptions.Item label="Last Run Status">
                    <Badge
                      status={getStatusColor(serviceData.last_run)}
                      text={
                        serviceData.last_run === "success"
                          ? "Success"
                          : serviceData.last_run === "failure"
                          ? "Failed"
                          : formatValue(serviceData.last_run)
                      }
                    />
                  </Descriptions.Item>
                </Descriptions>
              </Col>
            </Row>
          </Card>

          <Card
            title={
              <span>
                <ContainerOutlined style={{ marginRight: 8 }} />
                Infrastructure & Metadata
              </span>
            }
            size="small"
            style={{ marginBottom: 24 }}
          >
            <Descriptions column={1} size="small">
              <Descriptions.Item label="GitLab URL">
                {formatValue(serviceData.gitlab_url) !== "-" ? (
                  <a href={serviceData.gitlab_url} target="_blank" rel="noopener noreferrer">
                    {serviceData.gitlab_url}
                  </a>
                ) : (
                  "-"
                )}
              </Descriptions.Item>
              <Descriptions.Item label="Docker Container">
                <Text code>{formatValue(serviceData.docker_container_name)}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Kubernetes Pod">
                <Text code>{formatValue(serviceData.kubernetes_pod_name)}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Tags">
                {serviceData.tags && serviceData.tags.length > 0 ? (
                  <Space wrap>
                    {serviceData.tags.map((tag: any) => (
                      <Tag key={tag} icon={<TagsOutlined />}>
                        {tag}
                      </Tag>
                    ))}
                  </Space>
                ) : (
                  "-"
                )}
              </Descriptions.Item>
              <Descriptions.Item label="Description">
                <Paragraph style={{ margin: 0 }}>
                  {formatValue(serviceData.description)}
                </Paragraph>
              </Descriptions.Item>
              <Descriptions.Item label="Modified">
                {serviceData.has_been_modified !== undefined ? (
                  <Tag color={serviceData.has_been_modified ? "orange" : "default"}>
                    {serviceData.has_been_modified ? "Yes" : "No"}
                  </Tag>
                ) : (
                  "-"
                )}
              </Descriptions.Item>
              <Descriptions.Item label="Created At">
                <Text>{formatDate(serviceData.created_at)}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Updated At">
                <Text>{formatDate(serviceData.updated_at)}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Last Changed By">
                <Text>
                  <UserOutlined style={{ marginRight: 4 }} />
                  {formatValue(serviceData.last_changed_by)}
                </Text>
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </div>
      </Modal>

      <ServiceCreationModal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        onSubmit={handleUpdateSubmit}
        formData={formData}
        setFormData={setFormData}
        loading={isUpdating}
        isEditMode={true}
      />
    </>
  );
};

export default ServiceDetailsModal;