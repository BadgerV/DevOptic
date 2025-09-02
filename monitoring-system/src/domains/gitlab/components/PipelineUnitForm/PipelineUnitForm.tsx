import { Form, Select, Button, Card, Typography, message, Space } from "antd";
import { useCreatePipelineUnitMutation } from "../../store";
import { useEffect } from "react";
import { Service } from "../../types";

const { Option } = Select;
const { Title, Text } = Typography;

interface PipelineUnitFormProps {
  services: Service[];
}

interface FormValues {
  macro_service_id: string;
  micro_service_ids: string[];
}

const PipelineUnitForm = ({ services }: PipelineUnitFormProps) => {
  console.log(services);
  const [form] = Form.useForm();
  const [createPipelineUnit, { isLoading, error }] =
    useCreatePipelineUnitMutation();
  const [messageApi, contextHolder] = message.useMessage();

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: "error",
        content: "Failed to create pipeline unit. Please try again.",
      });
    }
  }, [error, messageApi]);

  const onFinish = async (values: FormValues) => {
    try {
      await createPipelineUnit({
        macro_service_id: values.macro_service_id,
        micro_service_ids: values.micro_service_ids || [],
      }).unwrap();

      messageApi.open({
        type: "success",
        content: "Pipeline unit created successfully!",
      });
      form.resetFields();
    } catch {
      // Error handled by useEffect
    }
  };

  const macroService = Form.useWatch("macro_service_id", form);
  const microServicesSelected = Form.useWatch("micro_service_ids", form) || [];

  // Filter services into macro and micro
  const macroServices = services.filter((s) => s.type === "macro");
  const microServices = services.filter((s) => s.type === "micro");

  return (
    <>
      {contextHolder}
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        style={{ maxWidth: 600 }}
      >
        {/* Macro Service */}
        <Form.Item
          label="Macro Service"
          name="macro_service_id"
          rules={[{ required: true, message: "Please select a macro service" }]}
        >
          <Select
            placeholder="Select macro service"
            disabled={macroServices.length === 0}
          >
            {macroServices.length > 0 ? (
              macroServices.map((service) => (
                <Option key={service.id} value={service.id}>
                  {service.name}
                </Option>
              ))
            ) : (
              <Option disabled value="">
                No macro services available
              </Option>
            )}
          </Select>
        </Form.Item>

        {/* Microservices */}
        <Form.Item label="Microservices" name="micro_service_ids">
          <Select
            mode="multiple"
            placeholder="Select microservices"
            tagRender={({ label }) => (
              <span
                style={{
                  marginRight: 8,
                  background: "#e6f7ff",
                  padding: "2px 8px",
                  borderRadius: 4,
                }}
              >
                {label}
              </span>
            )}
          >
            {microServices.map((service) => (
              <Option key={service.id} value={service.id}>
                {service.name}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={isLoading}>
            Create Pipeline Unit
          </Button>
        </Form.Item>
      </Form>

      {/* Preview Card */}
      {(macroService || microServicesSelected.length > 0) && (
        <Card style={{ marginTop: 16 }}>
          <Title level={4}>Preview</Title>
          <Space direction="vertical">
            {macroService && (
              <Text>
                <strong>Macro Service:</strong>{" "}
                {services.find((s) => s.id === macroService)?.name || "Unknown"}
              </Text>
            )}
            {microServicesSelected.length > 0 && (
              <Text>
                <strong>Microservices:</strong>{" "}
                {microServicesSelected
                  .map(
                    (id: string) =>
                      services.find((s) => s.id === id)?.name || "Unknown"
                  )
                  .join(", ")}
              </Text>
            )}
          </Space>
        </Card>
      )}
    </>
  );
};

export default PipelineUnitForm;
