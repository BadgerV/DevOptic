import { Layout, Result, Spin } from "antd";
import PipelineUnitForm from "../../components/PipelineUnitForm/PipelineUnitForm";
import {  useGetServicesQuery } from "../../store";
import { useCheckSuperAdminQuery } from "@/domains/rbac/store/rbacApi";
import { useEffect } from "react";
import { useSelector } from "react-redux";
import { RootState } from "@/app/store";

const PipelineUnitCreationPage = () => {
  const { data: services, isLoading, error } = useGetServicesQuery();
  const userId =
    useSelector((state: RootState) => state.auth.user?.userId) || "";
  const {
    data: superAdminData,
    errorPlus,
    isError,
  } = useCheckSuperAdminQuery(userId);

  useEffect(() => {
    console.log(superAdminData);
  }, [superAdminData]);

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
          subTitle="You need super admin privileges to create pipeline units."
        />
      </Layout.Content>
    );
  }

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

  if (error) {
    return (
      <Layout.Content>
        <Result
          status="error"
          title="Error"
          subTitle="Failed to load services. Please try again."
        />
      </Layout.Content>
    );
  }

  return (
    <Layout.Content style={{ padding: "24px" }}>
      <PipelineUnitForm services={services.services || []} />
    </Layout.Content>
  );
};

export default PipelineUnitCreationPage;
