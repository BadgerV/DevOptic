import { useState } from "react";
import EndpointMainLeft from "../endpointMainLeft/EndpointMainLeft";
import EndpointMainRight from "../endpointMainRight/EndpointMainRight";
import Header from "../header/Header";
import "./endpointMain.css";

interface EndpointMainProps {
  isSidebarOpen: boolean;
}

const EndpointMain: React.FC<EndpointMainProps> = ({ isSidebarOpen }) => {
  const [filter, setFilter] = useState("");

  return (
    <div
      className="endpoint-main"
      style={{
        marginLeft: isSidebarOpen ? "17em" : "0",
        transition: "margin-left 0.3s ease",
      }}
    >
      <Header filter={filter} onFilterChange={setFilter} />

      <div className="endpoint-main-container">
        <EndpointMainLeft filter={filter} />
        <EndpointMainRight />
      </div>
    </div>
  );
};

export default EndpointMain;
