import { useSelector } from "react-redux";
import { Layout, Result, Spin } from "antd";
import { ServiceForm } from "../../components/ServiceForm";
import { RootState } from "@/app/store";
import { useCheckSuperAdminQuery } from "@/domains/rbac/store/rbacApi";

const ServiceRegistrationPage = () => {
  const userId = useSelector((state: RootState) => state.auth.user?.userId) || "";
  const {
    data: superAdminData,
    error,
    isLoading,
    isError,
  } = useCheckSuperAdminQuery(userId);
  if (isLoading) {
    return (
      <Layout.Content
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <Spin size="large" />
      </Layout.Content>
    );
  }

  if (error || !superAdminData?.is_super_admin) {
    return (
      <Layout.Content>
        <Result
          status="403"
          title="Access Denied"
          subTitle="You need super admin privileges to register a service."
        />
      </Layout.Content>
    );
  }

  return (
    <Layout.Content style={{ padding: "24px" }}>
      <ServiceForm />
    </Layout.Content>
  );
};

export default ServiceRegistrationPage;
