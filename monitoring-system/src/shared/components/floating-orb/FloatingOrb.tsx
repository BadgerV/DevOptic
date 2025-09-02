import React, { useState, useEffect, useRef } from "react";
import { Layout, Typography, Menu } from "antd";
import {
  MenuOutlined,
  CloseOutlined,
  MonitorOutlined,
  CodeOutlined,
  DesktopOutlined,
  ContainerOutlined,
  CloudOutlined,
} from "@ant-design/icons";
import "./floating-orb.css";

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

// FloatingOrb Component
export const FloatingOrb: React.FC<FloatingOrbProps> = ({ onClick, isOpen }) => {
  return (
    <button
      className="floating-orb"
      onClick={onClick}
      aria-label="Open modules sidebar"
      type="button"
    >
      {isOpen ? <CloseOutlined /> : <MenuOutlined />}
    </button>
  );
};
