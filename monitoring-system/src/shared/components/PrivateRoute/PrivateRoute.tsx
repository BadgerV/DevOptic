import React, { ReactNode, useEffect } from "react";
import { Navigate } from "react-router-dom";
import { Spin } from "antd";
import { useCheckSuperAdminQuery } from "../../../domains/rbac/store/rbacApi";
import { message } from "antd";
import { useSelector } from "react-redux";
import { RootState } from "@/app/store";

interface RoleGuardProps {
  children: ReactNode;
}

const RoleGuard: React.FC<RoleGuardProps> = ({ children }) => {
  // Check if user has a token in localStorage
  const userId = useSelector((state: RootState) => state.auth.user?.userId) || "";

  // Skip the super admin check if no token is present
  const {
    data: superAdminData,
    error,
    isLoading,
    isError,
  } = useCheckSuperAdminQuery(userId);

  // If no token, redirect to unauthorized
  if (!userId) {
    return <Navigate to="/unauthorized" replace />;
  } else {
  }

  useEffect(() => {
    if (isError && error) {
      let errorMessage = "Failed to verify permissions";

      if (
        typeof error === "object" &&
        error !== null &&
        "data" in error &&
        typeof (error as any).data === "object" &&
        (error as any).data !== null &&
        "message" in (error as any).data
      ) {
        errorMessage = (error as any).data.message;
      }

      message.error(errorMessage);
    }
  }, [isError, error]);

  // Show loading spinner while checking permissions
  if (isLoading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "200px",
        }}
      >
        <Spin size="large" tip="Verifying permissions..." />
      </div>
    );
  }

  // If there's an error or user is not a super admin, redirect to unauthorized
  if (isError || !superAdminData?.is_super_admin) {
    return <Navigate to="/unauthorized" replace />;
  }

  // If user is a super admin, render the protected content
  return <>{children}</>;
};

export default RoleGuard;
