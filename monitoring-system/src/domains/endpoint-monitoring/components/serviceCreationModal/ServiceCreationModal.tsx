import React, { useEffect, useState } from 'react';
import {
  Modal,
  Input,
  Select,
  Button,
  Row,
  Col,
  Typography,
  Divider,
  message,
  Collapse,
  Space,
  Tag,
} from 'antd';
import {
  CloseOutlined,
  ApiOutlined,
  SettingOutlined,
  DownOutlined,
  UpOutlined,
  GitlabOutlined,
  ContainerOutlined,
  TagsOutlined,
  FileTextOutlined,
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;
const { TextArea } = Input;
const { Panel } = Collapse;

type ErrorFields = {
  service_name?: string;
  server_name?: string;
  url?: string;
  expected_status?: string;
};

interface ServiceCreationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (formData: any) => void;
  formData: any;
  setFormData: React.Dispatch<React.SetStateAction<any>>;
  loading?: boolean;
  isEditMode?: boolean;
}

const ServiceCreationModal: React.FC<ServiceCreationModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  formData,
  setFormData,
  loading = false,
  isEditMode = false,
}) => {
  const [errors, setErrors] = useState<ErrorFields>({});
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [tagInput, setTagInput] = useState('');

  const validateForm = () => {
    const newErrors: ErrorFields = {};

    if (!formData.service_name?.trim()) {
      newErrors.service_name = 'Service name is required';
    } else if (formData.service_name.length < 2) {
      newErrors.service_name = 'Service name must be at least 2 characters';
    }

    if (!formData.server_name?.trim()) {
      newErrors.server_name = 'Server name is required';
    }

    if (!formData.url?.trim()) {
      newErrors.url = 'URL is required';
    } else if (!/^https?:\/\/.+/.test(formData.url)) {
      newErrors.url = 'Please enter a valid URL';
    }

    if (!formData.expected_status?.trim()) {
      newErrors.expected_status = 'Expected status is required';
    } else if (!/^[1-5][0-9][0-9]$/.test(formData.expected_status)) {
      newErrors.expected_status = 'Please enter a valid HTTP status code';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = () => {
    if (!validateForm()) {
      message.error('Please fix the validation errors');
      return;
    }

    const completeFormData = {
      ...formData,
      gitlab_url: formData.gitlab_url || null,
      docker_container_name: formData.docker_container_name || null,
      kubernetes_pod_name: formData.kubernetes_pod_name || null,
      tags: formData.tags || [],
      description: formData.description || null,
    };

    onSubmit(completeFormData);
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev: any) => ({ ...prev, [field]: value }));
    if (errors[field as keyof ErrorFields]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const handleTagAdd = () => {
    if (tagInput.trim() && !formData.tags?.includes(tagInput.trim())) {
      const newTags = [...(formData.tags || []), tagInput.trim()];
      setFormData((prev: any) => ({ ...prev, tags: newTags }));
      setTagInput('');
    }
  };

  const handleTagRemove = (tagToRemove: string) => {
    const newTags = formData.tags?.filter((tag: string) => tag !== tagToRemove) || [];
    setFormData((prev: any) => ({ ...prev, tags: newTags }));
  };

  const handleTagInputKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      handleTagAdd();
    }
  };

  return (
    <Modal
      title={null}
      open={isOpen}
      onCancel={onClose}
      footer={null}
      width={700}
      centered
      closable={false}
      style={{ borderRadius: '12px', padding: 0 }}
      destroyOnClose
    >
      <div style={{
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '24px',
        borderRadius: '12px 12px 0 0',
        position: 'relative',
      }}>
        <Button
          type="text"
          icon={<CloseOutlined />}
          onClick={onClose}
          style={{
            position: 'absolute',
            top: '16px',
            right: '16px',
            color: 'white',
            border: 'none',
            fontSize: '16px',
          }}
        />
        <Title level={2} style={{ color: 'white', margin: 0, fontWeight: '600' }}>
          <ApiOutlined style={{ marginRight: '12px' }} />
          {isEditMode ? 'Edit Service' : 'Create New Service'}
        </Title>
        <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: '16px' }}>
          {isEditMode ? 'Update your service endpoint details below' : 'Configure your new service endpoint with the details below'}
        </Text>
      </div>

      <div style={{ padding: '32px' }}>
        <div style={{ marginBottom: '24px' }}>
          <Title level={4} style={{ marginBottom: '16px', color: '#1f2937' }}>
            Basic Configuration
          </Title>

          <Row gutter={[24, 0]}>
            <Col span={12}>
              <div style={{ marginBottom: '16px' }}>
                <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                  Service Name *
                </Text>
                <Input
                  placeholder="Enter service name"
                  value={formData.service_name}
                  onChange={(e) => handleInputChange('service_name', e.target.value)}
                  status={errors.service_name ? 'error' : ''}
                  style={{ borderRadius: '6px', height: '44px' }}
                />
                {errors.service_name && (
                  <Text type="danger" style={{ fontSize: '12px', marginTop: '4px' }}>
                    {errors.service_name}
                  </Text>
                )}
              </div>
            </Col>

            <Col span={12}>
              <div style={{ marginBottom: '16px' }}>
                <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                  Server Name *
                </Text>
                <Input
                  placeholder="Enter server name"
                  value={formData.server_name}
                  onChange={(e) => handleInputChange('server_name', e.target.value)}
                  status={errors.server_name ? 'error' : ''}
                  style={{ borderRadius: '6px', height: '44px' }}
                />
                {errors.server_name && (
                  <Text type="danger" style={{ fontSize: '12px', marginTop: '4px' }}>
                    {errors.server_name}
                  </Text>
                )}
              </div>
            </Col>
          </Row>

          <div style={{ marginBottom: '16px' }}>
            <Text strong style={{ display: 'block', marginBottom: '8px' }}>
              Service URL *
            </Text>
            <Input
              placeholder="https://api.example.com/endpoint"
              value={formData.url}
              onChange={(e) => handleInputChange('url', e.target.value)}
              status={errors.url ? 'error' : ''}
              style={{ borderRadius: '6px', height: '44px' }}
            />
            {errors.url && (
              <Text type="danger" style={{ fontSize: '12px', marginTop: '4px' }}>
                {errors.url}
              </Text>
            )}
          </div>

          <Row gutter={[24, 0]}>
            <Col span={12}>
              <div style={{ marginBottom: '16px' }}>
                <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                  API Method *
                </Text>
                <Select
                  value={formData.api_method}
                  onChange={(value) => handleInputChange('api_method', value)}
                  style={{ width: '100%', borderRadius: '6px' }}
                  size="large"
                >
                  <Option value="GET">GET</Option>
                  <Option value="POST">POST</Option>
                  <Option value="PUT">PUT</Option>
                  <Option value="DELETE">DELETE</Option>
                </Select>
              </div>
            </Col>

            <Col span={12}>
              <div style={{ marginBottom: '16px' }}>
                <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                  Expected Status Code *
                </Text>
                <Input
                  placeholder="200"
                  value={formData.expected_status}
                  onChange={(e) => handleInputChange('expected_status', e.target.value)}
                  status={errors.expected_status ? 'error' : ''}
                  style={{ borderRadius: '6px', height: '44px' }}
                />
                {errors.expected_status && (
                  <Text type="danger" style={{ fontSize: '12px', marginTop: '4px' }}>
                    {errors.expected_status}
                  </Text>
                )}
              </div>
            </Col>
          </Row>
        </div>

        <Button
          type="link"
          onClick={() => setShowAdvanced(!showAdvanced)}
          style={{
            padding: 0,
            height: 'auto',
            fontSize: '16px',
            fontWeight: '500',
            color: '#667eea',
          }}
          icon={<SettingOutlined />}
        >
          Advanced Configuration (Optional)
          {showAdvanced ? (
            <UpOutlined style={{ marginLeft: '8px', fontSize: '12px' }} />
          ) : (
            <DownOutlined style={{ marginLeft: '8px', fontSize: '12px' }} />
          )}
        </Button>

        {showAdvanced && (
          <div style={{
            marginBottom: '24px',
            padding: '24px',
            backgroundColor: '#f8f9fa',
            borderRadius: '8px',
            border: '1px solid #e9ecef',
          }}>
            <Title level={4} style={{ marginBottom: '16px', color: '#1f2937' }}>
              <SettingOutlined style={{ marginRight: '8px' }} />
              Advanced Configuration
            </Title>
            <Text type="secondary" style={{ display: 'block', marginBottom: '20px' }}>
              These fields are optional and can be left empty if not applicable to your service.
            </Text>

            <Row gutter={[24, 16]}>
              <Col span={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                    <GitlabOutlined style={{ marginRight: '6px', color: '#fc6d26' }} />
                    GitLab URL
                  </Text>
                  <Input
                    placeholder="https://gitlab.com/your-project"
                    value={formData.gitlab_url || ''}
                    onChange={(e) => handleInputChange('gitlab_url', e.target.value)}
                    style={{ borderRadius: '6px', height: '40px' }}
                  />
                </div>
              </Col>

              <Col span={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                    <ContainerOutlined style={{ marginRight: '6px', color: '#0db7ed' }} />
                    Docker Container Name
                  </Text>
                  <Input
                    placeholder="my-service-container"
                    value={formData.docker_container_name || ''}
                    onChange={(e) => handleInputChange('docker_container_name', e.target.value)}
                    style={{ borderRadius: '6px', height: '40px' }}
                  />
                </div>
              </Col>

              <Col span={24}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                    <ContainerOutlined style={{ marginRight: '6px', color: '#326ce5' }} />
                    Kubernetes Pod Name
                  </Text>
                  <Input
                    placeholder="my-service-deployment-7b8c9d"
                    value={formData.kubernetes_pod_name || ''}
                    onChange={(e) => handleInputChange('kubernetes_pod_name', e.target.value)}
                    style={{ borderRadius: '6px', height: '40px' }}
                  />
                </div>
              </Col>

              <Col span={24}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                    <TagsOutlined style={{ marginRight: '6px', color: '#52c41a' }} />
                    Tags
                  </Text>
                  <div style={{ marginBottom: '8px' }}>
                    <Input
                      placeholder="Enter tag and press Enter or comma"
                      value={tagInput}
                      onChange={(e) => setTagInput(e.target.value)}
                      onKeyPress={handleTagInputKeyPress}
                      onBlur={handleTagAdd}
                      style={{ borderRadius: '6px', height: '40px' }}
                      suffix={
                        <Button
                          type="link"
                          size="small"
                          onClick={handleTagAdd}
                          disabled={!tagInput.trim()}
                        >
                          Add
                        </Button>
                      }
                    />
                  </div>
                  {formData.tags && formData.tags.length > 0 && (
                    <Space wrap>
                      {formData.tags.map((tag: string, index: number) => (
                        <Tag
                          key={index}
                          closable
                          onClose={() => handleTagRemove(tag)}
                          style={{ marginBottom: '4px' }}
                        >
                          {tag}
                        </Tag>
                      ))}
                    </Space>
                  )}
                </div>
              </Col>

              <Col span={24}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ display: 'block', marginBottom: '8px' }}>
                    <FileTextOutlined style={{ marginRight: '6px', color: '#722ed1' }} />
                    Description
                  </Text>
                  <TextArea
                    placeholder="Brief description of what this service does..."
                    value={formData.description || ''}
                    onChange={(e) => handleInputChange('description', e.target.value)}
                    rows={3}
                    style={{ borderRadius: '6px', resize: 'vertical' }}
                    maxLength={500}
                    showCount
                  />
                </div>
              </Col>
            </Row>
          </div>
        )}

        <Divider />

        <Row justify="end" gutter={[16, 0]}>
          <Col>
            <Button onClick={onClose} size="large">
              Cancel
            </Button>
          </Col>
          <Col>
            <Button
              type="primary"
              onClick={handleSubmit}
              loading={loading}
              size="large"
              style={{
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                border: 'none',
              }}
            >
              {loading ? (isEditMode ? 'Updating...' : 'Creating...') : (isEditMode ? 'Update Service' : 'Create Service')}
            </Button>
          </Col>
        </Row>
      </div>
    </Modal>
  );
};

export default ServiceCreationModal;