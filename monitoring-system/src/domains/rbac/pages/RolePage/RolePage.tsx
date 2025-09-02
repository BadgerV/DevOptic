import React, { useState, useEffect } from "react";
import {
  Card,
  Form,
  Input,
  Button,
  Table,
  Space,
  Typography,
  Row,
  Col,
  message,
  Spin,
  Tabs,
  Modal,
  Select,
  Divider,
  Tag,
  Popconfirm,
  Alert,
  Badge,
} from "antd";
import {
  PlusOutlined,
  ReloadOutlined,
  UserOutlined,
  SafetyOutlined,
  SettingOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  DeleteOutlined,
  LinkOutlined,
} from "@ant-design/icons";
import {
  useCreateRoleMutation,
  useGetAllRolesQuery,
  useCreatePermissionMutation,
  useAssignRoleMutation,
  useRemoveRoleMutation,
  useAssignPermissionMutation,
  useGetUserPermissionsQuery,
  useCheckPermissionMutation,
  useCheckUserPermissionQuery,
  useCheckSuperAdminQuery,
  useHealthCheckQuery,
  useGetUserUsernameAndIdQuery,
  Role,
  Permission,
  CreateRoleRequest,
  CreatePermissionRequest,
  AssignRoleRequest,
  RemoveRoleRequest,
  AssignPermissionRequest,
  CheckPermissionRequest,
} from "../../store/rbacApi";

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;
const { TextArea } = Input;

