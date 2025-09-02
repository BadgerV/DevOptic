import React, { useEffect, useState } from "react";
import "./endpointMainRight.css";
import {
  Card,
  Divider,
  Timeline,
  Statistic,
  Progress,
  List,
  Typography,
} from "antd";
import { getOverallStats } from "../../services/endpointMonitoringApiCalls";
import { OverallStats } from "../../types";
import { useAppDispatch } from "@/shared/hooks/redux";
import { setTotalNumberOfEndpoints } from "../../store/endpointMonitoringSlice";

const { Title, Text } = Typography;

const EndpointMainRight = () => {
  const [overview, setOverview] = useState<OverallStats | null>(null);
  const [loading, setLoading] = useState(false);

  const dispatch = useAppDispatch();

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        const data = await getOverallStats();
        if (typeof data === "object" && data !== null && "data" in data) {
          const mainData = (data as { data: OverallStats }).data;
          setOverview(mainData);
          dispatch(setTotalNumberOfEndpoints(mainData.total_endpoints));
        } else {
          throw new Error("Unexpected response format");
        }
      } catch (err) {
        console.error("Failed to fetch overall stats:", err);
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, []);

  // Sample log data
  const liveLogs = [
    { timestamp: "2024-04-24 13:45:01", message: "Server started" },
    { timestamp: "2024-04-24 13:45:15", message: "Incoming request processed" },
  ];

  // Sample incident data
  const incidents = [
    { timestamp: "2024-04-24 13:42:11", message: "Timeout" },
    { timestamp: "2024-04-24 13:41:55", message: "Connection lost" },
    { timestamp: "2024-04-24 13:41:45", message: "Request failed" },
  ];

  if (!overview) {
    return (
      <div className="right-sidebar">
        <Card className="sidebar-card" bordered={false}>
          <Title level={4} className="section-title">
            Dashboard Overview
          </Title>
          <p>Loading stats...</p>
        </Card>
      </div>
    );
  }

  // Calculate success rate
  const successRate = overview.total_checks
    ? Number(
        ((overview.successful_checks / overview.total_checks) * 100).toFixed(1)
      )
    : 0;

  // Format uptime and latency for display
  const formattedUptime = Number(overview.overall_uptime.toFixed(1));
  const formattedLatency = Number(overview.average_latency.toFixed(1));

  return (
    <div className="right-sidebar">
      <Card className="sidebar-card" bordered={false}>
        <Title level={4} className="section-title">
          Dashboard Overview
        </Title>

        {/* Live Logs */}
        {/* <Divider orientation="left">Live Logs</Divider>
        <Timeline className="logs-timeline">
          {liveLogs.map((log, index) => (
            <Timeline.Item
              key={index}
              dot={<ClockCircleOutlined style={{ color: "#1890ff" }} />}
            >
              <Text className="log-text">
                <span className="log-timestamp">{log.timestamp}</span> -{" "}
                {log.message}
              </Text>
            </Timeline.Item>
          ))}
        </Timeline> */}

        {/* Historical Uptime */}
        <Divider orientation="left">Historical Uptime</Divider>
        <div className="uptime-stats">
          <Statistic title="Overall Uptime" value={`${formattedUptime}%`} />
          <Progress
            percent={formattedUptime}
            strokeColor={
              formattedUptime > 95
                ? "#52c41a"
                : formattedUptime > 80
                ? "#fa8c16"
                : "#ff4d4f"
            }
            showInfo={false}
            strokeWidth={8}
            className="uptime-progress"
          />
          <div className="uptime-details">
            <Statistic title="Total Outages" value={overview.down_time_count} />
            <Statistic title="Max Uptime" value="100%" />
          </div>
        </div>

        {/* Performance Metrics */}
        <Divider orientation="left">Performance Metrics</Divider>
        <div className="performance-stats">
          <Statistic title="Total Endpoints" value={overview.total_endpoints} />
          <Statistic title="Total Checks" value={overview.total_checks} />
          <Statistic title="Success Rate" value={`${successRate}%`} />
          <Statistic title="Average Latency" value={`${formattedLatency} ms`} />
        </div>

        {/* Recent Incidents */}
        {/* <Divider orientation="left">Recent Incidents</Divider>
        <List
          className="incidents-list"
          dataSource={incidents}
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                avatar={<WarningOutlined style={{ color: "#ff4d4f" }} />}
                title={
                  <Text className="incident-timestamp">{item.timestamp}</Text>
                }
                description={
                  <Text className="incident-message">{item.message}</Text>
                }
              />
            </List.Item>
          )}
        /> */}
      </Card>
    </div>
  );
};

export default EndpointMainRight;
