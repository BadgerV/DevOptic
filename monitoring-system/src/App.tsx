import React from "react";
import { Provider } from "react-redux";
import { ConfigProvider } from "antd";
import { BrowserRouter } from "react-router-dom";
import { store } from "@app/store";
// import AppRouter from '@app/router/AppRouter';
import { AppLayout } from "./app/router/AppRouter";
import "antd/dist/reset.css";
import "@styles/globals.css";

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <ConfigProvider
        theme={{
          token: {
            colorPrimary: "#1890ff",
            borderRadius: 6,
          },
        }}
      >
        <BrowserRouter>
          <AppLayout />
          {/* <AppRouter /> */}
        </BrowserRouter>
      </ConfigProvider>
    </Provider>
  );
};

export default App;
