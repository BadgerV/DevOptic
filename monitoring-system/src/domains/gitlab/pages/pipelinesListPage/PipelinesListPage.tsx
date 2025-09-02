import { Layout, Result, Spin, Pagination } from 'antd';
import PipelineUnitCard from '../../components/PipelineUnitCard/PipelineUnitCard';
import { useGetPipelineUnitsQuery } from '../../store';
import { useState } from 'react';

const PipelinesListPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 10;
  const { data: pipelineUnits, isLoading, error } = useGetPipelineUnitsQuery();

  if (isLoading) {
    return (
      <Layout.Content style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Spin size="large" />
      </Layout.Content>
    );
  }

  if (error || !pipelineUnits) {
    return (
      <Layout.Content>
        <Result
          status="error"
          title="Error"
          subTitle="Failed to load pipeline units. Please try again."
        />
      </Layout.Content>
    );
  }

  const paginatedUnits = pipelineUnits.slice((currentPage - 1) * pageSize, currentPage * pageSize);

  return (
    <Layout.Content style={{ padding: '24px' }}>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: "2rem" }}>
        {paginatedUnits.map((unit : any) => (
          <PipelineUnitCard key={unit.pipeline_unit.id} unit={unit} />
        ))}
      </div>
      <Pagination
        current={currentPage}
        pageSize={pageSize}
        total={pipelineUnits.length}
        onChange={(page) => setCurrentPage(page)}
        style={{ marginTop: 16, textAlign: 'center' }}
      />
    </Layout.Content>
  );
};

export default PipelinesListPage;