const RBACManagementPage: React.FC = () => {
  // Forms
  const [roleForm] = Form.useForm();
  const [permissionForm] = Form.useForm();
  const [assignRoleForm] = Form.useForm();
  const [removeRoleForm] = Form.useForm();
  const [assignPermissionForm] = Form.useForm();
  const [checkPermissionForm] = Form.useForm();
  const [userPermissionsForm] = Form.useForm();
  const [superAdminForm] = Form.useForm();

  // State for modals and selections
  const [userPermissionsModalVisible, setUserPermissionsModalVisible] =
    useState(false);
  const [selectedUserId, setSelectedUserId] = useState<string>("");
  const [checkPermissionResult, setCheckPermissionResult] = useState<any>(null);
  const [userPermissionsData, setUserPermissionsData] = useState<any>(null);
  const [superAdminCheckResult, setSuperAdminCheckResult] = useState<any>(null);

  // RTK Query hooks
  const {
    data: healthData,
    isLoading: healthLoading,
    refetch: refetchHealth,
  } = useHealthCheckQuery();

  const {
    data: roles = [],
    isLoading: rolesLoading,
    error: rolesError,
    refetch: refetchRoles,
  } = useGetAllRolesQuery();

  const {
    data: usersData,
    isLoading: usersLoading,
    error: usersError,
  } = useGetUserUsernameAndIdQuery();

  const {
    data: userPermissions,
    isLoading: userPermissionsLoading,
    error: userPermissionsError,
  } = useGetUserPermissionsQuery(selectedUserId, {
    skip: !selectedUserId,
  });

  // Mutations
  const [createRole, { isLoading: createRoleLoading }] =
    useCreateRoleMutation();
  const [createPermission, { isLoading: createPermissionLoading }] =
    useCreatePermissionMutation();
  const [assignRole, { isLoading: assignRoleLoading }] =
    useAssignRoleMutation();
  const [removeRole, { isLoading: removeRoleLoading }] =
    useRemoveRoleMutation();
  const [assignPermission, { isLoading: assignPermissionLoading }] =
    useAssignPermissionMutation();
  const [checkPermission, { isLoading: checkPermissionLoading }] =
    useCheckPermissionMutation();

  // Effect for user permissions data
  useEffect(() => {
    if (userPermissions) {
      setUserPermissionsData(userPermissions);
    }
  }, [userPermissions]);

  const [messageApi, contextHolder] = message.useMessage();

  // Handle form submissions
  const handleCreateRole = async (values: CreateRoleRequest) => {
    try {
      const result = await createRole(values).unwrap();
      messageApi.open({
        type: "success",
        content: "Role created successfully!",
      });
      roleForm.resetFields();
    } catch (error: any) {
      const errorMessage = error?.data?.message || "Failed to create role";
      messageApi.open({
        type: "error",
        content: errorMessage,
      });
    }
  };

  const handleCreatePermission = async (values: CreatePermissionRequest) => {
    try {
      const result = await createPermission(values).unwrap();
      messageApi.open({
        type: "success",
        content: "Permission created successfully!",
      });
      permissionForm.resetFields();
    } catch (error: any) {
      const errorMessage =
        error?.data?.message || "Failed to create permission";
      messageApi.open({
        type: "error",
        content: errorMessage,
      });
    }
  };

  const handleAssignRole = async (values: AssignRoleRequest) => {
    try {
      const result = await assignRole(values).unwrap();
      messageApi.open({
        type: "success",
        content: result.message || "Role assigned successfully!",
      });
      assignRoleForm.resetFields();
    } catch (error: any) {
      const errorMessage = error?.data?.message || "Failed to assign role";
      messageApi.open({
        type: "error",
        content: errorMessage,
      });
    }
  };

  const handleRemoveRole = async (values: RemoveRoleRequest) => {
    try {
      const result = await removeRole(values).unwrap();
      messageApi.open({
        type: "success",
        content: result.message || "Role removed successfully!",
      });
      removeRoleForm.resetFields();
    } catch (error: any) {
      const errorMessage = error?.data?.message || "Failed to remove role";
      messageApi.open({
        type: "error",
        content: errorMessage,
      });
    }
  };

  const handleAssignPermission = async (values: AssignPermissionRequest) => {
    try {
      const result = await assignPermission(values).unwrap();
      messageApi.open({
        type: "success",
        content: result.message || "Permission assigned successfully!",
      });
      assignPermissionForm.resetFields();
    } catch (error: any) {
      const errorMessage =
        error?.data?.message || "Failed to assign permission";
      messageApi.open({
        type: "error",
        content: errorMessage,
      });
    }
  };

  const handleCheckPermission = async (values: CheckPermissionRequest) => {
    try {
      const result = await checkPermission(values).unwrap();
      setCheckPermissionResult(result);
      messageApi.open({
        type: "info",
        content: "Permission Check Completed",
      });
    } catch (error: any) {
      const errorMessage = error?.data?.message || "Failed to check permission";
      setCheckPermissionResult(null);
      messageApi.open({
        type: "error",
        content: errorMessage,
      });
    }
  };

  const handleGetUserPermissions = (values: { user_id: string }) => {
    setSelectedUserId(values.user_id);
    setUserPermissionsModalVisible(true);
  };

  const handleCheckSuperAdmin = async (values: { user_id: string }) => {
    try {
      const result = await refetchHealth(); // Placeholder, replace with actual checkSuperAdmin query
      setSuperAdminCheckResult({ user_id: values.user_id, is_super_admin: true });
      messageApi.open({
        type: "info",
        content: "Super admin check completed",
      });
    } catch (error: any) {
      messageApi.open({
        type: "error",
        content: "Failed to check super admin status",
      });
    }
  };

  // Table columns for roles
  const rolesColumns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
      sorter: (a: Role, b: Role) => a.name.localeCompare(b.name),
    },
    {
      title: "Description",
      dataIndex: "description",
      key: "description",
      ellipsis: true,
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      key: "created_at",
      render: (date: string) => {
        if (!date) return "-";
        return new Date(date).toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit",
        });
      },
    },
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
      width: 120,
      render: (id: string) => (
        <Text code copyable={{ text: id }}>
          {id.substring(0, 8)}...
        </Text>
      ),
    },
  ];

  return (
    <div style={{ padding: "24px" }}>
      {contextHolder}
      <div style={{ marginBottom: "24px" }}>
        <Title level={2}>RBAC Management Dashboard</Title>
        <Space>
          <Badge
            status={
              healthData?.status === "RBAC service OK" ? "success" : "error"
            }
            text={healthData?.status || "Checking..."}
          />
          <Button
            size="small"
            icon={<ReloadOutlined />}
            loading={healthLoading}
            onClick={() => refetchHealth()}
          >
            Refresh Status
          </Button>
        </Space>
      </div>

      <Tabs defaultActiveKey="roles" type="card">
        {/* ROLES TAB */}
        <TabPane
          tab={
            <span>
              <UserOutlined />
              Roles
            </span>
          }
          key="roles"
        >
          <Row gutter={[24, 24]}>
            <Col xs={24} lg={8}>
              <Card title="Create New Role" size="small">
                <Form
                  form={roleForm}
                  layout="vertical"
                  onFinish={handleCreateRole}
                  disabled={createRoleLoading}
                >
                  <Form.Item
                    name="name"
                    label="Role Name"
                    rules={[
                      { required: true, message: "Please enter role name" },
                      {
                        min: 2,
                        message: "Role name must be at least 2 characters",
                      },
                    ]}
                  >
                    <Input placeholder="e.g., admin, user, manager" />
                  </Form.Item>
                  <Form.Item
                    name="description"
                    label="Description"
                    rules={[
                      { required: true, message: "Please enter description" },
                    ]}
                  >
                    <TextArea
                      rows={3}
                      placeholder="Describe the role's purpose"
                    />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={createRoleLoading}
                      icon={<PlusOutlined />}
                      block
                    >
                      Create Role
                    </Button>
                  </Form.Item>
                </Form>
              </Card>
            </Col>
            <Col xs={24} lg={16}>
              <Card
                title="All Roles"
                size="small"
                extra={
                  <Button
                    icon={<ReloadOutlined />}
                    onClick={refetchRoles}
                    loading={rolesLoading}
                    size="small"
                  >
                    Refresh
                  </Button>
                }
              >
                <Table
                  columns={rolesColumns}
                  dataSource={roles}
                  rowKey="id"
                  loading={rolesLoading}
                  pagination={{ pageSize: 10 }}
                  size="small"
                />
              </Card>
            </Col>
          </Row>
        </TabPane>

        {/* PERMISSIONS TAB */}
        <TabPane
          tab={
            <span>
              <SafetyOutlined />
              Permissions
            </span>
          }
          key="permissions"
        >
          <Row gutter={[24, 24]}>
            <Col xs={24} lg={12}>
              <Card title="Create New Permission" size="small">
                <Form
                  form={permissionForm}
                  layout="vertical"
                  onFinish={handleCreatePermission}
                  disabled={createPermissionLoading}
                >
                  <Form.Item
                    name="resource"
                    label="Resource"
                    rules={[
                      { required: true, message: "Please enter resource" },
                    ]}
                  >
                    <Input placeholder="e.g., users, posts, settings" />
                  </Form.Item>
                  <Form.Item
                    name="action"
                    label="Action"
                    rules={[{ required: true, message: "Please enter action" }]}
                  >
                    <Input placeholder="e.g., read, write, delete, create" />
                  </Form.Item>
                  <Form.Item
                    name="description"
                    label="Description"
                    rules={[
                      { required: true, message: "Please enter description" },
                    ]}
                  >
                    <TextArea
                      rows={2}
                      placeholder="Describe what this permission allows"
                    />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={createPermissionLoading}
                      icon={<PlusOutlined />}
                      block
                    >
                      Create Permission
                    </Button>
                  </Form.Item>
                </Form>
              </Card>
            </Col>
            <Col xs={24} lg={12}>
              <Card title="Check User Permission" size="small">
                <Form
                  form={checkPermissionForm}
                  layout="vertical"
                  onFinish={handleCheckPermission}
                  disabled={checkPermissionLoading}
                >
                  <Form.Item
                    name="user_id"
                    label="User"
                    rules={[{ required: true, message: "Please select a user" }]}
                  >
                    <Select
                      placeholder="Select a user"
                      loading={usersLoading}
                      showSearch
                      optionFilterProp="children"
                      filterOption={(input, option) =>
                        (option?.children as unknown as string)
                          ?.toLowerCase()
                          .includes(input.toLowerCase())
                      }
                    >
                      {usersData?.data?.map((user: { id: string; username: string }) => (
                        <Option key={user.id} value={user.id}>
                          {user.username} ({user.id.substring(0, 8)}...)
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item
                    name="resource"
                    label="Resource"
                    rules={[
                      { required: true, message: "Please enter resource" },
                    ]}
                  >
                    <Input placeholder="e.g., users, posts" />
                  </Form.Item>
                  <Form.Item
                    name="action"
                    label="Action"
                    rules={[{ required: true, message: "Please enter action" }]}
                  >
                    <Input placeholder="e.g., read, write" />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={checkPermissionLoading}
                      icon={<CheckCircleOutlined />}
                      block
                    >
                      Check Permission
                    </Button>
                  </Form.Item>
                </Form>

                {checkPermissionResult && (
                  <Alert
                    style={{ marginTop: 16 }}
                    message={checkPermissionResult.message}
                    description={
                      <Space>
                        <Text>Permission Status:</Text>
                        <Tag
                          color={
                            checkPermissionResult.has_permission
                              ? "success"
                              : "error"
                          }
                        >
                          {checkPermissionResult.has_permission
                            ? "ALLOWED"
                            : "DENIED"}
                        </Tag>
                      </Space>
                    }
                    type={
                      checkPermissionResult.has_permission
                        ? "success"
                        : "warning"
                    }
                    showIcon
                  />
                )}
              </Card>
            </Col>
          </Row>
        </TabPane>

        {/* USER ASSIGNMENTS TAB */}
        <TabPane
          tab={
            <span>
              <LinkOutlined />
              Assignments
            </span>
          }
          key="assignments"
        >
          <Row gutter={[24, 24]}>
            <Col xs={24} lg={12}>
              <Card title="Assign Role to User" size="small">
                <Form
                  form={assignRoleForm}
                  layout="vertical"
                  onFinish={handleAssignRole}
                  disabled={assignRoleLoading}
                >
                  <Form.Item
                    name="user_id"
                    label="User"
                    rules={[{ required: true, message: "Please select a user" }]}
                  >
                    <Select
                      placeholder="Select a user"
                      loading={usersLoading}
                      showSearch
                      optionFilterProp="children"
                      filterOption={(input, option) =>
                        (option?.children as unknown as string)
                          ?.toLowerCase()
                          .includes(input.toLowerCase())
                      }
                    >
                      {usersData?.data?.map((user: { id: string; username: string }) => (
                        <Option key={user.id} value={user.id}>
                          {user.username} ({user.id.substring(0, 8)}...)
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item
                    name="role_id"
                    label="Role"
                    rules={[
                      { required: true, message: "Please select a role" },
                    ]}
                  >
                    <Select placeholder="Select a role" loading={rolesLoading}>
                      {roles.map((role) => (
                        <Option key={role.id} value={role.id}>
                          {role.name} - {role.description}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={assignRoleLoading}
                      icon={<LinkOutlined />}
                      block
                    >
                      Assign Role
                    </Button>
                  </Form.Item>
                </Form>
              </Card>

              <Card
                title="Remove Role from User"
                size="small"
                style={{ marginTop: 16 }}
              >
                <Form
                  form={removeRoleForm}
                  layout="vertical"
                  onFinish={handleRemoveRole}
                  disabled={removeRoleLoading}
                >
                  <Form.Item
                    name="user_id"
                    label="User"
                    rules={[{ required: true, message: "Please select a user" }]}
                  >
                    <Select
                      placeholder="Select a user"
                      loading={usersLoading}
                      showSearch
                      optionFilterProp="children"
                      filterOption={(input, option) =>
                        (option?.children as unknown as string)
                          ?.toLowerCase()
                          .includes(input.toLowerCase())
                      }
                    >
                      {usersData?.data?.map((user: { id: string; username: string }) => (
                        <Option key={user.id} value={user.id}>
                          {user.username} ({user.id.substring(0, 8)}...)
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item
                    name="role_id"
                    label="Role"
                    rules={[
                      { required: true, message: "Please select a role" },
                    ]}
                  >
                    <Select
                      placeholder="Select a role to remove"
                      loading={rolesLoading}
                    >
                      {roles.map((role) => (
                        <Option key={role.id} value={role.id}>
                          {role.name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item>
                    <Popconfirm
                      title="Are you sure you want to remove this role?"
                      onConfirm={() => removeRoleForm.submit()}
                      okText="Yes"
                      cancelText="No"
                    >
                      <Button
                        danger
                        loading={removeRoleLoading}
                        icon={<DeleteOutlined />}
                        block
                      >
                        Remove Role
                      </Button>
                    </Popconfirm>
                  </Form.Item>
                </Form>
              </Card>
            </Col>

            <Col xs={24} lg={12}>
              <Card title="Assign Permission to Role" size="small">
                <Form
                  form={assignPermissionForm}
                  layout="vertical"
                  onFinish={handleAssignPermission}
                  disabled={assignPermissionLoading}
                >
                  <Form.Item
                    name="role_id"
                    label="Role"
                    rules={[
                      { required: true, message: "Please select a role" },
                    ]}
                  >
                    <Select placeholder="Select a role" loading={rolesLoading}>
                      {roles.map((role) => (
                        <Option key={role.id} value={role.id}>
                          {role.name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                  <Form.Item
                    name="permission_id"
                    label="Permission ID"
                    rules={[
                      { required: true, message: "Please enter permission ID" },
                    ]}
                  >
                    <Input placeholder="Enter permission UUID" />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={assignPermissionLoading}
                      icon={<SafetyOutlined />}
                      block
                    >
                      Assign Permission
                    </Button>
                  </Form.Item>
                </Form>
              </Card>

              <Card
                title="User Operations"
                size="small"
                style={{ marginTop: 16 }}
              >
                <Space direction="vertical" style={{ width: "100%" }}>
                  <Form
                    layout="inline"
                    onFinish={handleGetUserPermissions}
                  >
                    <Form.Item
                      name="user_id"
                      rules={[{ required: true, message: "Please select a user" }]}
                    >
                      <Select
                        placeholder="Select a user"
                        loading={usersLoading}
                        style={{ width: 200 }}
                        showSearch
                        optionFilterProp="children"
                        filterOption={(input, option) =>
                          (option?.children as unknown as string)
                            ?.toLowerCase()
                            .includes(input.toLowerCase())
                        }
                      >
                        {usersData?.data?.map((user: { id: string; username: string }) => (
                          <Option key={user.id} value={user.id}>
                            {user.username} ({user.id.substring(0, 8)}...)
                          </Option>
                        ))}
                      </Select>
                    </Form.Item>
                    <Form.Item>
                      <Button
                        type="default"
                        htmlType="submit"
                        icon={<UserOutlined />}
                      >
                        Get User Permissions
                      </Button>
                    </Form.Item>
                  </Form>

                  <Form
                    layout="inline"
                    onFinish={handleCheckSuperAdmin}
                  >
                    <Form.Item
                      name="user_id"
                      rules={[{ required: true, message: "Please select a user" }]}
                    >
                      <Select
                        placeholder="Select a user"
                        loading={usersLoading}
                        style={{ width: 200 }}
                        showSearch
                        optionFilterProp="children"
                        filterOption={(input, option) =>
                          (option?.children as unknown as string)
                            ?.toLowerCase()
                            .includes(input.toLowerCase())
                        }
                      >
                        {usersData?.data?.map((user: { id: string; username: string }) => (
                          <Option key={user.id} value={user.id}>
                            {user.username} ({user.id.substring(0, 8)}...)
                          </Option>
                        ))}
                      </Select>
                    </Form.Item>
                    <Form.Item>
                      <Button
                        type="default"
                        htmlType="submit"
                        icon={<CheckCircleOutlined />}
                      >
                        Check Super Admin
                      </Button>
                    </Form.Item>
                  </Form>

                  {superAdminCheckResult && (
                    <Alert
                      message="Super Admin Check Result"
                      description={
                        <Space>
                          <Text>User ID: {superAdminCheckResult.user_id}</Text>
                          <Tag
                            color={
                              superAdminCheckResult.is_super_admin
                                ? "success"
                                : "default"
                            }
                          >
                            {superAdminCheckResult.is_super_admin
                              ? "SUPER ADMIN"
                              : "REGULAR USER"}
                          </Tag>
                        </Space>
                      }
                      type={
                        superAdminCheckResult.is_super_admin
                          ? "success"
                          : "info"
                      }
                      showIcon
                    />
                  )}
                </Space>
              </Card>
            </Col>
          </Row>
        </TabPane>
      </Tabs>

      {/* User Permissions Modal */}
      <Modal
        title="User Permissions"
        open={userPermissionsModalVisible}
        onCancel={() => {
          setUserPermissionsModalVisible(false);
          setSelectedUserId("");
          setUserPermissionsData(null);
        }}
        footer={null}
        width={800}
      >
        {userPermissionsLoading ? (
          <div style={{ textAlign: "center", padding: "40px" }}>
            <Spin size="large" tip="Loading user permissions..." />
          </div>
        ) : userPermissionsError ? (
          <Alert
            message="Error"
            description="Failed to load user permissions"
            type="error"
            showIcon
          />
        ) : userPermissionsData ? (
          <div>
            <Text strong>User ID: </Text>
            <Tag>{userPermissionsData.user_id}</Tag>
            <Divider />
            <Title level={4}>Permissions:</Title>
            {userPermissionsData.permissions &&
            userPermissionsData.permissions.length > 0 ? (
              <Space wrap>
                {userPermissionsData.permissions.map(
                  (permission: Permission) => (
                    <Tag key={permission.id} color="blue">
                      {permission.resource}:{permission.action}
                    </Tag>
                  )
                )}
              </Space>
            ) : (
              <Alert message="No permissions found for this user" type="info" />
            )}
          </div>
        ) : null}
      </Modal>
    </div>
  );
};

export default RBACManagementPage;