import React, { useState, useEffect, useRef } from "react";
import { Layout, Typography, Menu, Button } from "antd";
import {
  MenuOutlined,
  CloseOutlined,
  MonitorOutlined,
  CodeOutlined,
  DesktopOutlined,
  ContainerOutlined,
  CloudOutlined,
  LogoutOutlined,
} from "@ant-design/icons";
import "./rightSidebar.css";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { clearUserState } from "@/shared/services/api/authApi"; // Adjust path based on your project structure

const { Content } = Layout;
const { Title, Paragraph } = Typography;

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

// Sidebar Component
export const Sidebar: React.FC<SidebarProps> = ({ isOpen, onClose, links }) => {
  const sidebarRef = useRef<HTMLDivElement>(null);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  // Handle ESC key press
  useEffect(() => {
    const handleEscKey = (event: KeyboardEvent) => {
      if (event.key === "Escape" && isOpen) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscKey);
      // Focus the sidebar for accessibility
      sidebarRef.current?.focus();
    }

    return () => {
      document.removeEventListener("keydown", handleEscKey);
    };
  }, [isOpen, onClose]);

  // Handle click outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        sidebarRef.current &&
        !sidebarRef.current.contains(event.target as Node)
      ) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen, onClose]);

  const getIcon = (label: string) => {
    switch (label.toLowerCase()) {
      case "monitoring module":
        return <MonitorOutlined />;
      case "gitlab":
        return <CodeOutlined />;
      case "linux":
        return <DesktopOutlined />;
      case "docker":
        return <ContainerOutlined />;
      case "kubernetes":
        return <CloudOutlined />;
      default:
        return <MenuOutlined />;
    }
  };

  const handleLinkClick = (url: string) => {
    console.log(`Navigating to: ${url}`);
    navigate(url);
  };

  const handleLogout = () => {
    // Clear token from localStorage
    localStorage.removeItem("token");
    // Dispatch clearUserState action
    dispatch(clearUserState());
    // Navigate to login page or home
    navigate("/login");
  };

  return (
    <>
      {/* Overlay */}
      <div className={`sidebar-overlay ${isOpen ? "open" : ""}`} />

      {/* Sidebar */}
      <div
        ref={sidebarRef}
        className={`right-app-sidebar ${isOpen ? "open" : ""}`}
        tabIndex={-1}
      >
        <div className="sidebar-header">
          <Title level={4} style={{ margin: 0, color: "#fff" }}>
            Modules
          </Title>
          <button
            className="sidebar-close-btn"
            onClick={onClose}
            aria-label="Close sidebar"
          >
            <CloseOutlined />
          </button>
        </div>

        <div className="sidebar-content" style={{ display: "flex", flexDirection: "column", height: "100%" }}>
          <Menu
            mode="vertical"
            theme="dark"
            style={{ background: "transparent", border: "none", flex: 1 }}
          >
            {links.map((link, index) => (
              <Menu.Item
                key={index}
                icon={getIcon(link.label)}
                onClick={() => handleLinkClick(link.url)}
                style={{
                  color: "#fff",
                  borderRadius: "8px",
                  margin: "4px 0",
                }}
              >
                {link.label}
              </Menu.Item>
            ))}
          </Menu>
          <div style={{ padding: "16px", borderTop: "1px solid rgba(255, 255, 255, 0.1)" }}>
            <Button
              type="text"
              icon={<LogoutOutlined />}
              onClick={handleLogout}
              style={{
                width: "100%",
                color: "#fff",
                borderRadius: "8px",
                textAlign: "left",
                padding: "0 16px",
                height: "40px",
                background: "transparent",
                display: "flex",
                alignItems: "center",
                justifyContent: "flex-start",
                fontSize: "16px",
              }}
              aria-label="Log out"
            >
              Log Out
            </Button>
          </div>
        </div>
      </div>
    </>
  );
};