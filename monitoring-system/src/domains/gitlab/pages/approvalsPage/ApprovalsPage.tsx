import { Layout, Select, Spin, Result } from 'antd';
import ApprovalsTable from '../../components/ApprovalsTable/ApprovalsTable';
import { useGetAuthorizationRequestsQuery, setApprovalFilter } from '../../store';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '@/app/store';

const { Option } = Select;

const ApprovalsPage = () => {
  const dispatch = useDispatch();
  const approvalFilter = useSelector((state: RootState) => state.gitlab.approvalFilter);
  const { data: requests, isLoading, error } = useGetAuthorizationRequestsQuery();

  if (isLoading) {
    return (
      <Layout.Content style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Spin size="large" />
      </Layout.Content>
    );
  }

  if (error) {
    return (
      <Layout.Content>
        <Result
          status="error"
          title="Error"
          subTitle="Failed to load authorization requests. Please try again."
        />
      </Layout.Content>
    );
  }

  return (
    <Layout.Content style={{ padding: '24px' }}>
      <Select
        value={approvalFilter}
        onChange={(value) => dispatch(setApprovalFilter(value))}
        style={{ width: 200, marginBottom: 16 }}
      >
        <Option value="all">All</Option>
        <Option value="pending">Pending</Option>
        <Option value="accepted">Approved</Option>
        <Option value="rejected">Rejected</Option>
      </Select>
      <ApprovalsTable requests={requests || []} filter={approvalFilter} />
    </Layout.Content>
  );
};

export default ApprovalsPage;