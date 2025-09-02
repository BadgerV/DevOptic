import React, { useEffect } from "react";
import "./endpointMainLeft.css";
import { useGetAllEndpointEssentialsQuery } from "../../store/endpointMonitoringApi";
import { EndpointsEssentials } from "../../types";
import ServiceCard from "../serviceCard/ServiceCard";
import { useAppSelector, useAppDispatch } from "@/shared/hooks/redux";
import { setRunningEndpoints } from "../../store/endpointMonitoringSlice";
import { Spin, Alert, Empty } from "antd";

type MonitoringDashboardProps = {
  filter: string;
};

const EndpointMainLeft: React.FC<MonitoringDashboardProps> = ({ filter }) => {
  const {
    data: endpointEssentials,
    isLoading,
    error,
  } = useGetAllEndpointEssentialsQuery();

  const filterByRunning = useAppSelector(
    (state) => state.endpointMonitoring.filterTypeByRunning
  );

  const dispatch = useAppDispatch();

  // ✅ useEffect to set runningEndpoints from ALL endpoints
  useEffect(() => {
    if (endpointEssentials?.data) {
      const runningCount = endpointEssentials.data.filter(
        (endpoint: EndpointsEssentials) => endpoint.last_run === true
      ).length;

      dispatch(setRunningEndpoints(runningCount));
    }
  }, [endpointEssentials, dispatch]);

  // normalize filter for case-insensitive matching
  const normalizedFilter = filter.toLowerCase();

  const filteredEndpoints = endpointEssentials?.data
    ?.filter((endpoint: EndpointsEssentials) => {
      if (filterByRunning === "all") return true;
      if (filterByRunning === true) return endpoint.last_run === true;
      if (filterByRunning === false) return endpoint.last_run === false;
      return true;
    })
    ?.filter((endpoint: EndpointsEssentials) =>
      endpoint.service_name.toLowerCase().includes(normalizedFilter)
    );

  // ✅ Conditional rendering AFTER hooks
  if (isLoading)
    return (
      <div className="endpoint-main-left-isloading">
        <Spin tip="Loading endpoints..." />
      </div>
    );
  if (error)
    return (
      <div className="endpoint-main-left-isloading">
        <Alert
          type="error"
          message="Failed to load endpoints"
          description="Please try refreshing the page."
          showIcon
        />
      </div>
    );

  return (
    <div className="endpoint-main-left">
      {filteredEndpoints?.length ? (
        filteredEndpoints.map((endpoint: EndpointsEssentials) => (
          <ServiceCard
            key={endpoint.endpointID}
            service_name={endpoint.service_name}
            server_name={endpoint.server_name}
            uptime={+endpoint.uptime_percentage}
            callType="API"
            total_checks={+endpoint.total_checks}
            successful_checks={+endpoint.successful_checks}
            downtime_count={+endpoint.downtime_count}
            avg_latency={+endpoint.avg_latency}
            severity={+endpoint.failure_count}
            id={endpoint.id}
          />
        ))
      ) : (
        <Empty description={`No services found matching "${filter}"`} />
      )}
    </div>
  );
};

export default EndpointMainLeft;
