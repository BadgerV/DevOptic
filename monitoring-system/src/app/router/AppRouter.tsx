import React, { lazy, Suspense, useEffect, useState } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import RoleGuard from "../../shared/components/PrivateRoute/PrivateRoute";
import Unauthorized from "../../pages/Unauthorized/Unauthorized";
import RolesPage from "../../domains/rbac/pages/RolePage/RolePage";
import Dashboard from "@/pages/Dashboard/Dashboard";
import LoginPage from "@/pages/Login/Login";
import RBACManagementPage from "../../domains/rbac/pages/RolePage/RolePage";
import RegisterPage from "@/pages/Register/Register";
import { Spin } from "antd";
import { FloatingOrb } from "@/shared/components/floating-orb/FloatingOrb";
import { Sidebar } from "@/shared/components/rightSidebar/RightSidebar";
import {
  setUserState,
  useValidateTokenQuery,
} from "@/shared/services/api/authApi";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../store";
import SettingsPage from "@/pages/Settings/SettingsPage";

// Mock PrivateRoute component - replace with your actual PrivateRoute implementation
interface PrivateRouteProps {
  children: React.ReactNode;
}

const LazyGitlabRoutes = lazy(() =>
  import("@/domains/gitlab/index").then((module) => ({
    default: module.GitlabRoutes
  }))
);

interface UserState {
  email: string;
  username: string;
  is_active: boolean;
  userId: string;
  isAdmin: boolean;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const dispatch = useDispatch();
  // If no token at all, block access
  // if (!token) {
  //   return <Navigate to="/login" replace />;
  // }

  // Run validateToken query
  const { data, error, isLoading, isError } = useValidateTokenQuery();

  const isAdmin = useSelector((state: RootState) => state.auth);

  useEffect(() => {
    if (data) {
      const newData: any = {
        email: data.user.User.email,
        username: data.user.User.username,
        is_active: data.user.User.is_active,
        userId: data.user.User.id,
        isAdmin: data.user.IsAdmin,
      };

      dispatch(setUserState(newData));
    }
  }, [data]);

  if (isLoading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <Spin size="large" tip="Validating session..." />
      </div>
    );
  }

  if (isError || !data?.user.User.id) {
    // message.error("Your session has expired, please log in again.");
    localStorage.removeItem("token");
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

const AppRoutes: React.FC = () => {
  return (
    <Routes>
      {/* Public routes */}
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/unauthorized" element={<Unauthorized />} />

      {/* Protected routes */}
      <Route
        path="/"
        element={
          <PrivateRoute>
            <Dashboard />
          </PrivateRoute>
        }

      />
      <Route
        path="/settings"
        element={
          <PrivateRoute>
            <SettingsPage />
          </PrivateRoute>
        }
      />

      {/* RBAC protected routes - requires super admin permissions */}
      <Route
        path="/rbac/roles"
        element={
          <PrivateRoute>
            <RoleGuard>
              <RBACManagementPage />
            </RoleGuard>
          </PrivateRoute>
        }
      />

      {/* Additional RBAC routes can be added here */}
      <Route
        path="/rbac/*"
        element={
          <PrivateRoute>
            <RoleGuard>
              <div style={{ padding: "24px" }}>
                <h2>RBAC Module</h2>
                <p>More RBAC features coming soon...</p>
                <ul>
                  <li>
                    <a href="/rbac/roles">Manage Roles</a>
                  </li>
                  <li>Manage Permissions (Coming Soon)</li>
                  <li>User Role Assignments (Coming Soon)</li>
                </ul>
              </div>
            </RoleGuard>
          </PrivateRoute>
        }
      />

      <Route
        path="/gitlab/*"
        element={
          <PrivateRoute>
            {/* <RoleGuard roles={["devops", "admin"]}> */}
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
              <LazyGitlabRoutes />
            </Suspense>
            {/* </RoleGuard> */}
          </PrivateRoute>
        }
      />

      {/* Catch-all route - redirect to home */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
};

// Types
interface SidebarLink {
  label: string;
  url: string;
}

interface FloatingOrbProps {
  onClick: () => void;
  isOpen: boolean;
}

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
  links: SidebarLink[];
}

// AppLayout Component (wraps your router with floating components)
export const AppLayout: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const user = useSelector((state: RootState) => state.auth.user);

  const sidebarLinks: SidebarLink[] = [
    { label: "Monitoring Module", url: "/monitoring" },
    { label: "GitLab", url: "/gitlab" },
    { label: "Linux", url: "/linux" },
    { label: "Docker", url: "/docker" },
    { label: "Kubernetes", url: "/kubernetes" },
    { label: "RBAC", url: "/rbac/roles" },
  ];

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const closeSidebar = () => {
    setSidebarOpen(false);
  };

  return (
    <div className="app">
      {/* Your existing router */}
      <AppRoutes />

      {user && <FloatingOrb onClick={toggleSidebar} isOpen={sidebarOpen} />}

      {/* Floating Orb */}

      {/* Sidebar */}

      {user && (
        <Sidebar
          isOpen={sidebarOpen}
          onClose={closeSidebar}
          links={sidebarLinks}
        />
      )}
    </div>
  );
};

export default AppRoutes;
