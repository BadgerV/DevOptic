import { Layout, Select, Input, Spin, Result } from 'antd';
import HistoryTable from '../../components/HistoryTable/HistoryTable';
import PipelineStatusModal from '../../components/PipelineStatusModal/PipelineStatusModal';
import { useGetExecutionHistoryQuery, setHistoryStatusFilter, setSearchTerm } from '../../store';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '@/app/store';

const { Option } = Select;
const { Search } = Input;

const ExecutionHistoryPage = () => {
  const dispatch = useDispatch();
  const historyStatusFilter = useSelector(
    (state: RootState) => state.gitlab.historyStatusFilter
  );
  const searchTerm = useSelector(
    (state: RootState) => state.gitlab.searchTerm
  );

  const { data: history, isLoading, error } = useGetExecutionHistoryQuery();

  if (isLoading) {
    return (
      <Layout.Content
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100vh',
        }}
      >
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
          subTitle="Failed to load execution history. Please try again."
        />
      </Layout.Content>
    );
  }

  return (
    <Layout.Content style={{ padding: '24px' }}>
      {/* Filters */}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Select
          value={historyStatusFilter}
          onChange={(value) => dispatch(setHistoryStatusFilter(value))}
          style={{ width: 200 }}
        >
          <Option value="all">All</Option>
          <Option value="running">Running</Option>
          <Option value="success">Success</Option>
          <Option value="failed">Failed</Option>
          <Option value="pending">Pending</Option>
        </Select>
        <Search
          placeholder="Search by requester or pipeline run ID"
          value={searchTerm}
          onChange={(e) => dispatch(setSearchTerm(e.target.value))}
          style={{ width: 300 }}
        />
      </div>

      {/* History Table */}
      <HistoryTable
        history={history || []}
        filter={historyStatusFilter}
        searchTerm={searchTerm}
      />
      {/* Pipeline Status Modal */}
      <PipelineStatusModal />
    </Layout.Content>
  );
};

export default ExecutionHistoryPage;