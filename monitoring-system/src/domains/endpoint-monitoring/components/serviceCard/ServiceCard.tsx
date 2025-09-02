import React, { useState } from "react";
import { Card, Tag, Progress, Statistic, Row, Col } from "antd";
import ServiceDetailsModal from "../serviceModal/ServiceModal"; // Import the modal component
import "./serviceCard.css";

type MonitoringCardProps = {
  service_name: string;
  server_name: string;
  uptime: number;
  callType: "API" | "Service";
  total_checks: number;
  successful_checks: number;
  downtime_count: number;
  avg_latency: number;
  severity: number;

  // Additional optional props for the modal
  id?: string;
  url?: string;
  api_method?: string;
  expected_status_code?: number;
  endpoint_id?: string;
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
};

const ServiceCard: React.FC<MonitoringCardProps> = ({
  service_name,
  server_name,
  uptime,
  callType,
  total_checks,
  successful_checks,
  downtime_count,
  avg_latency,
  severity,
  // Additional props for modal
  id,
  url,
  api_method,
  expected_status_code,
  endpoint_id,
  uptime_percentage,
  last_run,
  failure_count,
  gitlab_url,
  docker_container_name,
  kubernetes_pod_name,
  tags,
  description,
  has_been_modified,
  created_at,
  updated_at,
  last_changed_by,
}) => {
  const [isModalVisible, setIsModalVisible] = useState(false);

  const envColor = server_name === "QA" ? "blue" : "purple";
  const typeColor = callType === "API" ? "green" : "orange";
  const progressColor =
    uptime > 95 ? "#52c41a" : uptime > 80 ? "#fa8c16" : "#ff4d4f";

  // Original uptime-based styling (for severity 0)
  const getOriginalStyles = () => ({
    borderColor: progressColor,
    boxShadow: "0 4px 6px rgba(0, 0, 0, 0.1)",
    background: "linear-gradient(145deg, #ffffff, #f9fafb)",
    animation: "none",
  });

  // Severity-based styling configuration (for severity 1-4 and higher)
  const getSeverityStyles = (severity: number) => {
    const configs = {
      1: {
        borderColor: "rgba(255, 0, 0, 0.15)",
        boxShadow: "0 4px 6px rgba(0, 0, 0, 0.1)",
        background: "linear-gradient(145deg, #ffffff, #f9fafb)",
        animation: "none",
      },
      2: {
        borderColor: "rgba(255, 0, 0, 0.35)",
        boxShadow:
          "0 4px 8px rgba(255, 0, 0, 0.15), 0 0 12px rgba(255, 0, 0, 0.1)",
        background: "linear-gradient(145deg, #fefefe, #fdfbfb)",
        animation: "pulse-soft 3s ease-in-out infinite",
      },
      3: {
        borderColor: "rgba(255, 0, 0, 0.55)",
        boxShadow:
          "0 4px 12px rgba(255, 0, 0, 0.25), 0 0 16px rgba(255, 0, 0, 0.15)",
        background: "linear-gradient(145deg, #fefcfc, #fcf8f8)",
        animation: "pulse-medium 2s ease-in-out infinite",
      },
      4: {
        borderColor: "rgba(255, 0, 0, 0.8)",
        boxShadow:
          "0 4px 16px rgba(255, 0, 0, 0.4), 0 0 24px rgba(255, 0, 0, 0.25)",
        background: "linear-gradient(145deg, #fef9f9, #fcf4f4)",
        animation: "pulse-strong 1.2s ease-in-out infinite",
      },
    };
    // Use severity 4 styles for any severity >= 4
    return severity >= 4
      ? configs[4]
      : configs[severity as keyof typeof configs] || configs[1];
  };

  // Determine which styling to use based on severity
  const cardStyleConfig =
    severity === 0 ? getOriginalStyles() : getSeverityStyles(severity);

  const cardStyle: React.CSSProperties = {
    borderRadius: "12px",
    overflow: severity === 0 ? "hidden" : "visible",
    border: `2px solid ${cardStyleConfig.borderColor}`,
    boxShadow: cardStyleConfig.boxShadow,
    animation: cardStyleConfig.animation,
    background: cardStyleConfig.background,
    transition: "transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out",
    position: "relative",
    cursor: "pointer", // Add cursor pointer to indicate clickability
  };

  // Prepare data for the modal
  const modalData = {
    service_name,
    server_name,
    uptime,
    callType,
    total_checks,
    successful_checks,
    downtime_count,
    avg_latency,
    severity,
    id,
    url,
    api_method,
    expected_status_code,
    endpoint_id,
    uptime_percentage,
    last_run,
    failure_count,
    gitlab_url,
    docker_container_name,
    kubernetes_pod_name,
    tags,
    description,
    has_been_modified,
    created_at,
    updated_at,
    last_changed_by,
  };

  const handleCardClick = () => {
    setIsModalVisible(true);
  };

  const handleModalClose = () => {
    setIsModalVisible(false);
  };

  return (
    <>
      {/* CSS Keyframes for animations - only needed for severity 1-4 */}
      {severity > 0 && (
        <style>{`
          @keyframes pulse-soft {
            0%, 100% {
              border-color: rgba(255, 0, 0, 0.35);
              box-shadow: 
                0 4px 8px rgba(255, 0, 0, 0.15), 
                0 0 12px rgba(255, 0, 0, 0.1);
              background: linear-gradient(145deg, #fefefe, #fdfbfb);
            }
            50% {
              border-color: rgba(255, 0, 0, 0.5);
              box-shadow: 
                0 4px 10px rgba(255, 0, 0, 0.2), 
                0 0 16px rgba(255, 0, 0, 0.15);
              background: linear-gradient(145deg, #fefcfc, #fdf9f9);
            }
          }

          @keyframes pulse-medium {
            0%, 100% {
              border-color: rgba(255, 0, 0, 0.55);
              box-shadow: 
                0 4px 12px rgba(255, 0, 0, 0.25), 
                0 0 16px rgba(255, 0, 0, 0.15);
              background: linear-gradient(145deg, #fefcfc, #fcf8f8);
            }
            50% {
              border-color: rgba(255, 0, 0, 0.75);
              box-shadow: 
                0 4px 16px rgba(255, 0, 0, 0.35), 
                0 0 24px rgba(255, 0, 0, 0.25);
              background: linear-gradient(145deg, #fef9f9, #fcf5f5);
            }
          }

          @keyframes pulse-strong {
            0%, 100% {
              border-color: rgba(255, 0, 0, 0.8);
              box-shadow: 
                0 4px 16px rgba(255, 0, 0, 0.4), 
                0 0 24px rgba(255, 0, 0, 0.25);
              background: linear-gradient(145deg, #fef9f9, #fcf4f4);
            }
            50% {
              border-color: rgba(255, 0, 0, 1);
              box-shadow: 
                0 4px 20px rgba(255, 0, 0, 0.5), 
                0 0 32px rgba(255, 0, 0, 0.35);
              background: linear-gradient(145deg, #fef7f7, #fcf1f1);
            }
          }

          .service-card:hover {
            transform: translateY(-4px);
          }

          .service-card .ant-card-head {
            background: #f0f2f5;
            padding: 16px 15px;
            border-bottom: 1px solid #e8e8e8;
          }

          .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
          }

          .card-title {
            font-size: 1em;
            font-weight: 600;
            color: #1f2937;
          }

          .tag-container {
            display: flex;
            gap: 3px;
          }

          .env-tag, .type-tag {
            padding: 2px 4px;
            border-radius: 5px;
            font-size: 0.75em;
            font-weight: 500;
          }

          .stat-row {
            margin-top: 8px;
          }

          .stat-item {
            text-align: center;
          }

          .stat-title {
            font-size: 14px;
            color: #6b7280;
          }

          .stat-item .ant-statistic-content {
            font-size: 18px;
            font-weight: 600;
            color: #1f2937;
          }

          .progress-bar {
            margin-top: 8px;
            width: 100% !important;
          }
        `}</style>
      )}

      <Card
        className="service-card"
        title={
          <div className="card-header">
            <span className="card-title">{service_name}</span>
            <div className="tag-container">
              <Tag color={envColor} className="env-tag">
                {server_name}
              </Tag>
              <Tag color={typeColor} className="type-tag">
                {callType}
              </Tag>
            </div>
          </div>
        }
        bordered={false}
        style={cardStyle}
        onClick={handleCardClick} // Add click handler
      >
        {/* Card content with higher z-index for severity > 0 */}
        <div
          style={{ position: "relative", zIndex: severity > 0 ? 2 : "auto" }}
        >
          <Row gutter={[16, 16]} className="stat-row">
            <Col span={12}>
              <Statistic
                title={<span className="stat-title">Uptime</span>}
                value={`${uptime}%`}
                className="stat-item"
              />
            </Col>
            <Col span={12}>
              <Statistic
                title={<span className="stat-title">Avg Latency</span>}
                value={`${avg_latency} ms`}
                className="stat-item"
              />
            </Col>
            <Progress
              percent={uptime}
              strokeColor={progressColor}
              showInfo={false}
              strokeWidth={8}
              className="progress-bar"
            />
          </Row>
          <Row gutter={[16, 16]} className="stat-row">
            <Col span={8}>
              <Statistic
                title={<span className="stat-title">Total Calls</span>}
                value={total_checks}
                className="stat-item"
              />
            </Col>
            <Col span={8}>
              <Statistic
                title={<span className="stat-title">Successes</span>}
                value={successful_checks}
                className="stat-item"
              />
            </Col>
            <Col span={8}>
              <Statistic
                title={<span className="stat-title">Errors</span>}
                value={downtime_count}
                valueStyle={{ color: "#ff4d4f" }}
                className="stat-item"
              />
            </Col>
          </Row>
        </div>
      </Card>

      {/* Service Details Modal */}
      <ServiceDetailsModal
        visible={isModalVisible}
        onClose={handleModalClose}
        id={id}
      />
    </>
  );
};

export default ServiceCard;
