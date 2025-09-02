import { Form, Input, Select, Button, message } from "antd";
import { useEffect } from "react";
import { useCreateServiceMutation } from "../../store";
const { Option } = Select;

import "./serviceForm.css";

const ServiceForm = () => {
  const [form] = Form.useForm();
  const [createService, { isLoading, error }] = useCreateServiceMutation();
  const [messageApi, contextHolder] = message.useMessage();

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: "error",
        content: "Failed to create service. Please try again.",
      });
    }
  }, [error, messageApi]);

  const onFinish = async (values: any) => {
    console.log(values);
    try {
      await createService(values).unwrap();
      messageApi.open({
        type: "success",
        content: "Service created successfully!",
      });
      form.resetFields();
    } catch {
      // Error handled by useEffect
    }
  };

  return (
    <>
      {contextHolder}
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        style={{ maxWidth: 600 }}
      >
        <Form.Item
          label="GitLab Repo ID"
          name="gitlab_repo_id"
          rules={[
            { required: true, message: "Please enter the GitLab Repo ID" },
          ]}
          tooltip="GitLab Repo ID can be found in your repo settings"
        >
          <Input placeholder="Enter GitLab Repo ID" />
        </Form.Item>
        <Form.Item
          label="Name"
          name="name"
          rules={[{ required: true, message: "Please enter the service name" }]}
        >
          <Input placeholder="Enter service name" />
        </Form.Item>
        <Form.Item
          label="URL"
          name="url"
          rules={[
            { required: true, message: "Please enter the service URL" },
            { type: "url", message: "Please enter a valid URL" },
          ]}
        >
          <Input placeholder="Enter service URL" />
        </Form.Item>
        <Form.Item
          label="Type"
          name="type"
          rules={[
            { required: true, message: "Please select the service type" },
          ]}
        >
          <Select placeholder="Select service type">
            <Option value="micro">Micro</Option>
            <Option value="macro">Macro</Option>
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={isLoading}>
            Register Service
          </Button>
        </Form.Item>
      </Form>
    </>
  );
};

export default ServiceForm;
