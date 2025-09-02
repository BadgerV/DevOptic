import React, { useState } from "react";
import { Spin } from "antd";
// import { useGetAllEndpointStatusesQuery } from "@domains/endpoint-monitoring/store/endpointMonitoringApi";

import "./dashboard.css";
import Sidebar from "@/shared/components/sidebar/Sidebar";
import EndpointMain from "@/domains/endpoint-monitoring/components/endpointMain/EndpointMain";

const Dashboard: React.FC = () => {
  // const {
  //   data: statusData,
  //   isLoading,
  //   error,
  // } = useGetAllEndpointStatusesQuery();

  // if (isLoading) {
  //   return (
  //     <div>
  //       <Spin size="large" />
  //     </div>
  //   );
  // }

  // if (error) {
  //   return <div>Error loading dashboard data</div>;
  // }

  const [isSidebarOpen, setSidebarOpen] = useState(true);

  const toggleSidebar = () => {
    setSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="dashboard">
      <Sidebar isOpen={isSidebarOpen} onClose={toggleSidebar} />
      <EndpointMain isSidebarOpen={isSidebarOpen} />
    </div>
  );
};

export default Dashboard;
