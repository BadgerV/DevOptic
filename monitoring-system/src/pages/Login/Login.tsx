import React, { useEffect, useState } from "react";
import { Button, Typography, Card, Divider, message, Input } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import {
  clearUserState,
  useLoginMutation,
} from "@/shared/services/api/authApi";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
const { Title, Text, Link } = Typography;

interface FormData {
  username: string;
  password: string;
}

interface FormErrors {
  username?: string;
  password?: string;
}

// Login Page Component
const LoginPage = () => {
  const dispatch = useDispatch();

  dispatch(clearUserState());
  // localStorage.removeItem("token");
  const [formData, setFormData] = useState<FormData>({
    username: "",
    password: "",
  });
  const [errors, setErrors] = useState<FormErrors>({});
  const [login, { isLoading }] = useLoginMutation();

  const handleInputChange = (field: keyof FormData, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};
    if (!formData.username) newErrors.username = "Please enter your username!";
    if (!formData.password) newErrors.password = "Please enter your password!";
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  const onSubmit = async () => {
    if (!validateForm()) return;

    try {
      const res: any = await login({
        username: formData.username,
        password: formData.password,
      }).unwrap();

      messageApi.open({
        type: "success",
        content: "Login successful!",
      });
      console.log("Login success:", res.user.id);

      if (res?.token) {
        console.log(res.token);
        localStorage.setItem("token", res.token);
      }

      navigate("/dashboard");
    } catch (err: any) {
      if (err.data === "Request failed with status code 401") {
        messageApi.open({
          type: "error",
          content: "User not found!",
        });
      } else {
        messageApi.open({
          type: "error",
          content: "Login failed!",
        });
        console.log("Login failed:", err);
      }
    }
  };

  return (
    <div
      style={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
        padding: "20px",
      }}
    >
      {contextHolder}

      <Card
        style={{
          width: "100%",
          maxWidth: 400,
          boxShadow: "0 20px 40px rgba(0,0,0,0.1)",
          borderRadius: "16px",
          border: "none",
        }}
        bodyStyle={{
          padding: "40px 32px",
        }}
      >
        <div style={{ textAlign: "center", marginBottom: "32px" }}>
          <div
            style={{
              width: "60px",
              height: "60px",
              background: "linear-gradient(135deg, #667eea, #764ba2)",
              borderRadius: "50%",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              margin: "0 auto 16px",
              color: "white",
              fontSize: "24px",
            }}
          >
            <UserOutlined />
          </div>
          <Title level={2} style={{ margin: "0 0 8px", color: "#1f2937" }}>
            Welcome Back
          </Title>
          <Text type="secondary" style={{ fontSize: "16px" }}>
            Sign in to your account
          </Text>
        </div>

        <div>
          {/* Username Field */}
          <div style={{ marginBottom: "20px" }}>
            <label
              style={{
                display: "block",
                marginBottom: "8px",
                color: "#374151",
                fontWeight: 500,
              }}
            >
              Username
            </label>
            <Input
              prefix={<UserOutlined style={{ color: "#9ca3af" }} />}
              placeholder="Enter your username"
              value={formData.username}
              onChange={(e) => handleInputChange("username", e.target.value)}
              style={{
                borderRadius: "8px",
                border: errors.username
                  ? "1px solid #ef4444"
                  : "1px solid #e5e7eb",
                padding: "12px 16px",
                fontSize: "16px",
              }}
              size="large"
            />
            {errors.username && (
              <div
                style={{ color: "#ef4444", fontSize: "14px", marginTop: "4px" }}
              >
                {errors.username}
              </div>
            )}
          </div>

          {/* Password Field */}
          <div style={{ marginBottom: "24px" }}>
            <label
              style={{
                display: "block",
                marginBottom: "8px",
                color: "#374151",
                fontWeight: 500,
              }}
            >
              Password
            </label>
            <Input.Password
              prefix={<LockOutlined style={{ color: "#9ca3af" }} />}
              placeholder="Enter your password"
              value={formData.password}
              onChange={(e) => handleInputChange("password", e.target.value)}
              style={{
                borderRadius: "8px",
                border: errors.password
                  ? "1px solid #ef4444"
                  : "1px solid #e5e7eb",
                padding: "12px 16px",
                fontSize: "16px",
              }}
              size="large"
            />
            {errors.password && (
              <div
                style={{ color: "#ef4444", fontSize: "14px", marginTop: "4px" }}
              >
                {errors.password}
              </div>
            )}
          </div>

          {/* Login Button */}
          <div style={{ marginBottom: "16px" }}>
            <Button
              type="primary"
              loading={isLoading}
              onClick={onSubmit}
              style={{
                width: "100%",
                height: "48px",
                borderRadius: "8px",
                background: "linear-gradient(135deg, #667eea, #764ba2)",
                border: "none",
                fontSize: "16px",
                fontWeight: 600,
              }}
            >
              Login
            </Button>
          </div>

          {/* Forgot Password */}
          <div style={{ textAlign: "center", marginBottom: "20px" }}>
            {/* <Button
              type="link"
              style={{
                padding: 0,
                color: "#667eea",
                fontSize: "14px",
              }}
            >
              Forgot Password?
            </Button> */}
          </div>
        </div>

        <Divider style={{ margin: "24px 0", color: "#9ca3af" }}>or</Divider>

        <div style={{ textAlign: "center" }}>
          <Text style={{ color: "#6b7280", fontSize: "14px" }}>
            Don't have an account yet? Please{" "}
            <Link
              style={{
                color: "#667eea",
                fontWeight: 500,
                textDecoration: "none",
              }}
              onClick={() => navigate("/register")}
            >
              Sign Up
            </Link>
          </Text>
        </div>
      </Card>
    </div>
  );
};

export default LoginPage;
