import React, { useState } from "react";
import { Button, Typography, Card, Divider, message, Input } from "antd";
import { UserOutlined, LockOutlined, MailOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import {
  clearUserState,
  useRegisterMutation,
} from "@/shared/services/api/authApi";
import { formatDefaultLocale } from "d3";
import { useDispatch } from "react-redux";

const { Title, Text, Link } = Typography;

interface RegisterPageProps {
  onSwitchToLogin: () => void;
}

interface RegisterFormData {
  username: string;
  email: string;
  password: string;
}

interface RegisterFormErrors {
  username?: string;
  email?: string;
  password?: string;
}

// Register Page Component
const RegisterPage = () => {
  const dispatch = useDispatch();

  dispatch(clearUserState());
  const [loading, setLoading] = useState<boolean>(false);
  const [formData, setFormData] = useState<RegisterFormData>({
    username: "",
    email: "",
    password: "",
  });
  const [errors, setErrors] = useState<RegisterFormErrors>({});

  const handleInputChange = (field: keyof RegisterFormData, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  const validateForm = (): boolean => {
    const newErrors: RegisterFormErrors = {};

    if (!formData.username) {
      newErrors.username = "Please enter your username!";
    } else if (formData.username.length < 3) {
      newErrors.username = "Username must be at least 3 characters!";
    }

    if (!formData.email) {
      newErrors.email = "Please enter your email!";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = "Please enter a valid email!";
    }

    if (!formData.password) {
      newErrors.password = "Please enter your password!";
    } else if (formData.password.length < 6) {
      newErrors.password = "Password must be at least 6 characters!";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  const [register, { isLoading }] = useRegisterMutation();

  const onSubmit = async () => {
    if (!validateForm()) return;

    try {
      const res: any = await register({
        username: formData.username,
        password: formData.password,
        email: formData.email,
      }).unwrap();

      messageApi.open({
        type: "success",
        content: "Login successful!",
      });
      console.log("Login success:", res.data.token);

      if (res?.data.token) {
        console.log(res.token);
        localStorage.setItem("token", res.data.token);
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
          content: "Creation of account failed!",
        });
        console.log("Creation of account failed:", err);
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
            Create Account
          </Title>
          <Text type="secondary" style={{ fontSize: "16px" }}>
            Sign up for a new account
          </Text>
        </div>

        <div>
          {/* Username */}
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

          {/* Email */}
          <div style={{ marginBottom: "20px" }}>
            <label
              style={{
                display: "block",
                marginBottom: "8px",
                color: "#374151",
                fontWeight: 500,
              }}
            >
              Email
            </label>
            <Input
              prefix={<MailOutlined style={{ color: "#9ca3af" }} />}
              placeholder="Enter your email"
              value={formData.email}
              onChange={(e) => handleInputChange("email", e.target.value)}
              style={{
                borderRadius: "8px",
                border: errors.email
                  ? "1px solid #ef4444"
                  : "1px solid #e5e7eb",
                padding: "12px 16px",
                fontSize: "16px",
              }}
              size="large"
            />
            {errors.email && (
              <div
                style={{ color: "#ef4444", fontSize: "14px", marginTop: "4px" }}
              >
                {errors.email}
              </div>
            )}
          </div>

          {/* Password */}
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

          {/* Register Button */}
          <div style={{ marginBottom: "20px" }}>
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
              Register
            </Button>
          </div>
        </div>

        <Divider style={{ margin: "24px 0", color: "#9ca3af" }}>or</Divider>

        <div style={{ textAlign: "center" }}>
          <Text style={{ color: "#6b7280", fontSize: "14px" }}>
            Already have an account?{" "}
            <Link
              //   onClick={onSwitchToLogin}
              style={{
                color: "#667eea",
                fontWeight: 500,
                textDecoration: "none",
              }}
              onClick={() => navigate("/login")}
            >
              Sign In
            </Link>
          </Text>
        </div>
      </Card>
    </div>
  );
};

export default RegisterPage;
