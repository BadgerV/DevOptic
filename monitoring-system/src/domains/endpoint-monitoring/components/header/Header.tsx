import Input from "@/shared/components/input/Input";
import "./header.css";
import { useState, useEffect } from "react";
import {
  startEndpointCheck,
  stopEndpointCheck,
  getSchedulerStatus,
  createEndpoint,
} from "../../services/endpointMonitoringApiCalls";
import { Button, message } from "antd";
import ServiceCreationModal from "../serviceCreationModal/ServiceCreationModal";
import "antd/dist/reset.css";

type HeaderProps = {
  filter: string;
  onFilterChange: (value: string) => void;
};

const Header: React.FC<HeaderProps> = ({ filter, onFilterChange }) => {
  const [schedulerRunning, setSchedulerRunning] = useState<boolean | null>(
    null
  );
  const [loading, setLoading] = useState(false);

  // ✅ Always confirm initial status from backend
  useEffect(() => {
    const fetchStatus = async () => {
      try {
        setLoading(true);
        const status = await getSchedulerStatus();
        // Type assertion to fix 'unknown' error
        setSchedulerRunning((status as any).scheduler_running);
      } catch {
        message.error("❌ Failed to fetch scheduler status");
      } finally {
        setLoading(false);
      }
    };

    fetchStatus();
  }, []);

  const [messageApi, contextHolder] = message.useMessage();

  const handleToggle = async () => {
    try {
      setLoading(true);

      if (schedulerRunning) {
        await stopEndpointCheck();
        messageApi.open({
          type: "success",
          content: "Monitoring stopped successfully",
        });
      } else {
        await startEndpointCheck();
        messageApi.open({
          type: "success",
          content: "Monitoring started successfully",
        });
      }

      // ✅ Re-fetch actual status from backend after toggle
      const status: any = await getSchedulerStatus();
      setSchedulerRunning(status.scheduler_running);
    } catch (err: any) {
      console.log(err);
      if (err.message === "Request failed with status code 403") {
        messageApi.open({
          type: "error",
          content: "Unauthorized",
        });
      } else {
        messageApi.open({
          type: "error",
          content: "Failed to update monitoring state",
        });
      }
    } finally {
      setLoading(false);
    }
  };

  //Service creation Logic
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [formData, setFormData] = useState({
    service_name: "",
    server_name: "",
    url: "",
    api_method: "GET",
    expected_status: "200",
  });

  const handleOpenModal = () => {
    setFormData({
      service_name: "",
      server_name: "",
      url: "",
      api_method: "GET",
      expected_status: "200",
    }); // reset when open
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  const handleCreateService = async (serviceData: any) => {
    setIsCreating(true);
    try {
      // Call API
      await createEndpoint({
        service_name: serviceData.service_name,
        url: serviceData.url,
        server_name: serviceData.server_name,
        api_method: serviceData.api_method,
        expected_status_code: Number(serviceData.expected_status_code),
        gitlab_url: serviceData.gitlab_url || null,
        docker_container_name: serviceData.docker_container_name || null,
        kubernetes_pod_name: serviceData.kubernetes_pod_name || null,
        tags: serviceData.tags || [],
        description: serviceData.description || null,
        last_changed_by: serviceData.last_changed_by || null,
      });

      messageApi.open({
        type: "success",
        content: "Service created successfully!",
      });
      setIsModalOpen(false);
    } catch (err: any) {
      console.error("Create endpoint failed:", err);
      messageApi.open({
        type: "error",
        content: "Failed to create service",
      });
    } finally {
      setIsCreating(false);
    }
  };
  return (
    <div className="main-endpoint-header">
      {contextHolder}

      <div className="main-endpoint-header-left">
        <Input
          containerClassName="endpoint-main-header-input-container"
          inputClassName="endpoint-main-header-input-input"
          placeholder="Search"
          value={filter}
          onChange={(e: any) => onFilterChange(e.target.value)}
        />

        <Button
          type={schedulerRunning ? "default" : "primary"}
          danger={schedulerRunning ?? false} // red if running
          onClick={handleToggle}
          loading={loading}
          style={{ marginRight: "8px" }}
          disabled={schedulerRunning === null} // disabled until status is fetched
        >
          {schedulerRunning ? "Stop Checks" : "Start Checks"}
        </Button>
      </div>

      <div className="main-endpoint-header-right">
        <Button
          color="purple"
          variant="solid"
          onClick={() => window.location.reload()}
        >
          Refresh Endpoints
        </Button>

        <Button
          color="purple"
          variant="solid"
          onClick={() => handleOpenModal()}
        >
          Add Service
        </Button>
      </div>

      {/* Modal Component */}
      <ServiceCreationModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSubmit={handleCreateService}
        formData={formData}
        setFormData={setFormData}
        loading={isCreating}
      />
    </div>
  );
};

export default Header;
