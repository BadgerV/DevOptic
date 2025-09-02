import React from 'react';
import { Result, Button } from 'antd';
import { useNavigate } from 'react-router-dom';

const Unauthorized: React.FC = () => {
  const navigate = useNavigate();

  const handleGoBack = () => {
    // Go back to the previous page, or to home if no history
    navigate(-1);
  };

  const handleGoHome = () => {
    navigate('/');
  };

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      minHeight: '60vh',
      padding: '20px'
    }}>
      <Result
        status="403"
        title="403"
        subTitle="Sorry, you are not authorized to access this page. You need super admin privileges to view this content."
        extra={
          <div style={{ display: 'flex', gap: '8px', justifyContent: 'center' }}>
            <Button type="primary" onClick={handleGoHome}>
              Go Home
            </Button>
            <Button onClick={handleGoBack}>
              Go Back
            </Button>
          </div>
        }
      />
    </div>
  );
};

export default Unauthorized;