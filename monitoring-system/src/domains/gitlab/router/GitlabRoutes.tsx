import { Routes, Route, Navigate, useLocation, useNavigate } from "react-router-dom";
import { Suspense, lazy } from "react";
import { Layout, Menu, Result, Spin } from "antd";
import { useDispatch, useSelector } from "react-redux";
import { toggleSidebar } from "@/domains/gitlab/store";
import { RootState } from "@/app/store";

import {
  AppstoreOutlined,
  HistoryOutlined,
  CheckCircleOutlined,
  ApartmentOutlined,
  PlayCircleOutlined,
} from "@ant-design/icons";

// Lazy-loaded page components
const ServiceRegistrationPage = lazy(
  () => import("../pages/serviceRegistrationPage/ServiceRegistrationPage")
);
const PipelineUnitCreationPage = lazy(
  () => import("../pages/pipelineUnitCreationPage/PipelineUnitCreationPage")
);
const PipelinesListPage = lazy(
  () => import("../pages/pipelinesListPage/PipelinesListPage")
);
const ExecutionHistoryPage = lazy(
  () => import("../pages/executionHistoryPage/ExecutionHistoryPage")
);
const ApprovalsPage = lazy(
  () => import("../pages/approvalsPage/ApprovalsPage")
);

const GitlabSidebar = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const isSidebarOpen = useSelector(
    (state: RootState) => state.gitlab.ui.isSidebarOpen
  );
  const location = useLocation();

  const menuItems = [
    {
      key: "/gitlab/services",
      icon: <AppstoreOutlined />,
      label: "Services",
    },
    {
      key: "/gitlab/pipeline-units",
      icon: <ApartmentOutlined />,
      label: "Pipeline Units",
    },
    {
      key: "/gitlab/pipelines",
      icon: <PlayCircleOutlined />,
      label: "Pipelines",
    },
    {
      key: "/gitlab/history",
      icon: <HistoryOutlined />,
      label: "Execution History",
    },
    {
      key: "/gitlab/approvals",
      icon: <CheckCircleOutlined />,
      label: "Approvals",
    },
  ];

  return (
    <Layout.Sider
      width={200}
      theme="light"
      collapsible
      collapsed={!isSidebarOpen}
      onCollapse={() => dispatch(toggleSidebar())}
      style={{ height: "100vh", position: "fixed", left: 0 }}
    >
      <div style={{ padding: "16px", textAlign: "center" }}>
        <h3
          style={{ margin: 0, color: isSidebarOpen ? "#000" : "transparent" }}
        >
          GitLab Automation
        </h3>
      </div>
      <Menu
        mode="inline"
        selectedKeys={[location.pathname]}
        items={menuItems}
        style={{ borderRight: 0 }}
        onClick={({ key }) => navigate(key)}
      />
    </Layout.Sider>
  );
};

export const GitlabRoutes = () => {
  const isSidebarOpen = useSelector(
    (state: RootState) => state.gitlab.ui.isSidebarOpen
  );

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <GitlabSidebar />
      <Layout style={{ marginLeft: isSidebarOpen ? 200 : 80 }}>
        <Layout.Content style={{ padding: "24px" }}>
          <Suspense
            fallback={
              <Spin
                size="large"
                style={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "center",
                  height: "100vh",
                }}
              />
            }
          >
            <Routes>
              <Route path="services" element={<ServiceRegistrationPage />} />
              <Route
                path="pipeline-units"
                element={<PipelineUnitCreationPage />}
              />
              <Route path="pipelines" element={<PipelinesListPage />} />
              <Route path="history" element={<ExecutionHistoryPage />} />
              <Route path="approvals" element={<ApprovalsPage />} />
              <Route index element={<Navigate to="pipelines" replace />} />
              <Route
                path="*"
                element={
                  <Result status="404" title="404" subTitle="Page not found" />
                }
              />
            </Routes>
          </Suspense>
        </Layout.Content>
      </Layout>
    </Layout>
  );
};
