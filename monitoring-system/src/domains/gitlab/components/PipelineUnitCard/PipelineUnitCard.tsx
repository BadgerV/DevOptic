import {
  Card,
  Button,
  Tag,
  message,
  Badge,
  Typography,
  Divider,
  Space,
  Modal,
  Checkbox,
  List,
} from "antd";
import { useTriggerAuthorizationRequestMutation } from "../../store";
import { PipelineUnitWithServices } from "../../types";
import { PlayCircleOutlined, ClockCircleOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { RootState } from "@/app/store";
import { useSelector } from "react-redux";

const { Title, Text } = Typography;

interface PipelineUnitCardProps {
  unit: PipelineUnitWithServices;
}

const PipelineUnitCard = ({ unit }: PipelineUnitCardProps) => {
  const [isPending, setIsPending] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedMicroServices, setSelectedMicroServices] = useState<string[]>(
    []
  );

  useEffect(() => {
    console.log("Selected microservices:", selectedMicroServices);
  }, [selectedMicroServices]);

  const userId =
    useSelector((state: RootState) => state.auth.user?.userId) || "";
  const [triggerRequest, { isLoading }] =
    useTriggerAuthorizationRequestMutation();
  const [messageApi, contextHolder] = message.useMessage();

  const handleOpenModal = () => {
    // Start with no microservices selected
    setSelectedMicroServices([]);
    setIsModalVisible(true);
  };

  const handleModalOk = async () => {
    try {
      if (selectedMicroServices.length === 0) {
        messageApi.warning("Please select at least one microservice to run.");
        return;
      }

      await triggerRequest({
        id: unit.pipeline_unit.id,
        data: {
          requester_id: userId,
          selected_micro_service_ids: selectedMicroServices,
        },
      }).unwrap();

      setIsPending(true);
      setIsModalVisible(false);
      messageApi.success("Authorization request triggered!");
    } catch {
      messageApi.error("Failed to trigger request. Please try again.");
    }
  };

  const handleModalCancel = () => {
    setIsModalVisible(false);
  };

  const handleMicroServiceChange = (checkedValues: string[]) => {
    setSelectedMicroServices(checkedValues);
  };

  // Select/Deselect all helper
  const toggleSelectAll = () => {
    if (selectedMicroServices.length === unit.micro_services.length) {
      setSelectedMicroServices([]); // deselect all
    } else {
      setSelectedMicroServices(unit.micro_services.map((s) => s.id)); // select all
    }
  };

  return (
    <>
      {contextHolder}
      <Card
        hoverable
        style={{
          width: 340,
          borderRadius: 12,
          boxShadow: "0 2px 10px rgba(0,0,0,0.08)",
        }}
        bodyStyle={{ padding: "16px 20px" }}
        title={
          <Space direction="vertical" size={0}>
            <Title level={5} style={{ margin: 0 }}>
              Pipeline Unit
            </Title>
            <Text type="secondary" style={{ fontSize: 12 }}>
              ID: {unit.pipeline_unit.id}
            </Text>
          </Space>
        }
        actions={[
          <Button
            key="trigger"
            type={isPending ? "default" : "primary"}
            icon={isPending ? <ClockCircleOutlined /> : <PlayCircleOutlined />}
            onClick={handleOpenModal}
            disabled={isPending}
            style={{ borderRadius: 6 }}
          >
            {isPending ? "Pending Approval" : "Run Pipeline"}
          </Button>,
        ]}
      >
        {/* Macro Service */}
        <div style={{ marginBottom: 12 }}>
          <Text strong>Macro Service:</Text>{" "}
          <Tag color="blue" style={{ fontWeight: 500 }}>
            {unit.macro_service.name}
          </Tag>
        </div>
        <Divider style={{ margin: "8px 0" }} />

        {/* Microservices */}
        <div>
          <Text strong>Microservices:</Text>
          <div style={{ marginTop: 6 }}>
            {unit.micro_services.length > 0 ? (
              unit.micro_services.map((service) => (
                <Tag
                  key={service.id}
                  color="geekblue"
                  style={{
                    margin: "4px 4px 0 0",
                    borderRadius: 6,
                    fontSize: 13,
                  }}
                >
                  {service.name}
                </Tag>
              ))
            ) : (
              <Text type="secondary">None</Text>
            )}
          </div>
        </div>

        {/* Status Badge (if pending) */}
        {isPending && (
          <div style={{ marginTop: 16, textAlign: "center" }}>
            <Badge
              status="processing"
              text={<Text type="warning">Waiting for Approval</Text>}
            />
          </div>
        )}
      </Card>

      {/* Service Selection Modal */}
      <Modal
        title="Select Services to Run"
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={isLoading}
        okText="Run Pipeline"
        cancelText="Cancel"
        width={500}
      >
        <div style={{ marginBottom: 16 }}>
          <Text type="secondary">
            Select the services you want to include in this pipeline run:
          </Text>
        </div>

        {/* Macro Service - Always included */}
        <div style={{ marginBottom: 20 }}>
          <Title level={5} style={{ marginBottom: 8 }}>
            Macro Service (Always Included)
          </Title>
          <div
            style={{
              padding: 8,
              backgroundColor: "#e6f7ff",
              borderRadius: 6,
              border: "1px solid #91d5ff",
            }}
          >
            <Tag color="blue" style={{ marginLeft: 4 }}>
              {unit.macro_service.name}
            </Tag>
            <Text type="secondary" style={{ marginLeft: 8, fontSize: 12 }}>
              This service will always be included
            </Text>
          </div>
        </div>

        {/* Microservices */}
        {unit.micro_services.length > 0 && (
          <div>
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                marginBottom: 8,
              }}
            >
              <Title level={5} style={{ margin: 0 }}>
                Microservices
              </Title>
              <Button
                type="link"
                size="small"
                onClick={toggleSelectAll}
                style={{ padding: 0 }}
              >
                {selectedMicroServices.length === unit.micro_services.length
                  ? "Deselect All"
                  : "Select All"}
              </Button>
            </div>

            <div
              style={{
                maxHeight: 250,
                overflowY: "auto",
                paddingRight: 8,
                border: "1px solid #f0f0f0",
                borderRadius: 6,
              }}
            >
              <Checkbox.Group
                value={selectedMicroServices}
                onChange={handleMicroServiceChange}
                style={{ width: "100%" }}
              >
                <List
                  size="small"
                  dataSource={unit.micro_services}
                  renderItem={(service) => (
                    <List.Item
                      style={{
                        padding: "6px 8px",
                        borderBottom: "1px solid #f5f5f5",
                      }}
                    >
                      <Checkbox value={service.id}>
                        <Tag color="geekblue" style={{ marginLeft: 8 }}>
                          {service.name}
                        </Tag>
                      </Checkbox>
                    </List.Item>
                  )}
                />
              </Checkbox.Group>
            </div>
          </div>
        )}

        {/* Selection Summary */}
        <div
          style={{
            marginTop: 16,
            padding: 12,
            backgroundColor: "#f5f5f5",
            borderRadius: 6,
          }}
        >
          <Text strong>Selected Services: </Text>
          <Text>
            {selectedMicroServices.length} microservice
            {selectedMicroServices.length !== 1 ? "s" : ""} + 1 macro service
          </Text>
        </div>
      </Modal>
    </>
  );
};

export default PipelineUnitCard;
