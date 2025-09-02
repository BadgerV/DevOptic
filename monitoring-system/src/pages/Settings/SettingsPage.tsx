import React, { useEffect, useState } from "react";
import "./settingsPage.css";
import {
  Card,
  Tabs,
  Form,
  Input,
  Button,
  Tag,
  Typography,
  Space,
  Divider,
  message,
  Progress,
} from "antd";
import {
  UserOutlined,
  LockOutlined,
  MailOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons";
import { useSelector } from "react-redux";
import { RootState } from "@/app/store";
import {
  useChangePasswordMutation,
  useSetDeliveryEmailMutation,
  useGetDeliveryEmailQuery,
} from "@shared/services/api/authApi";
import { current } from "@reduxjs/toolkit";
import Sidebar from "@/shared/components/sidebar/Sidebar";

const { Title, Text } = Typography;
const { TabPane } = Tabs;

const SettingsPage = () => {
  const user = useSelector((state: RootState) => state.auth.user);

  const [passwordForm] = Form.useForm();
  const [emailForm] = Form.useForm();

  // Queries & Mutations
  const {
    data: deliveryEmailData,
    isLoading: emailLoading,
    refetch: refetchDeliveryEmail,
  } = useGetDeliveryEmailQuery();

  const [changePassword, { isLoading: changingPassword }] =
    useChangePasswordMutation();

  const [setDeliveryEmail, { isLoading: settingEmail }] =
    useSetDeliveryEmailMutation();

  // State management
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [passwordStrengthText, setPasswordStrengthText] = useState("");
  const [verificationSent, setVerificationSent] = useState(false);

  const currentEmail = deliveryEmailData?.delivery_email ?? user?.email;

  useEffect(() => {
    console.log("This is working", deliveryEmailData);
  }, [deliveryEmailData]);

  // Password strength calculation
  const calculatePasswordStrength = (password: string) => {
    if (!password) {
      setPasswordStrength(0);
      setPasswordStrengthText("");
      return;
    }

    let strength = 0;
    const checks = {
      length: password.length >= 8,
      lowercase: /[a-z]/.test(password),
      uppercase: /[A-Z]/.test(password),
      numbers: /\d/.test(password),
      special: /[!@#$%^&*(),.?":{}|<>]/.test(password),
    };

    strength = Object.values(checks).filter(Boolean).length;
    const strengthPercentage = (strength / 5) * 100;
    setPasswordStrength(strengthPercentage);

    const strengthLabels = ["Very Weak", "Weak", "Fair", "Good", "Strong"];
    setPasswordStrengthText(strengthLabels[strength - 1] || "");
  };

  const [messageApi, contextHolder] = message.useMessage();

  // Handle password change
  const handlePasswordChange = async (values: any) => {
    try {
      await changePassword({
        old_password: values.currentPassword,
        new_password: values.newPassword,
      }).unwrap();

      messageApi.open({
        type: "success",
        content: "Password changed successful!",
      });

      passwordForm.resetFields();
      setPasswordStrength(0);
      setPasswordStrengthText("");
    } catch (err: any) {
      messageApi.open({
        type: "error",
        content: "Failed to change password",
      });
    }
  };

  // Handle email update
  const handleEmailUpdate = async (values: any) => {
    try {
      await setDeliveryEmail({ delivery_email: values.email }).unwrap();

      messageApi.open({
        type: "success",
        content: "Email chnaged successfully",
      });

      emailForm.resetFields();
      setVerificationSent(false);
      refetchDeliveryEmail();
    } catch (err: any) {
      message.error(err?.data?.message || "Failed to update email");
    }
  };

  // Handle email verification (stub — implement if backend supports)
  const handleEmailVerification = async () => {
    try {
      // TODO: implement real verify API if available
      setVerificationSent(true);
      message.success("Verification email sent! Please check your inbox.");
    } catch {
      message.error("Failed to send verification email.");
    }
  };

  // Password strength color
  const getPasswordStrengthColor = () => {
    if (passwordStrength < 20) return "#ff4d4f";
    if (passwordStrength < 40) return "#ff7a45";
    if (passwordStrength < 60) return "#ffa940";
    if (passwordStrength < 80) return "#52c41a";
    return "#389e0d";
  };

    const [isSidebarOpen, setSidebarOpen] = useState(true);
  
    const toggleSidebar = () => {
      setSidebarOpen(!isSidebarOpen);
    };
  

  return (
    <div className="settings-main">
      <Sidebar isOpen={isSidebarOpen} onClose={toggleSidebar} />
      <div className="settings-container">
        {contextHolder}
        <div className="settings-header">
          <Title level={2}>Settings</Title>
          <Text type="secondary">
            Manage your account settings and preferences
          </Text>
        </div>

        <Tabs defaultActiveKey="profile" className="settings-tabs">
          <TabPane
            tab={
              <span>
                <UserOutlined />
                Profile
              </span>
            }
            key="profile"
          >
            <Card title="Profile Information" className="settings-card">
              <Form layout="vertical">
                <Form.Item label="Email">
                  <Input
                    value={user?.email}
                    disabled
                    prefix={<UserOutlined />}
                    className="readonly-input"
                  />
                </Form.Item>
                <Form.Item label="Username">
                  <Input
                    value={user?.username}
                    disabled
                    prefix={<UserOutlined />}
                    className="readonly-input"
                  />
                </Form.Item>
              </Form>
              <Text type="secondary" className="readonly-note">
                Profile information is read-only and cannot be modified.
              </Text>
            </Card>
          </TabPane>

          <TabPane
            tab={
              <span>
                <LockOutlined />
                Security
              </span>
            }
            key="security"
          >
            <Space
              direction="vertical"
              size="large"
              className="security-section"
            >
              {/* Change Password Card */}
              <Card title="Change Password" className="settings-card">
                <Form
                  form={passwordForm}
                  layout="vertical"
                  onFinish={handlePasswordChange}
                  autoComplete="off"
                >
                  <Form.Item
                    label="Current Password"
                    name="currentPassword"
                    rules={[
                      {
                        required: true,
                        message: "Please enter your current password",
                      },
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="Enter current password"
                      iconRender={(visible) =>
                        visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
                      }
                    />
                  </Form.Item>

                  <Form.Item
                    label="New Password"
                    name="newPassword"
                    rules={[
                      {
                        required: true,
                        message: "Please enter a new password",
                      },
                      {
                        min: 8,
                        message: "Password must be at least 8 characters long",
                      },
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="Enter new password"
                      onChange={(e) =>
                        calculatePasswordStrength(e.target.value)
                      }
                      iconRender={(visible) =>
                        visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
                      }
                    />
                  </Form.Item>

                  {passwordStrength > 0 && (
                    <div className="password-strength">
                      <Text type="secondary">Password Strength:</Text>
                      <Progress
                        percent={passwordStrength}
                        strokeColor={getPasswordStrengthColor()}
                        showInfo={false}
                        size="small"
                      />
                      <Text
                        style={{
                          color: getPasswordStrengthColor(),
                          fontWeight: 500,
                        }}
                        className="strength-text"
                      >
                        {passwordStrengthText}
                      </Text>
                    </div>
                  )}

                  <Form.Item
                    label="Confirm Password"
                    name="confirmPassword"
                    dependencies={["newPassword"]}
                    rules={[
                      {
                        required: true,
                        message: "Please confirm your new password",
                      },
                      ({ getFieldValue }) => ({
                        validator(_, value) {
                          if (
                            !value ||
                            getFieldValue("newPassword") === value
                          ) {
                            return Promise.resolve();
                          }
                          return Promise.reject(
                            new Error("Passwords do not match")
                          );
                        },
                      }),
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="Confirm new password"
                      iconRender={(visible) =>
                        visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
                      }
                    />
                  </Form.Item>

                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={changingPassword}
                      className="action-button"
                    >
                      Change Password
                    </Button>
                  </Form.Item>
                </Form>
              </Card>

              <Divider />

              {/* Delivery Email Card */}
              <Card title="Delivery Email" className="settings-card">
                <div className="email-status">
                  <Space align="center">
                    <Text strong>Current Email:</Text>
                    <Text>{emailLoading ? "Loading..." : currentEmail}</Text>
                    <Tag color={"warning"} icon={<ExclamationCircleOutlined />}>
                      Unverified
                    </Tag>
                  </Space>
                </div>

                <Form
                  form={emailForm}
                  layout="vertical"
                  onFinish={handleEmailUpdate}
                  initialValues={{ email: currentEmail }}
                  className="email-form"
                >
                  <Form.Item
                    label="Email Address"
                    name="email"
                    rules={[
                      {
                        required: true,
                        message: "Please enter your email address",
                      },
                      {
                        type: "email",
                        message: "Please enter a valid email address",
                      },
                    ]}
                  >
                    <Input
                      prefix={<MailOutlined />}
                      placeholder="Enter email address"
                      type="email"
                    />
                  </Form.Item>

                  <Form.Item>
                    <Space>
                      <Button
                        type="primary"
                        htmlType="submit"
                        loading={settingEmail}
                        className="action-button"
                      >
                        Update Email
                      </Button>
                      {/* <Button
                      onClick={handleEmailVerification}
                      disabled={verificationSent}
                      className="verify-button"
                    >
                      {verificationSent
                        ? "Resend Verification"
                        : "Verify Email"}
                    </Button> */}
                    </Space>
                  </Form.Item>
                </Form>

                {verificationSent && (
                  <div className="verification-notice">
                    <Text type="secondary">
                      Verification email sent! Please check your inbox.
                    </Text>
                  </div>
                )}
              </Card>
            </Space>
          </TabPane>
        </Tabs>
      </div>
    </div>
  );
};

export default SettingsPage;
