// store/endpointMonitoringSlice.ts
import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface EndpointMonitoringState {
  selectedEndpointId: string | null;
  isSidebarOpen: boolean;
  filterTypeByRunning: "all" | true | false;
  totalNumberOfEndpoints: number;
  runningEndpoints: number;
}

const initialState: EndpointMonitoringState = {
  selectedEndpointId: null,
  isSidebarOpen: false,
  filterTypeByRunning: false,
  totalNumberOfEndpoints: 0,
  runningEndpoints: 0,
};

const endpointMonitoringSlice = createSlice({
  name: "endpointMonitoring",
  initialState,
  reducers: {
    setSelectedEndpointId: (state, action: PayloadAction<string | null>) => {
      state.selectedEndpointId = action.payload;
    },
    toggleSidebar: (state) => {
      state.isSidebarOpen = !state.isSidebarOpen;
    },
    setFilterTypeByRunning: (state, action) => {
      state.filterTypeByRunning = action.payload;
    },
    setTotalNumberOfEndpoints: (state, action) => {
      state.totalNumberOfEndpoints = action.payload;
    },
    setRunningEndpoints: (state, action) => {
      state.runningEndpoints = action.payload;
    },
  },
});

export const {
  setSelectedEndpointId,
  toggleSidebar,
  setFilterTypeByRunning,
  setRunningEndpoints,
  setTotalNumberOfEndpoints,
} = endpointMonitoringSlice.actions;
export default endpointMonitoringSlice.reducer;
